package service

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/justIGreK/MoneyKeeper/Bot/internal/models"
	"github.com/justIGreK/MoneyKeeper/Bot/pkg/util"
)

func (s *Service) AddBudget(args []string, chatID string) (string, error) {
	if len(args) < 3 {
		return "", errors.New("missing arguments")
	}
	id, err := s.GetUserID(chatID)
	if err != nil {
		return "", err
	}
	name := args[0]
	amount, err := strconv.ParseFloat(args[1], 32)
	if err != nil {
		return "", errors.New("invalid amount")
	}
	options := make(map[string]interface{})
	periodSpecified := false
	endSpecified := false
	startSpecified := false
	dateSpecified := false
	for _, part := range args[2:] {
		key, value, ok := util.ParseKeyValue(part)
		if !ok {
			return "", fmt.Errorf("incorrect parameter format: %s", part)
		}
		switch key {
		case "period":
			if periodSpecified {
				return "", fmt.Errorf("period is given several times")
			}
			options["period"] = value
			periodSpecified = true
		case "start":
			if startSpecified {
				return "", fmt.Errorf("period is given several times")
			}
			startDate, err := util.ParseDate(value)
			if err != nil {
				return "", fmt.Errorf("incorrect start format: %s", value)
			}
			options["start"] = startDate
			startSpecified = true
		case "end":
			if endSpecified {
				return "", fmt.Errorf("period is given several times")
			}
			endDate, err := util.ParseDate(value)
			if err != nil {
				return "", fmt.Errorf("incorrect end format: %s", value)
			}
			options["end"] = endDate
			endSpecified = true
		default:
			return "", fmt.Errorf("unknown parameter: %s", key)
		}
	}
	if endSpecified && startSpecified {
		dateSpecified = true
	}
	if periodSpecified && dateSpecified {
		return "", fmt.Errorf("It is not possible to specify both the period and specific dates")
	}
	if !periodSpecified && !dateSpecified {
		return "", fmt.Errorf("You must specify either the period or the start and end dates")
	}
	createBudget := models.CreateBudget{
		UserID: id,
		Name:   name,
		Amount: float32(amount),
	}
	if value, exist := options["period"]; exist {
		createBudget.Period = value.(string)
	} else {
		createBudget.Period = ""
	}
	if value, exist := options["start"]; exist {
		start := value.(time.Time)
		createBudget.Start = start.Format(Dateformat)  
	}
	if value, exist := options["end"]; exist {
		end := value.(time.Time)
		createBudget.End = end.Format(Dateformat) 
	}
	budgetId, err := s.budget.AddBudget(s.ctx, createBudget)
	if err != nil {
		return "", err
	}
	message := fmt.Sprintf("budget is added: %v", budgetId)
	return message, nil
}

func (s *Service) AddCategory(args []string, chatID string) (string, error) {
	if len(args) != 3 {
		return "", errors.New("invalid response")
	}
	id, err := s.GetUserID(chatID)
	if err != nil {
		return "", err
	}
	limit, err := strconv.ParseFloat(args[2], 32)
	if err != nil {
		return "", errors.New("invalid amount")
	}
	budget, err := s.budget.AddCategory(s.ctx, models.CreateCategory{
		UserID:   id,
		BudgetID: args[0],
		Name:     args[1],
		Limit:    float32(limit),
	})
	if err != nil {
		return "", err
	}
	message := s.PrepareBudgets([]models.Budget{*budget})

	return "Updated budget:" + message, nil
}

func (s *Service) GetBudget(args []string, chatID string) (string, error) {
	if len(args) != 1 {
		return "", errors.New("invalid response")
	}
	id, err := s.GetUserID(chatID)
	if err != nil {
		return "", err
	}

	budget, err := s.budget.GetBudget(s.ctx, id, args[0])
	if err != nil {
		return "", err
	}
	message := s.PrepareBudgets([]models.Budget{*budget})

	return "Budget:" + message, nil
}
func (s *Service) DeleteBudget(args []string, chatID string) (string, error) {
	if len(args) != 1 {
		return "", errors.New("invalid response")
	}
	id, err := s.GetUserID(chatID)
	if err != nil {
		return "", err
	}
	err = s.budget.DeleteBudget(s.ctx, id, args[0])
	if err != nil {
		return "", err
	}
	return "Budget is deleted", nil
}

func (s *Service) DeleteCategory(args []string, chatID string) (string, error) {
	if len(args) != 2 {
		return "", errors.New("invalid response")
	}
	id, err := s.GetUserID(chatID)
	if err != nil {
		return "", err
	}
	err = s.budget.DeleteCategory(s.ctx, id, args[0], args[1])
	if err != nil {
		return "", err
	}
	return "Category is deleted", nil
}

