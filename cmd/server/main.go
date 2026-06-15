package main

import (
	"UserAgeAPI/config"
	db "UserAgeAPI/db/sqlc/generated"
	"UserAgeAPI/internal/handler"
	"UserAgeAPI/internal/logger"
	"UserAgeAPI/internal/repository"
	"UserAgeAPI/internal/routes"
	"UserAgeAPI/internal/service"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

func main() {
	zapLogger, err := logger.NewLogger()
	if err != nil {
		panic(err)
	}
	defer zapLogger.Sync()
	_ = godotenv.Load()

	pool, err := config.ConnectDB()
	if err != nil {
		zapLogger.Fatal("failed to connect database",
			zap.Error(err),
		)
	}
	defer pool.Close()

	queries := db.New(pool)

	repo := repository.NewUserRepository(queries)
	userService := service.NewUserService(repo, zapLogger)
	userHandler := handler.NewUserHandler(userService)

	app := fiber.New()

	routes.SetupRoutes(app, userHandler)

	zapLogger.Info("server running",
		zap.String("port", "3000"),
	)

	if err := app.Listen(":3000"); err != nil {
		zapLogger.Fatal("failed to start server",
			zap.Error(err),
		)
	}
}
