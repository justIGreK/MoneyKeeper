package client

import (
	"context"
	"errors"

	"github.com/justIGreK/MoneyKeeper/Bot/internal/models"
	budget "github.com/justIGreK/MoneyKeeper/BudgetService/pkg/go/budget"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type BudgetClient struct {
	client budget.BudgetServiceClient
}

func NewBudgetClient(serviceAddress string) (*BudgetClient, error) {
	conn, err := grpc.NewClient(serviceAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	return &BudgetClient{
		client: budget.NewBudgetServiceClient(conn),
	}, nil
}

func (bc *BudgetClient) CreateBudget(ctx context.Context, budgetreq models.CreateBudget) (string, error) {
	req := &budget.CreateBudgetRequest{
		UserId:  budgetreq.UserID,
		Name:    budgetreq.Name,
		Amount:  budgetreq.Amount,
		Endtime: budgetreq.EndTime,
	}
	res, err := bc.client.CreateBudget(ctx, req)
	if err != nil {
		return "", err
	}
	return res.BudgetId, nil
}

func (bc *BudgetClient) GetBudgetList(ctx context.Context, id string) ([]models.Budget, error) {
	req := &budget.GetBudgetListRequest{UserId: id}
	res, err := bc.client.GetBudgetList(ctx, req)
	if err != nil {
		return nil, err
	}
	if req == nil {
		return nil, errors.New("user is not find")
	}
	budgets := bc.convertToBudgets(res.Budgets)
	return budgets, nil
}

func (bc *BudgetClient)convertToBudgets(protoBudgets []*budget.Budget) []models.Budget {
	budgets := make([]models.Budget, len(protoBudgets))
	for i, b := range protoBudgets {
		budgets[i] = models.Budget{
			ID:          b.Id,
			Name:        b.Name,
			Amount:      b.Amount,
			DailyAmount: b.DailyAmount,
			StartDate:   b.StartDate,
			EndDate:     b.EndDate,
			CreatedAt:   b.CreatedAt,
			UpdatedAt:   b.UpdatedAt,
			IsActive:    b.IsActive,
		}
	}
	return budgets
}
