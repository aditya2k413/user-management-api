package main

import (
	"UserAgeAPI/config"
	db "UserAgeAPI/db/sqlc/generated"
	"UserAgeAPI/internal/models"
	"UserAgeAPI/internal/repository"
	"UserAgeAPI/internal/service"
	"context"
	"fmt"
	"log"

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

	user, err := userService.CreateUser(
		context.Background(),
		models.CreateUserRequest{
			Name: "Aditya",
			Dob:  "2003-01-13",
		},
	)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Created User: %+v\n", user)

	fetchedUser, err := userService.GetUser(
		context.Background(),
		user.ID,
	)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Fetched User: %+v\n", fetchedUser)

	users, err := userService.ListUsers(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("All Users: %+v\n", users)
}
