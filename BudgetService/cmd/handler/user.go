package handler

import (
	"github.com/justIGreK/MoneyKeeper/BudgetService/internal/models"
	"encoding/json"
	"log"
	"net/http"
	"strings"
)

func (h *Handler) CreateUser(w http.ResponseWriter, r *http.Request) {
	user := models.User{
		Name: strings.TrimSpace(r.URL.Query().Get("name")),
	}
	if user.Name == "" {
		http.Error(w, "invalid name", http.StatusBadRequest)
		return
	}
	id, err := h.UserSRV.CreateUser(r.Context(), &user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(id)
	if err != nil {
		log.Println("failed to encode JSON: %v", err)
		http.Error(w, "Failed to encode JSON response", http.StatusInternalServerError)
		return
	}
}

func (h *Handler) GetUser(w http.ResponseWriter, r *http.Request) {
	userID := strings.TrimSpace(r.URL.Query().Get("userID"))
	user, err := h.UserSRV.GetUser(r.Context(), userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(user)
	if err != nil {
		log.Println("failed to encode JSON: %v", err)
		http.Error(w, "Failed to encode JSON response", http.StatusInternalServerError)
		return
	}
}
