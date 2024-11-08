package main

import (
	"context"
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	grpchandler "github.com/justIGreK/MoneyKeeper/BudgetService/cmd/grpcServer/grpcHandler"
	mongorep "github.com/justIGreK/MoneyKeeper/BudgetService/internal/repository/mongo"
	"github.com/justIGreK/MoneyKeeper/BudgetService/internal/service"
	
)
func main() {
	ctx := context.Background()
	client := mongorep.CreateMongoClient(ctx)
	userRepo := mongorep.NewUserRepository(client)
	budgetRepo := mongorep.NewBudgetRepository(client)
	txRepo := mongorep.NewTransactionRepository(client)
	reportRepo := mongorep.NewReportRepository(client)

	userSRV := service.NewUserService(userRepo)
	budgetSRV := service.NewBudgetService(budgetRepo, userRepo)
	txSRV := service.NewTransactionService(txRepo, userRepo, budgetRepo)
	reportSRV := service.NewReportService(reportRepo, txRepo, budgetRepo, userRepo)
	lis, err := net.Listen("tcp", ":50052")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()

	handler := grpchandler.NewGrpcHandler(grpcServer, userSRV, budgetSRV, txSRV, reportSRV)
	handler.RegisterServices()
	reflection.Register(grpcServer)

	log.Printf("Starting gRPC server on :50052")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
