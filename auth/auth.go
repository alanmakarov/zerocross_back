package auth

import (
	"context"
	"log"

	"fiber_api_v1/db"
	"fiber_api_v1/models"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

func GeneratePasswdHash(pass string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)
	if err != nil {
		log.Fatal(err)
	}
	return string(hashedPassword), err
}

func RegisterHandler(c *fiber.Ctx) error {
	var request models.RegisterRequest
	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Invalid request body"})
	}

	var count uint64
	err := db.Pool.QueryRow(context.Background(), "SELECT count(*) FROM users WHERE username=$1", request.Username).Scan(&count)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Database error"})
	}

	if count != 0 {
		return c.Status(fiber.StatusNotAcceptable).JSON(fiber.Map{"message": "User already exists"})
	}

	hashedPassword, err := GeneratePasswdHash(request.Password)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "can generate passwd hash"})
	}

	commandTag, err := db.Pool.Exec(context.Background(), "INSERT INTO users (username, password) VALUES ($1,$2)", request.Username, hashedPassword)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Database error can't insert user"})
	}
	if commandTag.RowsAffected() != 1 {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Database error can't insert user"})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "user added"})

}

// LoginHandler — обработка запроса на авторизацию
func LoginHandler(c *fiber.Ctx) error {
	var request models.LoginRequest
	if err := c.BodyParser(&request); err != nil && request.Username != "" && request.Password != "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Invalid or empty request body"})
	}

	// Получаем хеш пароля из базы
	var storedHash string
	var user_id int
	err := db.Pool.QueryRow(context.Background(), "SELECT id, password FROM users WHERE username=$1", request.Username).Scan(&user_id, &storedHash)
	if err != nil {
		log.Fatal("Database error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Database error"})
	}

	// Сравниваем пароль с хешем
	err = bcrypt.CompareHashAndPassword([]byte(storedHash), []byte(request.Password))
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "Invalid credentials"})
	}

	// Генерируем JWT токен
	token, err := GenerateToken(request.Username, user_id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Could not generate token"})
	}

	// Возвращаем токен клиенту
	return c.Status(fiber.StatusOK).JSON(models.LoginResponse{
		Message: "Login successful",
		Token:   token, // Возвращаем токен
	})
	// Успех
	// return c.Status(fiber.StatusOK).JSON(models.LoginResponse{Message: "Login successful"})
}
