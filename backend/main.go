package main

import (
	"go-gin-auth/config"
	"go-gin-auth/helpers"
	"go-gin-auth/router"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	envErr := godotenv.Load(".env")
	if envErr != nil {
		panic("error loading .env file: " + envErr.Error())
	}

	config.ConnectDB()

	migrateErr := helpers.MigrateDB()
	if migrateErr != nil {
		panic("error migrating database: " + migrateErr.Error())
	}

	r := router.SetupRouter()
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	runErr := r.Run(":" + port)
	if runErr != nil {
		panic("error running server: " + runErr.Error())
	}
}
