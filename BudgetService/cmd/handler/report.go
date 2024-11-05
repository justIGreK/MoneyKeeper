package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
)

func (h *Handler) GetSummaryReport(w http.ResponseWriter, r *http.Request) {
	userID := strings.TrimSpace(r.URL.Query().Get("userID"))
	period := strings.TrimSpace(r.URL.Query().Get("period"))
	if userID == "" || period == "" {
		http.Error(w, "user_id and period are required", http.StatusBadRequest)
		return
	}
	report, err := h.ReportSRV.GetPeriodSummary(r.Context(), userID, period)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	response := map[string]interface{}{
		"report": report,
	}
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		log.Println("failed to encode JSON: %v", err)
		http.Error(w, "Failed to encode JSON response", http.StatusInternalServerError)
		return
	}
}

func (h *Handler) GetBudgetReport(w http.ResponseWriter, r *http.Request) {
	userID := strings.TrimSpace(r.URL.Query().Get("userID"))
	budgetID := strings.TrimSpace(r.URL.Query().Get("budgetID"))
	if userID == "" || budgetID == "" {
		http.Error(w, "user_id and budgetID are required", http.StatusBadRequest)
		return
	}
	report, err := h.ReportSRV.GetBudgetReport(r.Context(), userID, budgetID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	response := map[string]interface{}{
		"report": report,
	}
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		log.Println("failed to encode JSON: %v", err)
		http.Error(w, "Failed to encode JSON response", http.StatusInternalServerError)
		return
	}
}
