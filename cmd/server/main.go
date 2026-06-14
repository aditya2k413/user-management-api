package main

import (
	"UserAgeAPI/config"
	db "UserAgeAPI/db/sqlc/generated"
	"UserAgeAPI/internal/handler"
	"UserAgeAPI/internal/repository"
	"UserAgeAPI/internal/routes"
	"UserAgeAPI/internal/service"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	pool, err := config.ConnectDB()
	if err != nil {
		log.Fatal(err)
	}
	defer pool.Close()

	queries := db.New(pool)

	repo := repository.NewUserRepository(queries)
	userService := service.NewUserService(repo)
	userHandler := handler.NewUserHandler(userService)

	app := fiber.New()

	routes.SetupRoutes(app, userHandler)

	log.Println("Server running on :3000")

	log.Fatal(app.Listen(":3000"))
}
