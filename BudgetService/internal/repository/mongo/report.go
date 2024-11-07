package mongorep

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
)

type ReportRepo struct {
	ReportCollection *mongo.Collection
	TxCollection     *mongo.Collection
}

func NewReportRepository(db *mongo.Client) *ReportRepo {
	return &ReportRepo{
		ReportCollection: db.Database(dbname).Collection(reportCollection),
		TxCollection:     db.Database(dbname).Collection(transactionCollection),
	}
}

func (r *ReportRepo) AddReport() {

}

func (r *ReportRepo) BudgetSummary(ctx context.Context, userID string, start, end time.Time) {

}
