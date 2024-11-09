package service

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/justIGreK/MoneyKeeper/Bot/internal/models"
)

var (
	Dateformat string = "2006-01-02"
)
func (s *Service) AddBudget(args []string, chatID string) (string, error) {
	if len(args) < 3 {
		return "", errors.New("missing arguments")
	}
	id, err := s.GetUserID(chatID)
	if err != nil {
		return "", err
	}
	amount, err := strconv.ParseFloat(args[1], 32)
	if err != nil {
		return "", errors.New("invalid amount")
	}
	endDate, err := time.Parse(Dateformat, args[2])
	if err != nil {
		log.Println(err)
		return "", fmt.Errorf("invalid date or time format: %v", err)
	}
	createBudget := models.CreateBudget{
		UserID:  id,
		Name:    args[0],
		Amount:  float32(amount),
		EndTime: endDate.Format(Dateformat),
	}
	budgetId, err := s.budget.CreateBudget(s.ctx, createBudget)
	if err != nil {
		return "", err
	}
	message := fmt.Sprintf("budget is added: %v", budgetId)
	return message, nil
}

func (s *Service) GetBudgetList(args []string, chatID string) (string, error) {
	id, err := s.GetUserID(chatID)
	if err != nil {
		return "", err
	}
	budgets, err := s.budget.GetBudgetList(s.ctx, id)
	if err != nil {
		return "", err
	}
	var sb strings.Builder

	sb.WriteString("Budgets:\n\n")

	for _, budget := range budgets {
		sb.WriteString(fmt.Sprintf("*%s*\n", budget.Name))
		sb.WriteString(fmt.Sprintf("ID: %s\n", budget.ID))
		sb.WriteString(fmt.Sprintf("Amount: $%.2f\n", budget.Amount))
		sb.WriteString(fmt.Sprintf("Daily Limit: $%.2f\n", budget.DailyAmount))
		sb.WriteString(fmt.Sprintf("Start Date: %s\n", budget.StartDate))
		sb.WriteString(fmt.Sprintf("End Date: %s\n", budget.EndDate))
		sb.WriteString(fmt.Sprintf("Status: %v\n", budget.IsActive))
		sb.WriteString("\n")
	}
	return sb.String(), nil
}
