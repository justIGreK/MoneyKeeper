package service

import (
	"context"

	"github.com/justIGreK/MoneyKeeper/Bot/pkg/client"
)

type Service struct {
	ctx context.Context
	user   *client.UserClient
	userDB UserDB
	budget *client.BudgetClient
	tx *client.TransactionClient
	report *client.ReportClient
}

type UserDB interface {
	AddUser(ctx context.Context, chatID, userID string) error 
	GetUserID(ctx context.Context, chatID string) (string, error)
}

func NewService(context context.Context, user *client.UserClient, userDB UserDB,
	tx *client.TransactionClient, report *client.ReportClient, budget *client.BudgetClient) *Service {
	return &Service{
		ctx: context, 
		user: user, 
		userDB: userDB,
	 	budget: budget,
		tx: tx,
		report: report,
	}
}
