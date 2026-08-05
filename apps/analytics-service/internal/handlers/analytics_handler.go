package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/wgmasvix-hue/ChengetAi-learn-/apps/analytics-service/internal/services"
	"github.com/wgmasvix-hue/ChengetAi-learn-/platform/pkg/analytics"
	"go.uber.org/zap"
)

// AnalyticsHandler handles analytics endpoints
type AnalyticsHandler struct {
	analyticsService *services.AnalyticsService
	logger           *zap.SugaredLogger
}

// NewAnalyticsHandler creates a new analytics handler
func NewAnalyticsHandler(analyticsService *services.AnalyticsService, logger *zap.SugaredLogger) *AnalyticsHandler {
	return &AnalyticsHandler{
		analyticsService: analyticsService,
		logger:           logger,
	}
}

// TrackEvent handles POST /events/track
func (h *AnalyticsHandler) TrackEvent(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var event analytics.Event
	if err := json.NewDecoder(r.Body).Decode(&event); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid request body"})
		return
	}

	if err := h.analyticsService.TrackEvent(r.Context(), &event); err != nil {
		h.logger.Errorw("Error tracking event", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"status": "event tracked"})
}

// TrackBatchEvents handles POST /events/batch
func (h *AnalyticsHandler) TrackBatchEvents(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var events []analytics.Event
	if err := json.NewDecoder(r.Body).Decode(&events); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid request body"})
		return
	}

	if err := h.analyticsService.TrackBatchEvents(r.Context(), events); err != nil {
		h.logger.Errorw("Error tracking batch events", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status": "events tracked",
		"count":  len(events),
	})
}

// GetPlatformMetrics handles GET /metrics/platform
func (h *AnalyticsHandler) GetPlatformMetrics(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	days := 30
	if d := r.URL.Query().Get("days"); d != "" {
		if parsedDays, err := strconv.Atoi(d); err == nil {
			days = parsedDays
		}
	}

	metrics, err := h.analyticsService.GetPlatformMetrics(r.Context(), days)
	if err != nil {
		h.logger.Errorw("Error getting platform metrics", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(metrics)
}

// GetUserMetrics handles GET /metrics/user/{userId}
func (h *AnalyticsHandler) GetUserMetrics(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	userID := r.PathValue("userId")
	if userID == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "User ID required"})
		return
	}

	metrics, err := h.analyticsService.GetUserMetrics(r.Context(), userID)
	if err != nil {
		h.logger.Errorw("Error getting user metrics", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(metrics)
}

// GetResourceMetrics handles GET /metrics/resource/{resourceId}
func (h *AnalyticsHandler) GetResourceMetrics(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	resourceID := r.PathValue("resourceId")
	if resourceID == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Resource ID required"})
		return
	}

	metrics, err := h.analyticsService.GetResourceMetrics(r.Context(), resourceID)
	if err != nil {
		h.logger.Errorw("Error getting resource metrics", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(metrics)
}

// GetPlatformReport handles GET /reports/platform
func (h *AnalyticsHandler) GetPlatformReport(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	period := analytics.ReportPeriodMonthly
	if p := r.URL.Query().Get("period"); p != "" {
		period = analytics.ReportPeriod(p)
	}

	report, err := h.analyticsService.GeneratePlatformReport(r.Context(), period)
	if err != nil {
		h.logger.Errorw("Error generating platform report", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(report)
}

// GetUserReport handles GET /reports/user/{userId}
func (h *AnalyticsHandler) GetUserReport(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	userID := r.PathValue("userId")
	if userID == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "User ID required"})
		return
	}

	period := analytics.ReportPeriodMonthly
	if p := r.URL.Query().Get("period"); p != "" {
		period = analytics.ReportPeriod(p)
	}

	report, err := h.analyticsService.GenerateUserReport(r.Context(), userID, period)
	if err != nil {
		h.logger.Errorw("Error generating user report", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(report)
}

// GetResourceReport handles GET /reports/resource/{resourceId}
func (h *AnalyticsHandler) GetResourceReport(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	resourceID := r.PathValue("resourceId")
	if resourceID == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Resource ID required"})
		return
	}

	period := analytics.ReportPeriodMonthly
	if p := r.URL.Query().Get("period"); p != "" {
		period = analytics.ReportPeriod(p)
	}

	report, err := h.analyticsService.GenerateResourceReport(r.Context(), resourceID, period)
	if err != nil {
		h.logger.Errorw("Error generating resource report", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(report)
}

// GetPlatformInsights handles GET /insights/platform
func (h *AnalyticsHandler) GetPlatformInsights(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	insights, err := h.analyticsService.GetInsights(r.Context())
	if err != nil {
		h.logger.Errorw("Error getting insights", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(insights)
}

// GetDashboardSummary handles GET /dashboard/summary
func (h *AnalyticsHandler) GetDashboardSummary(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	summary, err := h.analyticsService.GetDashboardSummary(r.Context())
	if err != nil {
		h.logger.Errorw("Error getting dashboard summary", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(summary)
}
