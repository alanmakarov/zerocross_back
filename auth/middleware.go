package auth

import (
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"strings"
)

// Protected — middleware для защиты эндпоинтов с JWT
func Protected(c *fiber.Ctx) error {
	// Извлекаем токен из заголовка Authorization
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Missing Authorization header"})
	}

	// Токен должен быть в формате "Bearer <token>"
	tokenString := strings.TrimPrefix(authHeader, "Bearer ")

	// Парсим токен
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		// Проверка подписи с использованием нашего секретного ключа
		return SecretKey, nil
	})

	if err != nil || !token.Valid {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid or expired token"})
	}

	// Записываем информацию о пользователе из токена в контекст
	claims, ok := token.Claims.(*Claims)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid token claims"})
	}

	c.Locals("username", claims.Username) // Добавляем имя пользователя в контекст
	c.Locals("user_id", claims.UserID)
	return c.Next()
}
