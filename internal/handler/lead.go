package handler

import (
	"encoding/json"
	"kibit/internal/model"
	"kibit/internal/service"
	"net/http"

	"github.com/go-playground/validator/v10"
)

type LeadHandler struct {
	svc *service.LeadService
	val *validator.Validate
}

func NewLeadHandler(s *service.LeadService) *LeadHandler { return &LeadHandler{s, validator.New()} }
func (h *LeadHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req model.CreateLeadRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}
	if err := h.val.Struct(req); err != nil {
		http.Error(w, "validation failed", http.StatusUnprocessableEntity)
		return
	}
	if err := h.svc.Create(r.Context(), &req); err != nil {
		http.Error(w, "server error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"status": "success"})
}
