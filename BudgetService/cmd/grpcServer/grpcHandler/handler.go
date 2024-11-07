package grpchandler

import (
	userProto "budget/pkg/go/user"
	budgetProto "budget/pkg/go/budget"
	transactionProto "budget/pkg/go/transaction"
	reportProto "budget/pkg/go/report"
	"google.golang.org/grpc"
)
type GrpcHandler struct {
	server grpc.ServiceRegistrar
	user UserService
	budget BudgetService
	transaction TransactionService
	report ReportService
}

func NewGrpcHandler(grpcServer grpc.ServiceRegistrar, userSRV UserService,
	 budgetSRV BudgetService, txSRV TransactionService, reportSRV ReportService) *GrpcHandler{
	return &GrpcHandler{server: grpcServer, user: userSRV, 
		budget: budgetSRV, transaction: txSRV, report: reportSRV,}
}
func(h *GrpcHandler) RegisterServices(){
	h.registerUserService(h.server, h.user)
	h.registerBudgetService(h.server, h.budget)
	h.registerTxService(h.server, h.transaction)
	h.registerReportService(h.server, h.report)
}



func(h *GrpcHandler) registerUserService(server grpc.ServiceRegistrar, user UserService){
	userProto.RegisterUserServiceServer(server, &UserServiceServer{UserSRV: user})
}

func(h *GrpcHandler) registerBudgetService(server grpc.ServiceRegistrar, budget BudgetService){
	budgetProto.RegisterBudgetServiceServer(server, &BudgetServiceServer{BudgetSRV: budget})
}

func(h *GrpcHandler) registerTxService(server grpc.ServiceRegistrar, tx TransactionService){
	transactionProto.RegisterTransactionServiceServer(server, &TransactionServiceServer{TxSRV: tx})
}

func(h *GrpcHandler) registerReportService(server grpc.ServiceRegistrar, report ReportService){
	reportProto.RegisterReportServiceServer(server, &ReportServiceServer{ReportSRV: report})
}