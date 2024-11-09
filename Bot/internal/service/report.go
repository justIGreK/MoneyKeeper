package service

import (
	"errors"
	"fmt"
	"strings"

	"github.com/justIGreK/MoneyKeeper/Bot/internal/models"
)

func (s *Service) GetSummary(args []string, chatID string) (string, error) {
	if len(args) < 1 {
		return "", errors.New("missing arguments")
	}
	id, err := s.GetUserID(chatID)
	if err != nil {
		return "", err
	}
	report, err := s.report.GetSummaryReport(s.ctx, id, args[0])
	if err != nil {
		return "", err
	}
	summary := s.prepareSummaryReport(*report)
	return summary, nil
}

func (s *Service) prepareSummaryReport(report models.ReportResponse) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Summary for %s\n", report.Period))
	sb.WriteString(fmt.Sprintf("Total expenses: %.2f\n", report.TotalSpent))
	sb.WriteString(fmt.Sprintf("Total number of transactions: %d\n", report.TransactionCount))
	if len(report.Categories) != 0 {
		sb.WriteString("Transactions by categories:\n")
		sb.WriteString("---------------------------\n")
		for _, category := range report.Categories {
			sb.WriteString(fmt.Sprintf("Category: %s\n", category.Name))
			sb.WriteString(fmt.Sprintf("Total spent: %.2f\n", category.Total))
			sb.WriteString(fmt.Sprintf("Count: %d\n", category.Count))
			sb.WriteString("---------------------------\n")
		}
	}
	return sb.String()
}

func (s *Service) GetBudgetReport(args []string, chatID string) (string, error) {
	if len(args) < 1 {
		return "", errors.New("missing arguments")
	}
	id, err := s.GetUserID(chatID)
	if err != nil {
		return "", err
	}
	report, err := s.report.GetBudgetReport(s.ctx, id, args[0])
	if err != nil {
		return "", err
	}
	summary := s.prepareBudgetReport(*report)
	return summary, nil
}

func (s *Service) prepareBudgetReport(report models.BudgetReport) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Result of Budget: "))
	if report.BudgetExceeded {
		sb.WriteString("Failed\n")
	} else {
		sb.WriteString("Success\n")
	}
	sb.WriteString(fmt.Sprintf("Total expenses: %.2f\n", report.TotalSpent))
	sb.WriteString(fmt.Sprintf("Count of successful days: %d\n", report.SuccessfulDays))
	sb.WriteString(fmt.Sprintf("Count of failed days: %d\n", report.FailedDays))
	sb.WriteString(fmt.Sprintf("Name of most frequent transaction: %s\n", report.MostFrequent))
	sb.WriteString("Most expensive transaction:")
	sb.WriteString(s.PrepareTransactions([]models.Transaction{report.MostExpensive}))

	return sb.String()
}
