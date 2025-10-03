package main

import (
	"Project/internal/api"
	"Project/internal/app/dsn"
	"Project/internal/app/repository"
	"log"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	log.Println("Environment variables loaded successfully")

	repo, err := repository.NewRepository(dsn.FromEnv())
	if err != nil {
		log.Fatalf("failed to create repository: %v", err)
	}

	api.StartServer(repo)
}
