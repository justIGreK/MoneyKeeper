package client

import (
	"context"

	"github.com/justIGreK/MoneyKeeper/Bot/internal/models"
	report "github.com/justIGreK/MoneyKeeper/BudgetService/pkg/go/report"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type ReportClient struct {
	client report.ReportServiceClient
}

func NewReportClient(serviceAddress string) (*ReportClient, error) {
	conn, err := grpc.NewClient(serviceAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	return &ReportClient{
		client: report.NewReportServiceClient(conn),
	}, nil
}

func (rc *ReportClient) GetSummaryReport(ctx context.Context, userID, period string) (*models.ReportResponse, error) {
	req := &report.GetSummaryReportRequest{
		UserId: userID,
		Period: period,
	}
	res, err := rc.client.GetSummaryReport(ctx, req)
	if err != nil {
		return nil, err
	}
	report := rc.convertToReportResponse(res.Report)
	return report, nil
}

func (rc *ReportClient) convertToReportResponse(report *report.ReportResponse) *models.ReportResponse {
	categories := make([]*models.CategoryReport, len(report.Categories))
	for i, b := range report.Categories {
		categories[i] = &models.CategoryReport{
			Name:  b.Name,
			Total: float32(b.Total),
			Count: b.Count,
		}
	}
	return &models.ReportResponse{UserID: report.UserId, Period: report.Period,
		TotalSpent: float32(report.TotalSpent), TransactionCount: int32(report.TransactionCount),
		Categories: categories}
}

func (rc *ReportClient) GetBudgetReport(ctx context.Context, userID, budgetID string) (*models.BudgetReport, error) {
	req := &report.GetBudgetReportRequest{
		UserId:   userID,
		BudgetId: budgetID,
	}
	res, err := rc.client.GetBudgetReport(ctx, req)
	if err != nil {
		return nil, err
	}
	budgetReport := &models.BudgetReport{
		TotalSpent:     res.Report.TotalSpent,
		BudgetExceeded: res.Report.BudgetExceeded,
		FailedDays:     res.Report.FailedDays,
		SuccessfulDays: res.Report.SuccessfulDays,
		MostFrequent:   res.Report.MostFrequent,
		MostExpensive: models.Transaction{
			ID:       res.Report.MostExpensive.Id,
			Category: res.Report.MostExpensive.Category,
			Name:     res.Report.MostExpensive.Name,
			Cost:     res.Report.MostExpensive.Cost,
			Date:     res.Report.MostExpensive.Date,
		},
	}
	return budgetReport, nil
}
