package main

import (
	"backend/internal/agent"
	"backend/internal/chat"
	"backend/internal/http/rest"
	"backend/internal/storage"
	"log"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	storage := storage.NewStorage()
	agent := agent.NewAgent()

	chatService := chat.NewService(storage, agent)

	r := rest.Handler(chatService)

	r.Run(":8080") // Run on port 8080
}
