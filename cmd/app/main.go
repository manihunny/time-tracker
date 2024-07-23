package main

import (
	"github.com/joho/godotenv"
	"log/slog"
	"main/internal/app"
	"main/internal/config"
	"os"
)

//	@title			Time Tracker API
//	@version		1.0
//	@description	API Server for Time Tracker application

//	@schemes	http
//	@BasePath	/

func main() {
	err := godotenv.Load()
	if err != nil {
		slog.Error("Error loading .env file. Error message: ", "msg", err)
		os.Exit(1)
	}

	cfg := config.Config{
		ServerHost:      os.Getenv("SERVER_HOST"),
		ServerPort:      os.Getenv("SERVER_PORT"),
		DBDriver:        os.Getenv("DB_DRIVER"),
		DBHost:          os.Getenv("DB_HOST"),
		DBPort:          os.Getenv("DB_PORT"),
		DBUser:          os.Getenv("DB_USER"),
		DBPassword:      os.Getenv("DB_PASSWORD"),
		DBName:          os.Getenv("DB_NAME"),
		DBMigrationsDir: os.Getenv("DB_MIGRATIONS_DIR"),
	}

	app.Run(cfg)
}
