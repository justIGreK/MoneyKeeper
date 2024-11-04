package handler

import (
	"budget/internal/models"
	"context"

	"github.com/go-chi/chi/v5"
)

type UserService interface {
	CreateUser(ctx context.Context, user *models.User) (string, error)
	GetUser(ctx context.Context, userID string) (*models.User, error)
}

type BudgetService interface{
	AddBudget(ctx context.Context, budget models.CreateBudget) error
	GetBudgetList(ctx context.Context, userID string) ([]models.Budget, error)
}

type TransactionService interface{
	AddTransaction(ctx context.Context, transaction models.CreateTransaction) ([]string, error) 
	GetTransaction(ctx context.Context, transactionID, userID string) (*models.Transaction, error) 
	GetAllTransactions(ctx context.Context, userID string) ([]models.Transaction, error) 
	GetTXByTimeFrame(ctx context.Context, userID string, timeframe models.CreateTimeFrame) ([]models.Transaction, error) 
}

type Handler struct {
	UserSRV UserService
	BudgetSRV BudgetService
	TxSRV TransactionService
}

func NewHandler(user UserService, budget BudgetService, tx TransactionService) *Handler {
	return &Handler{UserSRV: user, BudgetSRV: budget, TxSRV: tx}
}

func (h *Handler) InitRoutes() *chi.Mux {
	r := chi.NewRouter()
	r.Route("/user", func(r chi.Router) {
		r.Post("/create", h.CreateUser)
		r.Get("/getUser", h.GetUser)
	})
	r.Route("/budget", func(r chi.Router) {
		r.Post("/add", h.CreateBudget)
		r.Get("/getlist", h.GetBudgetList)
	})
	r.Route("/tx", func(r chi.Router) {
		r.Post("/add", h.CreateTransaction)
		r.Get("/get", h.GetTransaction)
		r.Get("/getlist", h.GetTransactionList)
		r.Get("/getbytf", h.GetTXByTimeFrame)
	})
	return r
}
