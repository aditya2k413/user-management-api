package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"UserAgeAPI/config"
	db "UserAgeAPI/db/sqlc/generated"

	"github.com/jackc/pgx/v5/pgtype"
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

	user, err := queries.CreateUser(context.Background(), db.CreateUserParams{
		Name: "Aditya",
		Dob: pgtype.Date{
			Time:  time.Date(2003, 1, 13, 0, 0, 0, 0, time.UTC),
			Valid: true,
		},
	})

	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Created User: %+v\n", user)

	fetchedUser, err := queries.GetUser(context.Background(), user.ID)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Fetched User: %+v\n", fetchedUser)
}
