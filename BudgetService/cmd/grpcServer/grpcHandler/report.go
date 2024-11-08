package grpchandler

import (
	"github.com/justIGreK/MoneyKeeper/BudgetService/internal/models"
	reportProto "github.com/justIGreK/MoneyKeeper/BudgetService/pkg/go/report"
	"context"
)

type ReportServiceServer struct {
	reportProto.UnimplementedReportServiceServer
	ReportSRV ReportService
}

type ReportService interface {
	GetPeriodSummary(ctx context.Context, userID, period string) (*models.ReportResponse, error)
	GetBudgetReport(ctx context.Context, userID, budgetID string) (*models.BudgetReport, error)
}

func (s *ReportServiceServer) GetSummaryReport(ctx context.Context, req *reportProto.GetSummaryReportRequest) (*reportProto.GetSummaryReportResponse, error) {

	report, err := s.ReportSRV.GetPeriodSummary(ctx, req.UserId, req.Period)
	if err != nil {
		return nil, err
	}

	protoReport := convertToProtoReportResponse(report)
	return &reportProto.GetSummaryReportResponse{
		Report: protoReport,
	}, nil

}

func convertToProtoReportResponse(report *models.ReportResponse) *reportProto.ReportResponse {
	categories := make([]*reportProto.CategoryReport, len(report.Categories))
	for i, b := range report.Categories {
		categories[i] = &reportProto.CategoryReport{
			Name:  b.Name,
			Total: float32(b.Total),
			Count: int32(b.Count),
		}
	}

	return &reportProto.ReportResponse{UserId: report.UserID, Period: report.Period,
		TotalSpent: float32(report.TotalSpent), TransactionCount: int32(report.TransactionCount),
		Categories: categories}
}

func (s *ReportServiceServer) GetBudgetReport(ctx context.Context, req *reportProto.GetBudgetReportRequest) (*reportProto.GetBudgetReportResponse, error) {
	budgetReport, err := s.ReportSRV.GetBudgetReport(ctx, req.UserId, req.BudgetId)
	if err != nil {
		return nil, err
	}

	protoReport := convertToProtoBudgetReport(budgetReport)
	return &reportProto.GetBudgetReportResponse{
		Report: protoReport,
	}, nil

}

func convertToProtoBudgetReport(report *models.BudgetReport) *reportProto.BudgetReport {
	return &reportProto.BudgetReport{
		TotalSpent:     float32(report.TotalSpent),
		BudgetExceeded: report.BudgetExceeded,
		FailedDays:     int32(report.FailedDays),
		SuccessfulDays: int32(report.SuccessfulDays),
		MostFrequent:   report.MostFrequent,
		MostExpensive: &reportProto.Transaction{
			Id:       report.MostExpensive.ID,
			UserId:   report.MostExpensive.UserID,
			Category: report.MostExpensive.Category,
			Name:     report.MostExpensive.Name,
			Cost:     float32(report.MostExpensive.Cost),
			Date:     report.MostExpensive.Date.Format(Dateformat),
		}}
}
