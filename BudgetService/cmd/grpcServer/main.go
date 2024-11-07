package main

import (
	"log"
	"net"

	"google.golang.org/grpc"

	"github.com/justIGreK/MoneyKeeper/BudgetServiceContract/go/_go"
)

func main() {

	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen on port 50051: %v", err)
	}

	grpcServer := grpc.NewServer()

}
