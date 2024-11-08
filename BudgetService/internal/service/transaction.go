package service

import (
	"github.com/justIGreK/MoneyKeeper/BudgetService/internal/models"
	"context"
	"errors"
	"fmt"
	"log"
	"time"
)

type TransactionRepository interface {
	AddTransaction(ctx context.Context, transaction models.Transaction) (string, error)
	GetTransaction(ctx context.Context, transactionID, userID string) (*models.Transaction, error)
	GetAllTransactions(ctx context.Context, userID string) ([]models.Transaction, error)
	GetTXByTimeFrame(ctx context.Context, userID string, dateFrame models.TimeFrame) ([]models.Transaction, error)
}

type TransactionService struct {
	TransactionRepo TransactionRepository
	UserRepo        UserRepository
	BudgetRepo      BudgetRepository
}

func NewTransactionService(transRepo TransactionRepository, userRepo UserRepository, budgetRepo BudgetRepository) *TransactionService {
	return &TransactionService{TransactionRepo: transRepo,
		UserRepo: userRepo, BudgetRepo: budgetRepo}
}

const (
	NoCategory = "other"
)

func (s *TransactionService) AddTransaction(ctx context.Context, transaction models.CreateTransaction) ([]string, error) {
	if transaction.Cost < 0{
		return nil, errors.New("cost cant be below 0")
	}
	user, err := s.UserRepo.GetUser(ctx, transaction.UserID)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	if user == nil {
		return nil, errors.New("user not found")
	}
	if transaction.Category == "" {
		transaction.Category = NoCategory
	}
	createTransaction := models.Transaction{
		UserID:   transaction.UserID,
		Category: transaction.Category,
		Name:     transaction.Name,
		Cost:     transaction.Cost,
		Date:     time.Now().UTC(),
	}
	_, err = s.TransactionRepo.AddTransaction(ctx, createTransaction)
	if err != nil {
		return nil, err
	}
	notifications, err := s.checkLimits(ctx, transaction.UserID)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return notifications, nil
}

func (s *TransactionService) GetTransaction(ctx context.Context, transactionID, userID string) (*models.Transaction, error) {
	user, err := s.UserRepo.GetUser(ctx, userID)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	if user == nil {
		return nil, errors.New("user not found")
	}
	trans, err := s.TransactionRepo.GetTransaction(ctx, transactionID, userID)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return trans, nil
}

func (s *TransactionService) checkLimits(ctx context.Context, userID string) ([]string, error) {
	budgets, err := s.BudgetRepo.GetBudgetList(ctx, userID)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	if len(budgets) == 0 {
		return nil, nil
	}
	now := time.Now().UTC()
	CurrTframe := models.TimeFrame{
		StartDate: time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location()),
		EndDate:   time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 999999999, now.Location()),
	}
	txs, err := s.TransactionRepo.GetTXByTimeFrame(ctx, userID, CurrTframe)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	sum := 0.0
	for _, tx := range txs {
		sum += tx.Cost
	}
	warningNotifications := []string{}
	for _, budget := range budgets {
		if sum > budget.DailyAmount && now.Before(budget.EndDate) {
			notification := fmt.Sprintf("daily budget of %v is exceeded by %v", budget.Name, sum-budget.DailyAmount)
			warningNotifications = append(warningNotifications, notification)
		}
	}

	return warningNotifications, nil
}

func (s *TransactionService) GetAllTransactions(ctx context.Context, userID string) ([]models.Transaction, error) {
	user, err := s.UserRepo.GetUser(ctx, userID)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	if user == nil {
		return nil, errors.New("user not found")
	}

	txs, err := s.TransactionRepo.GetAllTransactions(ctx, userID)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return txs, nil
}
func (s *TransactionService) GetTXByTimeFrame(ctx context.Context, userID string, timeframe models.CreateTimeFrame) ([]models.Transaction, error) {
	user, err := s.UserRepo.GetUser(ctx, userID)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	if user == nil {
		return nil, errors.New("user not found")
	}
	var tf models.TimeFrame
	if timeframe.StartDate == "" {
		tf.StartDate = time.Unix(0, 0)
	} else {
		tf.StartDate, err = time.Parse(Dateformat, timeframe.StartDate)
		if err != nil{
			log.Println(err)
			return nil, err
		}
	}
	if timeframe.EndDate == "" {
		tf.EndDate = time.Now().AddDate(10000, 0, 0)
	} else {
		tf.EndDate, err = time.Parse(Dateformat, timeframe.EndDate)
		if err != nil{
			log.Println(err)
			return nil, err
		}
	}

	txs, err := s.TransactionRepo.GetTXByTimeFrame(ctx, userID, tf)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return txs, nil
}
