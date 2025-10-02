package main

import (
	"Project/internal/app/ds"
	"Project/internal/app/dsn"
	"log"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	// Загрузка .env
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	db, err := gorm.Open(postgres.Open(dsn.FromEnv()), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	// Миграции
	err = db.AutoMigrate(
		&ds.HeaterProduct{},
		&ds.HeatersProductRequest{},
		&ds.RequestHeater{},
		&ds.User{},
	)
	if err != nil {
		panic("cant migrate db")
	}
}
