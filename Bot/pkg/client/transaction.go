package client

import (
	"context"

	tx "github.com/justIGreK/MoneyKeeper-Transaction/pkg/go/transaction"
	"github.com/justIGreK/MoneyKeeper/Bot/internal/models"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

type TransactionClient struct {
	client tx.TransactionServiceClient
}

func NewTransactionClient(serviceAddress string) (*TransactionClient, error) {
	conn, err := grpc.NewClient(serviceAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	return &TransactionClient{
		client: tx.NewTransactionServiceClient(conn),
	}, nil
}

func (tc *TransactionClient) AddTransaction(ctx context.Context, txreq models.CreateTransaction) (string, error) {
	req := &tx.CreateTransactionRequest{
		UserId:   txreq.UserID,
		Name:     txreq.Name,
		Cost:     txreq.Cost,
		Category: txreq.Category,
	}
	if txreq.Date != nil {
		req.Date = wrapperspb.String(*txreq.Date)
	} else {
		req.Date = optionalSField
	}
	res, err := tc.client.CreateTransaction(ctx, req)
	
	if err != nil {
		return "", err
	}
	return res.TxId, nil
}

func (tc *TransactionClient) GetTransactionList(ctx context.Context, userID string) ([]models.Transaction, error) {
	req := &tx.GetTransactionListRequest{
		UserId: userID,
	}
	res, err := tc.client.GetTransactionList(ctx, req)
	if err != nil {
		return nil, err
	}
	txs := tc.convertToTxs(res.Transactions)
	return txs, nil
}

func (tc *TransactionClient) DeleteTransaction(ctx context.Context, userID, txID string) error {
	req := &tx.DeleteTransactionRequest{
		UserId: userID,
		TxId: txID,
	}
	_, err := tc.client.DeleteTransaction(ctx, req)
	if err != nil {
		return err
	}
	return nil
}

func (tc *TransactionClient) convertToTxs(protoBudgets []*tx.Transaction) []models.Transaction {
	txs := make([]models.Transaction, len(protoBudgets))
	for i, b := range protoBudgets {
		txs[i] = models.Transaction{
			ID:       b.Id,
			UserID:   b.UserId,
			Category: b.Category,
			Name:     b.Name,
			Cost:     b.Cost,
			Date:     b.Date,
		}
	}
	return txs
}

func (tc *TransactionClient) GetTransaction(ctx context.Context, userID, txID string) (*models.Transaction, error) {
	req := &tx.GetTransactionRequest{
		UserId: userID,
		TxId:   txID,
	}
	res, err := tc.client.GetTransaction(ctx, req)
	if err != nil {
		return nil, err
	}

	return &models.Transaction{
		ID:       res.Transaction.Id,
		UserID:   res.Transaction.UserId,
		Category: res.Transaction.Category,
		Name:     res.Transaction.Name,
		Cost:     res.Transaction.Cost,
		Date:     res.Transaction.Date,
	}, nil
}

func (tc *TransactionClient) UpdateTransaction(ctx context.Context, updates models.UpdateTransaction) (*models.Transaction, error) {
	req := &tx.UpdateTransactionRequest{
		UserId: updates.UserID,
		TxId:   updates.ID,
	}
	if updates.Name != nil{
		req.Name = wrapperspb.String(*updates.Name)
	}else{
		req.Name = optionalSField
	}
	if updates.Category != nil{
		req.Category = wrapperspb.String(*updates.Category)
	}else{
		req.Category = optionalSField
	}
	if updates.Cost!= nil{
		req.Cost = wrapperspb.Double(*updates.Cost)
	}else{
		req.Cost = optionalDField
	}
	if updates.Date!= nil{
		req.Date = wrapperspb.String(*updates.Date)
	}else{
		req.Date = optionalSField
	}
	if updates.Time!= nil{
		req.Time = wrapperspb.String(*updates.Time)
	}else{
		req.Time = optionalSField
	}
	res, err := tc.client.UpdateTransaction(ctx, req)
	if err != nil {
		return nil, err
	}
	
	return &models.Transaction{
		ID:       res.Transaction.Id,
		UserID:   res.Transaction.UserId,
		Category: res.Transaction.Category,
		Name:     res.Transaction.Name,
		Cost:     res.Transaction.Cost,
		Date:     res.Transaction.Date,
	}, nil
}

func (tc *TransactionClient) GetTXByTimeFrame(ctx context.Context, userID, start, end string) ([]models.Transaction, error) {
	req := &tx.GetTXByTimeFrameRequest{
		UserId:    userID,
		StartDate: start,
		EndDate:   end,
	}
	res, err := tc.client.GetTXByTimeFrame(ctx, req)
	if err != nil {
		return nil, err
	}

	txs := tc.convertToTxs(res.Transactions)

	return txs, nil
}
