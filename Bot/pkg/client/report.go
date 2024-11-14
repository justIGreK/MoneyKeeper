package client

import (
	"context"

	report "github.com/justIGreK/MoneyKeeper-Report/pkg/go/report"
	"github.com/justIGreK/MoneyKeeper/Bot/internal/models"
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

func (rc *ReportClient) GetSummaryReport(ctx context.Context, periodSum models.DatesSummary) (*models.ReportResponse, error) {
	req := &report.GetSummaryReportRequest{
		UserId: periodSum.UserID,
		Period: periodSum.Period,
		Start: periodSum.Start,
		End: periodSum.End,
	}
	res, err := rc.client.GetSummaryReport(ctx, req)
	if err != nil {
		return nil, err
	}
	return &models.ReportResponse{
		Period:           res.Report.Period,
		TotalSpent:       float32(res.Report.TotalSpent),
		TransactionCount: res.Report.TransactionCount,
		Categories:       rc.convertToCategoryReport(res.Report.Categories),
	}, nil
}

func (rc *ReportClient) convertToCategoryReport(categs []*report.CategoryReport) []models.CategoryReport {
	categories := make([]models.CategoryReport, len(categs))
	for i, b := range categs {
		categories[i] = models.CategoryReport{
			Name:  b.Name,
			Total: float32(b.Total),
			Count: b.Count,
		}
	}
	return categories
}
func (rc *ReportClient) convertToReqCategoryReport(categs []*report.RequiredCategoryReport) []models.RequiredCategoryReport {
	categories := make([]models.RequiredCategoryReport, len(categs))
	for i, b := range categs {
		categories[i] = models.RequiredCategoryReport{
			Name:  b.Name,
			Total: float32(b.Total),
			Limit: b.Limit,
			Count: b.Count,
		}
	}
	return categories
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
		BudgetName:         res.Report.BudgetName,
		Period:             res.Report.Period,
		LeftDays:           res.Report.LeftDays,
		Limit:              res.Report.Limit,
		TotalSpent:         res.Report.TotalSpent,
		TransactionCount:   int(res.Report.TransactionCount),
		RequiredCategories: rc.convertToReqCategoryReport(res.Report.ReqCategories),
		Categories:         rc.convertToCategoryReport(res.Report.Categories),
	}

	return budgetReport, nil
}
