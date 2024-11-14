package service

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/justIGreK/MoneyKeeper/Bot/internal/models"
	"github.com/justIGreK/MoneyKeeper/Bot/pkg/util"
)

func (s *Service) GetSummary(args []string, chatID string) (string, error) {
	if len(args) < 1 || len(args) > 3 {
		return "", errors.New("incorrect input format")
	}
	id, err := s.GetUserID(chatID)
	if err != nil {
		return "", err
	}
	options := make(map[string]interface{})
	periodSpecified := false
	endSpecified := false
	startSpecified := false
	dateSpecified := false
	for _, part := range args {
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
	if endSpecified || startSpecified {
		dateSpecified = true
	}
	if periodSpecified && dateSpecified {
		return "", fmt.Errorf("It is not possible to specify both the period and specific dates")
	}
	if !periodSpecified && !dateSpecified {
		return "", fmt.Errorf("You must specify either the period or the start and end dates")
	}
	dateSum := s.getDates(options)
	dateSum.UserID = id
	report, err := s.report.GetSummaryReport(s.ctx, dateSum)
	if err != nil {
		return "", err
	}
	summary := s.prepareSummaryReport(*report)
	return summary, nil
}

func (s *Service) getDates(dates map[string]interface{}) models.DatesSummary {
	dateSummary := models.DatesSummary{}
	for field, value := range dates {
		switch field {
		case "period":
			period := value.(string)
			dateSummary.Period = period
		case "start":
			start := value.(time.Time)
			startstr := start.Format(Dateformat)
			dateSummary.Start = startstr
		case "end":
			end := value.(time.Time)
			endstr := end.Format(Dateformat)
			dateSummary.End = endstr
		}
	}
	return dateSummary
}
func (s *Service) prepareSummaryReport(report models.ReportResponse) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Summary for %s\n", report.Period))
	sb.WriteString(fmt.Sprintf("Total expenses: %.2f\n", report.TotalSpent))
	sb.WriteString(fmt.Sprintf("Total number of transactions: %d\n", report.TransactionCount))
	if len(report.Categories) != 0 {
		sb.WriteString("Transactions by categories:\n")
		for _, category := range report.Categories {
			sb.WriteString("---------------------------\n")
			sb.WriteString(fmt.Sprintf("Category: %s\n", category.Name))
			sb.WriteString(fmt.Sprintf("Total spent: %.2f\n", category.Total))
			sb.WriteString(fmt.Sprintf("Count: %d\n", category.Count))
		}
		sb.WriteString("---------------------------\n")
	}
	return sb.String()
}

func (s *Service) GetBudgetReport(args []string, chatID string) (string, error) {
	if len(args) != 1 {
		return "", errors.New("incorrect input format")
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
	sb.WriteString(fmt.Sprintf("Budget *%v* report:\n ", report.BudgetName))
	sb.WriteString(fmt.Sprintf("Period: %v\n", report.Period))
	if report.LeftDays != 0 {
		sb.WriteString(fmt.Sprintf("%.2f days left until the end of the budget period\n", report.LeftDays))
	}
	sb.WriteString(fmt.Sprintf("Total expenses %.2f / %.2f\n", report.TotalSpent, report.Limit))
	sb.WriteString(fmt.Sprintf("Total count of transactions %d\n", report.TransactionCount))
	if len(report.RequiredCategories) != 0 {
		sb.WriteString(fmt.Sprintf("\nResult of categories with limit:\n"))
		sb.WriteString("---------------------------\n")
		for _, categ := range report.RequiredCategories {
			
				sb.WriteString(fmt.Sprintf("%v: %0.2f / %0.2f\n", categ.Name, categ.Total, categ.Limit))
				sb.WriteString(fmt.Sprintf("Total count: %d\n", categ.Count))
				sb.WriteString("---------------------------\n")
			
		}
	}
	if len(report.Categories) != 0 {
		sb.WriteString(fmt.Sprintf("\n Other categories:\n"))
		sb.WriteString("---------------------------\n")
		for _, categ := range report.Categories {
			sb.WriteString(fmt.Sprintf("%v: %0.2f\n", categ.Name, categ.Total))
			sb.WriteString(fmt.Sprintf("Total count: %d\n", categ.Count))
			sb.WriteString("---------------------------\n")
		}
	}
	return sb.String()
}
