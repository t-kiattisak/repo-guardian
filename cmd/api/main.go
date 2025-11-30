package main

import (
	"log"
	"time"

	"repo-guardian/internal/user/handler"
	"repo-guardian/internal/user/repository"
	"repo-guardian/internal/user/usecase"

	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()

	timeoutContext := 2 * time.Second

	userRepo := repository.NewMemoryUserRepository()
	userUsecase := usecase.NewUserUsecase(userRepo, timeoutContext)
	handler.NewUserHandler(app, userUsecase)

	log.Fatal(app.Listen(":3000"))
}
