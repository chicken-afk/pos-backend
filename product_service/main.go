package main

import (
	"log"

	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables from .env file
	if err := godotenv.Load(); err != nil {
		log.Println("⚠️ .env file not found, relying on system environment")
	}
}
