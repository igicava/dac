package main

import (
	"log"
	"path/filepath"

	"dac/orchestrator/http/server"

	"github.com/joho/godotenv"
)

func main() {
	LoadEnv()
	server.Run()
}

func LoadEnv() {
	envPath := filepath.Join("..", ".env")
	err := godotenv.Load(envPath)
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
}
