package handler

import (
	"budget/internal/models"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

func (h *Handler) CreateBudget(w http.ResponseWriter, r *http.Request) {
	amountStr := strings.TrimSpace(r.URL.Query().Get("amount"))
	amount, err := strconv.ParseFloat(amountStr, 64)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	budget := models.CreateBudget{
		UserID:  strings.TrimSpace(r.URL.Query().Get("userID")),
		Name:    strings.TrimSpace(r.URL.Query().Get("name")),
		Amount:  amount,
		EndTime: strings.TrimSpace(r.URL.Query().Get("end_date")),
	}
	if err := validate.Struct(budget); err != nil {
		http.Error(w, "Validation failed: "+err.Error(), http.StatusBadRequest)
		return
	}

	err = h.BudgetSRV.AddBudget(r.Context(), budget)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func (h *Handler) GetBudgetList(w http.ResponseWriter, r *http.Request) {
	userID := strings.TrimSpace(r.URL.Query().Get("userID"))
	if userID == "" {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	budgets, err := h.BudgetSRV.GetBudgetList(r.Context(), userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	response := map[string]interface{}{
		"budgets": budgets,
	}
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		log.Println("failed to encode JSON: %v", err)
		http.Error(w, "Failed to encode JSON response", http.StatusInternalServerError)
		return
	}
}
