package grpchandler

import (
	"github.com/justIGreK/MoneyKeeper/BudgetService/internal/models"
	transactionProto "github.com/justIGreK/MoneyKeeper/BudgetService/pkg/go/transaction"
	"context"
	"errors"
)

type TransactionServiceServer struct {
	transactionProto.UnimplementedTransactionServiceServer
	TxSRV TransactionService
}

type TransactionService interface {
	AddTransaction(ctx context.Context, transaction models.CreateTransaction) ([]string, error)
	GetTransaction(ctx context.Context, transactionID, userID string) (*models.Transaction, error)
	GetAllTransactions(ctx context.Context, userID string) ([]models.Transaction, error)
	GetTXByTimeFrame(ctx context.Context, userID string, timeframe models.CreateTimeFrame) ([]models.Transaction, error)
}

func (s *TransactionServiceServer) CreateTransaction(ctx context.Context, req *transactionProto.CreateTransactionRequest) (*transactionProto.CreateTransactionResponse, error) {
	tx := models.CreateTransaction{
		Category: req.Category,
		UserID:   req.UserId,
		Name:     req.Name,
		Cost:     float64(req.Cost),
	}
	notifications, err := s.TxSRV.AddTransaction(ctx, tx)
	if err != nil {
		return nil, err
	}
	if len(notifications) == 0 {
		return nil, nil
	}

	return &transactionProto.CreateTransactionResponse{
		Notifications: notifications,
	}, nil

}

func (s *TransactionServiceServer) GetTransaction(ctx context.Context, req *transactionProto.GetTransactionRequest) (*transactionProto.GetTransactionResponse, error) {
	tx, err := s.TxSRV.GetTransaction(ctx, req.TxId, req.UserId)
	if err != nil {
		return nil, err
	}
	if tx == nil {
		return nil, errors.New("transaction os not found")
	}

	return &transactionProto.GetTransactionResponse{
		Transaction: &transactionProto.Transaction{
			Id:       tx.ID,
			UserId:   tx.UserID,
			Category: tx.Category,
			Name:     tx.Name,
			Cost:     float32(tx.Cost),
			Date:     tx.Date.Format(DateTimeformat),
		},
	}, nil

}

func (s *TransactionServiceServer) GetTransactionList(ctx context.Context, req *transactionProto.GetTransactionListRequest) (*transactionProto.GetTransactionListResponse, error) {
	txs, err := s.TxSRV.GetAllTransactions(ctx, req.UserId)
	if err != nil {
		return nil, err
	}
	protoTxs := convertToProtoTxs(txs)
	return &transactionProto.GetTransactionListResponse{
		Transactions: protoTxs,
	}, nil

}

func convertToProtoTxs(txs []models.Transaction) []*transactionProto.Transaction {
	protoTxs := make([]*transactionProto.Transaction, len(txs))
	for i, b := range txs {
		protoTxs[i] = &transactionProto.Transaction{
			Id:       b.ID,
			UserId:   b.UserID,
			Category: b.Category,
			Name:     b.Name,
			Cost:     float32(b.Cost),
			Date:     b.Date.Format(DateTimeformat),
		}
	}
	return protoTxs
}

func (s *TransactionServiceServer) GetTXByTimeFrame(ctx context.Context, req *transactionProto.GetTXByTimeFrameRequest) (*transactionProto.GetTransactionListResponse, error) {

	txs, err := s.TxSRV.GetTXByTimeFrame(ctx, req.UserId, models.CreateTimeFrame{req.StartDate, req.EndDate})
	if err != nil {
		return nil, err
	}
	protoTxs := convertToProtoTxs(txs)
	return &transactionProto.GetTransactionListResponse{
		Transactions: protoTxs,
	}, nil
}
