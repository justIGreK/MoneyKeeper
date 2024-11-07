package handler

import (
	"budget/internal/models"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"
)

func (h *Handler) CreateTransaction(w http.ResponseWriter, r *http.Request) {
	costStr := strings.TrimSpace(r.URL.Query().Get("cost"))
	cost, err := strconv.ParseFloat(costStr, 64)

	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	tx := models.CreateTransaction{
		Category: strings.TrimSpace(r.URL.Query().Get("category")),
		UserID:   strings.TrimSpace(r.URL.Query().Get("userID")),
		Name:     strings.TrimSpace(r.URL.Query().Get("name")),
		Cost:     cost,
	}
	if err := validate.Struct(tx); err != nil {
		http.Error(w, "Validation failed: "+err.Error(), http.StatusBadRequest)
		return
	}
	notifications, err := h.TxSRV.AddTransaction(r.Context(), tx)
	if err != nil {
		http.Error(w, "operation is declined", http.StatusBadRequest)
		return
	}
	if len(notifications) == 0 {
		w.WriteHeader(http.StatusCreated)
		return
	}
	response := map[string]interface{}{
		"notifications": notifications,
	}
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		log.Println("failed to encode JSON: %v", err)
		http.Error(w, "Failed to encode JSON response", http.StatusInternalServerError)
		return
	}
}

func (h *Handler) GetTransaction(w http.ResponseWriter, r *http.Request) {
	userID := strings.TrimSpace(r.URL.Query().Get("userID"))
	txID := strings.TrimSpace(r.URL.Query().Get("txID"))

	tx, err := h.TxSRV.GetTransaction(r.Context(), txID, userID)
	if err != nil {
		http.Error(w, "operation is declined", http.StatusBadRequest)
		return
	}
	response := map[string]interface{}{
		"transaction": tx,
	}
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		log.Println("failed to encode JSON: %v", err)
		http.Error(w, "Failed to encode JSON response", http.StatusInternalServerError)
		return
	}
}

func (h *Handler) GetTransactionList(w http.ResponseWriter, r *http.Request) {
	userID := strings.TrimSpace(r.URL.Query().Get("userID"))

	txs, err := h.TxSRV.GetAllTransactions(r.Context(), userID)
	if err != nil {
		http.Error(w, "operation is declined", http.StatusBadRequest)
		return
	}
	response := map[string]interface{}{
		"transactions": txs,
	}
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		log.Println("failed to encode JSON: %v", err)
		http.Error(w, "Failed to encode JSON response", http.StatusInternalServerError)
		return
	}
}

func (h *Handler) GetTXByTimeFrame(w http.ResponseWriter, r *http.Request) {
	userID := strings.TrimSpace(r.URL.Query().Get("userID"))
	timeFrame := models.CreateTimeFrame{
		StartDate: strings.TrimSpace(r.URL.Query().Get("start")),
		EndDate:   strings.TrimSpace(r.URL.Query().Get("end")),
	}
	txs, err := h.TxSRV.GetTXByTimeFrame(r.Context(), userID, timeFrame)
	if err != nil {
		http.Error(w, "operation is declined", http.StatusBadRequest)
		return
	}
	response := map[string]interface{}{
		"transactions": txs,
	}
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		log.Println("failed to encode JSON: %v", err)
		http.Error(w, "Failed to encode JSON response", http.StatusInternalServerError)
		return
	}
}
