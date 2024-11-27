package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"

	"fiber_api_v1/auth"
	"fiber_api_v1/db"
	"fiber_api_v1/game"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/pprof"
)

func main() {

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)

	// Подключаемся к базе данных
	if err := db.InitDB(); err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer func() {
		db.CloseDB()
		fmt.Print("Server Done")
	}()

	res, _ := auth.GeneratePasswdHash("qwerty")
	fmt.Print(res)

	// Создаем новое приложение Fiber
	app := fiber.New()

	// Добавляем профайлер
	app.Use(pprof.New())
	// Логирование запросов
	app.Use(logger.New())
	app.Use(cors.New())

	// Роуты
	app.Post("/login", auth.LoginHandler)

	app.Post("/register", auth.RegisterHandler)

	// Защищенный эндпоинт, требует авторизации
	app.Get("/profile", auth.Protected, func(c *fiber.Ctx) error {
		// Здесь мы можем получить имя пользователя из контекста
		username := c.Locals("username").(string)
		id := c.Locals("user_id").(int)
		return c.JSON(fiber.Map{"message": "Welcome to your profile", "username": username, "user_id": id})
	})

	app.Get("/start", auth.Protected, game.StartHandler)
	app.Post("/step", auth.Protected, game.GameStepHandler)

	// Запуск сервера

	go func() {
		log.Println(app.Listen(":3000"))
	}()
	<-quit

	if err := app.Shutdown(); err != nil {
		fmt.Println("closing error", err)
	} else {
		fmt.Println("app closed correct")
	}
}
