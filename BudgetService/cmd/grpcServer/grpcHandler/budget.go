package grpchandler

import (
	budgetProto "github.com/justIGreK/MoneyKeeper/BudgetService/pkg/go/budget"
	"context"

	"github.com/justIGreK/MoneyKeeper/BudgetService/internal/models"
)

type BudgetServiceServer struct {
	budgetProto.UnimplementedBudgetServiceServer
	BudgetSRV BudgetService
}

type BudgetService interface {
	AddBudget(ctx context.Context, budget models.CreateBudget) (string, error)
	GetBudgetList(ctx context.Context, userID string) ([]models.Budget, error)
}

func (s *BudgetServiceServer) CreateBudget(ctx context.Context, req *budgetProto.CreateBudgetRequest) (*budgetProto.CreateBudgetResponse, error) {
	createBudget := models.CreateBudget{
		UserID:  req.UserId,
		Name:    req.Name,
		Amount:  float64(req.Amount),
		EndTime: req.Endtime,
	}
	budgetID, err := s.BudgetSRV.AddBudget(ctx, createBudget)
	if err != nil {
		return nil, err
	} else {
		return &budgetProto.CreateBudgetResponse{
			BudgetId: budgetID,
		}, nil
	}

}

func (s *BudgetServiceServer) GetBudgetList(ctx context.Context, req *budgetProto.GetBudgetListRequest) (*budgetProto.GetBudgetListResponse, error) {

	budgets, err := s.BudgetSRV.GetBudgetList(ctx, req.UserId)
	protobudgets := convertToProtoBudgets(budgets)
	if err != nil {
		return nil, err
	} else {
		return &budgetProto.GetBudgetListResponse{
			Budgets: protobudgets,
		}, nil
	}

}

var(
	Dateformat string = "2006-01-02"
	DateTimeformat string = "2006-01-02T15:04:05"
)
func convertToProtoBudgets(budgets []models.Budget) []*budgetProto.Budget {
	protoBudgets := make([]*budgetProto.Budget, len(budgets))
	for i, b := range budgets {
		protoBudgets[i] = &budgetProto.Budget{
			Id:          b.ID,
			Name:        b.Name,
			Amount:      float32(b.Amount),
			DailyAmount: float32(b.DailyAmount),
			StartDate:   b.StartDate.Format(Dateformat),
			EndDate:     b.EndDate.Format(Dateformat),
			CreatedAt:   b.CreatedAt.Format(Dateformat),
			UpdatedAt:   b.UpdatedAt.Format(Dateformat),
			IsActive:    b.IsActive,
		}
	}
	return protoBudgets
}
