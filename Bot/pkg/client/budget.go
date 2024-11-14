package client

import (
	"context"
	"errors"

	budget "github.com/justIGreK/MoneyKeeper-Budget/pkg/go/budget"
	"github.com/justIGreK/MoneyKeeper/Bot/internal/models"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/wrapperspb"
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

func (bc *BudgetClient) AddBudget(ctx context.Context, budgetreq models.CreateBudget) (string, error) {
	req := &budget.AddBudgetRequest{
		UserId: budgetreq.UserID,
		Name:   budgetreq.Name,
		Limit:  budgetreq.Amount,
		Period: budgetreq.Period,
		Start:  budgetreq.Start,
		End:    budgetreq.End,
	}
	res, err := bc.client.AddBudget(ctx, req)
	if err != nil {
		return "", err
	}
	return res.BudgetId, nil
}

func (bc *BudgetClient) AddCategory(ctx context.Context, catreq models.CreateCategory) (*models.Budget, error) {
	req := &budget.AddCategoryRequest{
		BudgetId: catreq.BudgetID,
		UserId:   catreq.UserID,
		Name:     catreq.Name,
		Limit:    catreq.Limit,
	}
	res, err := bc.client.AddCategory(ctx, req)
	if err != nil {
		return nil, err
	}

	return &models.Budget{
		ID:       res.Budget.BudgetId,
		Name:     res.Budget.Name,
		Limit:    res.Budget.Limit,
		Start:    res.Budget.Start,
		End:      res.Budget.End,
		Category: bc.convertToCategory(res.Budget.Category),
	}, nil
}

func (bc *BudgetClient) GetBudget(ctx context.Context, userID, budgetID string) (*models.Budget, error) {
	req := &budget.GetBudgetRequest{UserId: userID, BudgetId: budgetID}
	res, err := bc.client.GetBudget(ctx, req)
	if err != nil {
		return nil, err
	}
	if req == nil {
		return nil, errors.New("budget is not find")
	}
	return &models.Budget{
		ID:       res.Budget.BudgetId,
		Name:     res.Budget.Name,
		Limit:    res.Budget.Limit,
		Start:    res.Budget.Start,
		End:      res.Budget.End,
		Category: bc.convertToCategory(res.Budget.Category),
	}, nil
}

func (bc *BudgetClient) GetBudgetList(ctx context.Context, userID string) ([]models.Budget, error) {
	req := &budget.GetBudgetListRequest{UserId: userID}
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

func (bc *BudgetClient) convertToBudgets(protoBudgets []*budget.Budget) []models.Budget {
	budgets := make([]models.Budget, len(protoBudgets))
	for i, b := range protoBudgets {
		budgets[i] = models.Budget{
			ID:       b.BudgetId,
			Name:     b.Name,
			Limit:    b.Limit,
			Start:    b.Start,
			End:      b.End,
			Category: bc.convertToCategory(b.Category),
		}
	}
	return budgets
}

func (bc *BudgetClient) convertToCategory(protoCategs []*budget.Category) []models.Category {
	categs := make([]models.Category, len(protoCategs))
	for i, b := range protoCategs {
		categs[i] = models.Category{
			ID:    b.CategoryId,
			Name:  b.Name,
			Limit: b.Limit,
		}
	}
	return categs
}

func (bc *BudgetClient) DeleteBudget(ctx context.Context, userID, budgetID string) error {
	req := &budget.DeleteBudgetRequest{UserId: userID, BudgetId: budgetID}
	_, err := bc.client.DeleteBudget(ctx, req)
	if err != nil {
		return err
	}
	return nil
}

func (bc *BudgetClient) DeleteCategory(ctx context.Context, userID, budgetID, categoryID string) error {
	req := &budget.DeleteCategoryRequest{UserId: userID, BudgetId: budgetID, CategoryId: categoryID}
	_, err := bc.client.DeleteCategory(ctx, req)
	if err != nil {
		return err
	}
	return nil
}

var (
	optionalSField *wrapperspb.StringValue = nil
	optionalDField *wrapperspb.DoubleValue = nil
)

func (bc *BudgetClient) UpdateBudget(ctx context.Context, updates models.UpdateBudget)(*models.Budget, error) {
	req := &budget.UpdateBudgetRequest{
		Update: &budget.UpdateBudget{
			BudgetId: updates.BudgetID,
			UserId:   updates.UserID,
		},
	}
	if updates.Name != nil {
		req.Update.Name = wrapperspb.String(*updates.Name)
	} else {
		req.Update.Name = optionalSField
	}
	if updates.Limit != nil {
		req.Update.Limit = wrapperspb.Double(*updates.Limit)
	} else {
		req.Update.Limit = optionalDField
	}
	if updates.Start != nil {
		req.Update.Start = wrapperspb.String(*updates.Start)
	} else {
		req.Update.Start = optionalSField
	}
	if updates.End != nil {
		req.Update.End = wrapperspb.String(*updates.End)
	} else {
		req.Update.End = optionalSField
	}
	res, err := bc.client.UpdateBudget(ctx, req)
	if err != nil {
		return nil, err
	}
	return &models.Budget{
		ID:       res.Budget.BudgetId,
		Name:     res.Budget.Name,
		Limit:    res.Budget.Limit,
		Start:    res.Budget.Start,
		End:      res.Budget.End,
		Category: bc.convertToCategory(res.Budget.Category),
	}, nil
}

func (bc *BudgetClient) UpdateCategory(ctx context.Context, updates models.UpdateCategory) (*models.Budget, error) {
	req := &budget.UpdateCategoryRequest{
		Update: &budget.UpdateCategory{
			BudgetId:   updates.BudgetID,
			UserId:     updates.UserID,
			CategoryId: updates.CategoryID,
		},
	}
	if updates.Name != nil {
		req.Update.Name = wrapperspb.String(*updates.Name)
	} else {
		req.Update.Name = optionalSField
	}
	if updates.Limit != nil {
		req.Update.Limit = wrapperspb.Double(*updates.Limit)
	} else {
		req.Update.Limit = optionalDField
	}
	res, err := bc.client.UpdateCategory(ctx, req)
	if err != nil {
		return nil, err
	}
	return &models.Budget{
		ID:       res.Budget.BudgetId,
		Name:     res.Budget.Name,
		Limit:    res.Budget.Limit,
		Start:    res.Budget.Start,
		End:      res.Budget.End,
		Category: bc.convertToCategory(res.Budget.Category),
	}, nil
}