func (s *Service) PrepareBudgets(budgets []models.Budget) string {
	var sb strings.Builder
	sb.WriteString("\n")
	for _, budget := range budgets {
		sb.WriteString(fmt.Sprintf("*%s*\n", budget.Name))
		sb.WriteString(fmt.Sprintf("ID: %s\n", budget.ID))
		sb.WriteString(fmt.Sprintf("Limit: $%.2f\n", budget.Limit))
		sb.WriteString(fmt.Sprintf("Period: %v - %v\n", budget.Start, budget.End))
		if len(budget.Category) != 0 {
			sb.WriteString("Categories:\n")
			for _, category := range budget.Category {
				sb.WriteString("---------------------------\n")
				sb.WriteString(fmt.Sprintf("ID: %s\n", category.ID))
				sb.WriteString(fmt.Sprintf("Category: %s\n", category.Name))
				sb.WriteString(fmt.Sprintf("Limit: %.2f\n", category.Limit))
			}
			sb.WriteString("---------------------------\n")
		}
		sb.WriteString("\n")
	}
	return sb.String()
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

	message := s.PrepareBudgets(budgets)

	return "Budgets:" + message, nil
}

func (s *Service) UpdateBudget(args []string, chatID string) (string, error) {
	if len(args) < 2 || len(args) > 6 {
		return "", errors.New("incorrect input format")
	}
	id, err := s.GetUserID(chatID)
	if err != nil {
		return "", err
	}
	budgetID := args[0]
	options := make(map[string]interface{})
	nameSpecified := false
	limitSpecified := false
	startSpecified := false
	endSpecified := false
	for _, part := range args[1:] {
		key, value, ok := util.ParseKeyValue(part)
		if !ok {
			return "", fmt.Errorf("incorrect parameter format: %s", part)
		}
		switch key {
		case "name":
			if nameSpecified {
				return "", fmt.Errorf("name is given several times")
			}
			options["name"] = value
			nameSpecified = true
		case "limit":
			if limitSpecified {
				return "", fmt.Errorf("limit is given several times")
			}
			limit, err := strconv.ParseFloat(value, 64)
			if err != nil {
				return "", errors.New("invalid amount")
			}
			options["limit"] = limit
			limitSpecified = true
		case "start":
			if startSpecified {
				return "", fmt.Errorf("period is given several times")
			}
			startDate, err := util.ParseDate(value)
			if err != nil {
				return "", fmt.Errorf("incorrect start format: %s", value)
			}
			options["start"] = startDate
			startSpecified = true
		case "end":
			if endSpecified {
				return "", fmt.Errorf("period is given several times")
			}
			endDate, err := util.ParseDate(value)
			if err != nil {
				return "", fmt.Errorf("incorrect end format: %s", value)
			}
			options["end"] = endDate
			endSpecified = true
		default:
			return "", fmt.Errorf("unknown parameter: %s", key)
		}
	}
	updateBudget := s.prepareBudgetForUpdate(options)
	updateBudget.UserID, updateBudget.BudgetID = id, budgetID

	budget, err := s.budget.UpdateBudget(s.ctx, updateBudget)
	if err != nil {
		return "", err
	}
	message := s.PrepareBudgets([]models.Budget{*budget})

	return "Updated budget:" + message, nil
}

func (s *Service) UpdateCategory(args []string, chatID string) (string, error) {
	if len(args) < 3 || len(args) > 4 {
		return "", errors.New("incorrect input format")
	}
	id, err := s.GetUserID(chatID)
	if err != nil {
		return "", err
	}
	budgetID := args[0]
	categoryID := args[1]
	options := make(map[string]interface{})
	nameSpecified := false
	limitSpecified := false
	for _, part := range args[2:] {
		key, value, ok := util.ParseKeyValue(part)
		if !ok {
			return "", fmt.Errorf("incorrect parameter format: %s", part)
		}
		switch key {
		case "name":
			if nameSpecified {
				return "", fmt.Errorf("name is given several times")
			}
			options["name"] = value
			nameSpecified = true
		case "limit":
			if limitSpecified {
				return "", fmt.Errorf("limit is given several times")
			}
			limit, err := strconv.ParseFloat(value, 64)
			if err != nil {
				return "", errors.New("invalid amount")
			}
			options["limit"] = limit
			limitSpecified = true
		default:
			return "", fmt.Errorf("unknown parameter: %s", key)
		}
	}
	UpdateCategory := s.prepareCategoryForUpdate(options)
	UpdateCategory.UserID, UpdateCategory.BudgetID, UpdateCategory.CategoryID = id, budgetID, categoryID

	budget, err := s.budget.UpdateCategory(s.ctx, UpdateCategory)
	if err != nil {
		return "", err
	}
	message := s.PrepareBudgets([]models.Budget{*budget})

	return "Updated budget:" + message, nil
}
func (s *Service) prepareCategoryForUpdate(updates map[string]interface{}) models.UpdateCategory {
	var categoryUpdates models.UpdateCategory
	for field, value := range updates {
		switch field {
		case "name":
			name := value.(string)
			categoryUpdates.Name = &name
		case "limit":
			limit := value.(float64)
			categoryUpdates.Limit = &limit
		}
	}
	return categoryUpdates
}

func (s *Service) prepareBudgetForUpdate(updates map[string]interface{}) models.UpdateBudget {
	var budgetUpdates models.UpdateBudget
	for field, value := range updates {
		switch field {
		case "name":
			name := value.(string)
			budgetUpdates.Name = &name
		case "limit":
			limit := value.(float64)
			budgetUpdates.Limit = &limit
		case "start":
			start := value.(time.Time)
			startstr := start.Format(Dateformat)
			budgetUpdates.Start = &startstr
		case "end":
			end := value.(time.Time)
			endstr := end.Format(Dateformat)
			budgetUpdates.End = &endstr
		}
	}
	return budgetUpdates
}
