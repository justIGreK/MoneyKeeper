package main

import (
	"budget/cmd/config"
	"budget/cmd/handler"
	mongorep "budget/internal/repository/mongo"
	"budget/internal/service"
	"context"
	"log"
	"net/http"
	"os"
	"time"
)
const (
	RWTimeout   = 10
	IdleTimeout = 60
)

func main() {
	ctx := context.Background()
	config.LoadEnv()
	client := mongorep.CreateMongoClient(ctx)
	userRepo := mongorep.NewUserRepository(client)
	budgetRepo := mongorep.NewBudgetRepository(client)
	txRepo := mongorep.NewTransactionRepository(client)

	txSRV := service.NewTransactionService(txRepo, userRepo, budgetRepo)
	budgetSRV := service.NewBudgetService(budgetRepo, userRepo)
	userSRV := service.NewUserService(userRepo)
	handler := handler.NewHandler(userSRV, budgetSRV, txSRV)

	srv := &http.Server{
		Addr:         os.Getenv("PORT"),
		Handler:      handler.InitRoutes(),
		ReadTimeout:  RWTimeout * time.Second,
		WriteTimeout: RWTimeout * time.Second,
		IdleTimeout:  IdleTimeout * time.Second,
	}
	err := srv.ListenAndServe()
	if err != nil {
		log.Fatalf("Server error", err)
	}
}
