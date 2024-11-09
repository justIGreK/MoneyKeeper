package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/justIGreK/MoneyKeeper/Bot/cmd/config"
	handler "github.com/justIGreK/MoneyKeeper/Bot/cmd/handler"
	"github.com/justIGreK/MoneyKeeper/Bot/internal/repository"
	"github.com/justIGreK/MoneyKeeper/Bot/internal/service"
	"github.com/justIGreK/MoneyKeeper/Bot/pkg/client"
)

func main() {
	ctx := context.Background()
	config.LoadEnv()
	botToken := os.Getenv("BOT_TOKEN")
	if botToken == "" {
		log.Fatal("TELEGRAM_BOT_TOKEN is not set")
	}
	user, err := client.NewUserClient("localhost:50052")
	if err != nil{
		log.Fatal(err)
	}
	budget, err := client.NewBudgetClient("localhost:50052")
	if err != nil{
		log.Fatal(err)
	}
	tx, err := client.NewTransactionClient("localhost:50052")
	if err != nil{
		log.Fatal(err)
	}
	report, err := client.NewReportClient("localhost:50052")
	if err != nil{
		log.Fatal(err)
	}
	db := repository.CreateMongoClient(ctx)
	userDB := repository.NewUserRepo(db)
	srv := service.NewService(ctx, user, userDB, budget, tx, report)
	handler := handler.NewTelegramHandler(botToken, srv)
	http.HandleFunc("/webhook", handler.HandleWebhook)
	log.Println("Starting server on :7777")
	if err := http.ListenAndServe(":7777", nil); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
