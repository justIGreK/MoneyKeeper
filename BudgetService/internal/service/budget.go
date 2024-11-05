package service

import (
	"budget/internal/models"
	"context"
	"errors"
	"fmt"
	"log"
	"time"
)

type BudgetRepository interface {
	AddBudget(ctx context.Context, budget models.Budget) (string, error)
	CloseBudget(ctx context.Context, budgetID string) error
	GetBudgetList(ctx context.Context, userID string) ([]models.Budget, error)
	GetBudget(ctx context.Context, userID, budgetID string) (*models.Budget, error)
}

type BudgetService struct {
	BudgetRepo BudgetRepository
	UserRepo   UserRepository
}

func NewBudgetService(repo BudgetRepository, user UserRepository) *BudgetService {
	return &BudgetService{BudgetRepo: repo, UserRepo: user}
}

const (
	Dateformat string = "2006-01-02"
)

func (s *BudgetService) AddBudget(ctx context.Context, budget models.CreateBudget) error {
	user, err := s.UserRepo.GetUser(ctx, budget.UserID)
	if err != nil {
		log.Println(err)
		return err
	}
	if user == nil {
		return errors.New("user not found")
	}
	
	endDate, err := time.Parse(Dateformat, budget.EndTime)
	if err != nil {
		log.Println(err)
		return fmt.Errorf("invalid date or time format: %v", err)
	}
	now := time.Now().UTC()
	if endDate.Before(now){
		return errors.New("Past time")
	}
	duration := endDate.Sub(now)
	days := duration.Hours() / 24.0
	if days < 1 {
		return errors.New("Invalid duration")
	}
	createBudget := models.Budget{
		UserID:      budget.UserID,
		Name:        budget.Name,
		Amount:      budget.Amount,
		DailyAmount: budget.Amount / float64(days),
		StartDate:   now,
		EndDate:     endDate,
		CreatedAt:   now,
		UpdatedAt:   now,
		IsActive:    true,
	}
	_, err = s.BudgetRepo.AddBudget(ctx, createBudget)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func (s *BudgetService) GetBudgetList(ctx context.Context, userID string) ([]models.Budget, error) {
	user, err := s.UserRepo.GetUser(ctx, userID)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	if user == nil {
		return nil, errors.New("user not found")
	}
	budgetList, err := s.BudgetRepo.GetBudgetList(ctx, userID)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return budgetList, nil
}
