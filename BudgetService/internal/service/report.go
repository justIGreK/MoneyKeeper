package service

import (
	"budget/internal/models"
	"context"
	"errors"
	"fmt"
	"log"
	"time"
)

type ReportRepository interface {

}

type ReportService struct {
	ReportRepo  ReportRepository
	Transaction TransactionRepository
	BudgetRepo  BudgetRepository
	UserRepo    UserRepository
}

func NewReportService(report ReportRepository, tx TransactionRepository,
	budget BudgetRepository, user UserRepository) *ReportService {
	return &ReportService{ReportRepo: report, Transaction: tx,
		BudgetRepo: budget, UserRepo: user}
}

func (s *ReportService) GetPeriodSummary(ctx context.Context, userID, period string) (*models.ReportResponse, error) {
	user, err := s.UserRepo.GetUser(ctx, userID)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	if user == nil {
		return nil, errors.New("user not found")
	}
	startDate, endDate := s.getPeriodDates(period)
	if startDate.IsZero() {
		return nil, errors.New("invalid period")
	}
	txs, err := s.Transaction.GetTXByTimeFrame(ctx, userID, models.TimeFrame{startDate, endDate})
	if err != nil {
		log.Println(err)
		return nil, err
	}

	report := models.ReportResponse{
		UserID:     userID,
		Period:     fmt.Sprintf("%s - %s", startDate.Format("2006-01-02"), endDate.Format("2006-01-02")),
		Categories: []models.CategoryReport{},
	}

	categoryTotals := make(map[string]float64)
	categoryCounts := make(map[string]int)
	transactionCount := 0
	totalSpent := 0.0

	for _, txn := range txs {
        totalSpent += txn.Cost
        transactionCount++

        if txn.Category != "" {
            categoryTotals[txn.Category] += txn.Cost
			categoryCounts[txn.Category]++ 
        }
    }

    report.TotalSpent = totalSpent
    report.TransactionCount = transactionCount

    for category, total := range categoryTotals {
        report.Categories = append(report.Categories, models.CategoryReport{
            Name:  category,
            Total: total,
			Count: categoryCounts[category],
        })
    }

    return &report, nil
}


func (s *ReportService) getPeriodDates(period string) (time.Time, time.Time) {
	now := time.Now().UTC()
	switch period {
	case "day":
		start := now
		return start.AddDate(0, 0, -1), start
	case "week":
		start := now
		return start.AddDate(0, 0, -7), start
	case "month":
		start := now
		return start, start.AddDate(0, -1, 0)
	default:
		return time.Time{}, time.Time{}
	}
}

func (s *ReportService) GetBudgetReport(ctx context.Context, userID, budgetID string) (*models.BudgetReport, error) {
	budget, err := s.BudgetRepo.GetBudget(ctx, userID, budgetID)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	if budget == nil {
		return nil, errors.New("budget is not found")
	}
	now := time.Now()
	if now.Before(budget.EndDate) {
		return nil, errors.New("budget is not over yet")
	}
	txs, err := s.Transaction.GetTXByTimeFrame(ctx, userID, models.TimeFrame{budget.StartDate, budget.EndDate})
	if err != nil {
		log.Println(err)
		return nil, err
	}

	var totalSpent float64
	var failedDays, successfulDays int
	dailyExpenses := make(map[string]float64)
	itemCounts := make(map[string]int)
	var mostExpensive models.Transaction
	var mostFrequent string

	for _, tx := range txs {
		totalSpent += tx.Cost

		day := tx.Date.Format("2006-01-02")
		dailyExpenses[day] += tx.Cost

		if tx.Cost > mostExpensive.Cost {
			mostExpensive = tx
		}

		itemCounts[tx.Name]++
		if itemCounts[tx.Name] > itemCounts[mostFrequent] {
			mostFrequent = tx.Name
		}
	}
	for _, dailyTotal := range dailyExpenses {
		if dailyTotal > budget.DailyAmount {
			failedDays++
		} else {
			successfulDays++
		}
	}

	budgetExceeded := totalSpent > budget.Amount
	report := models.BudgetReport{
		TotalSpent:     totalSpent,
		BudgetExceeded: budgetExceeded,
		FailedDays:     failedDays,
		SuccessfulDays: successfulDays,
		MostExpensive:  mostExpensive,
		MostFrequent:   mostFrequent,
	}
	return &report, nil
}
