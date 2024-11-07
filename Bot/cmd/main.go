package main

import (
	"Bot/cmd/config"
	handler "Bot/cmd/handler"
	"log"
	"net/http"
	"os"
)

func main() {
	config.LoadEnv()
	botToken := os.Getenv("BOT_TOKEN")
	if botToken == "" {
		log.Fatal("TELEGRAM_BOT_TOKEN is not set")
	}
	handler := handler.NewTelegramHandler(botToken)
	http.HandleFunc("/webhook", handler.HandleWebhook)
	log.Println("Starting server on :7777")
	if err := http.ListenAndServe(":7777", nil); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
