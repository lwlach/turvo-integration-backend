package load

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/lwlach/turvo-integration-backend/internal/models"
	"github.com/lwlach/turvo-integration-backend/internal/service/load"
)

type Handler struct {
	service *load.Service
}

func NewHandler(service *load.Service) *Handler {
	return &Handler{
		service: service,
	}
}

// RegisterRoutes registers the load routes with the chi router
func (h *Handler) RegisterRoutes(r chi.Router) {
	r.Get("/loads", h.GetLoads)
	r.Post("/loads", h.CreateLoad)
}

// GetLoads handles GET /loads - returns filtered and paginated loads from Turvo
func (h *Handler) GetLoads(w http.ResponseWriter, r *http.Request) {
	filters := parseFilters(r)

	response, err := h.service.GetLoads(filters)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
		return
	}
}

// parseFilters parses query parameters into LoadFilters
func parseFilters(r *http.Request) models.LoadFilters {
	filters := models.LoadFilters{
		Page:  1,
		Limit: 20,
	}

	// Parse status filter
	if status := r.URL.Query().Get("status"); status != "" {
		filters.Status = status
	}

	// Parse customerId filter
	if customerID := r.URL.Query().Get("customerId"); customerID != "" {
		filters.CustomerID = customerID
	}

	// Parse pickupDateSearchFrom (start of day in UTC)
	if dateFromStr := r.URL.Query().Get("pickupDateSearchFrom"); dateFromStr != "" {
		if dateFrom, err := time.Parse(time.RFC3339, dateFromStr); err == nil {
			// Set to start of day in UTC
			startOfDay := time.Date(dateFrom.Year(), dateFrom.Month(), dateFrom.Day(), 0, 0, 0, 0, time.UTC)
			filters.PickupDateFrom = &startOfDay
		}
	}

	// Parse pickupDateSearchTo (end of day in UTC)
	if dateToStr := r.URL.Query().Get("pickupDateSearchTo"); dateToStr != "" {
		if dateTo, err := time.Parse(time.RFC3339, dateToStr); err == nil {
			// Set to end of day in UTC
			endOfDay := time.Date(dateTo.Year(), dateTo.Month(), dateTo.Day(), 23, 59, 59, 999999999, time.UTC)
			filters.PickupDateTo = &endOfDay
		}
	}

	// Parse page (default: 1, min: 1)
	if pageStr := r.URL.Query().Get("page"); pageStr != "" {
		if page, err := strconv.Atoi(pageStr); err == nil && page >= 1 {
			filters.Page = page
		}
	}

	// Parse limit (default: 20, min: 1, max: 100)
	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		if limit, err := strconv.Atoi(limitStr); err == nil {
			if limit < 1 {
				limit = 1
			}
			if limit > 100 {
				limit = 100
			}
			filters.Limit = limit
		}
	}

	// Parse includeDetails flag
	if includeDetails := r.URL.Query().Get("includeDetails"); includeDetails != "" {
		if includeDetails == "true" || includeDetails == "1" {
			filters.IncludeDetails = true
		}
	}

	return filters
}

// CreateLoad handles POST /loads - creates a new load in Turvo
func (h *Handler) CreateLoad(w http.ResponseWriter, r *http.Request) {
	var load models.Load

	if err := json.NewDecoder(r.Body).Decode(&load); err != nil {
		http.Error(w, "invalid request body: "+err.Error(), http.StatusBadRequest)
		return
	}

	response, err := h.service.CreateLoad(&load)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
		return
	}
}
