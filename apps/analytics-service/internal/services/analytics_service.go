package services

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/wgmasvix-hue/ChengetAi-learn-/platform/pkg/analytics"
	"go.uber.org/zap"
)

// AnalyticsService orchestrates analytics operations
type AnalyticsService struct {
	tracker  *analytics.Tracker
	reporter *analytics.Reporter
	db       *pgxpool.Pool
	logger   *zap.SugaredLogger
}

// NewAnalyticsService creates a new analytics service
func NewAnalyticsService(
	tracker *analytics.Tracker,
	reporter *analytics.Reporter,
	db *pgxpool.Pool,
	logger *zap.SugaredLogger,
) *AnalyticsService {
	return &AnalyticsService{
		tracker:  tracker,
		reporter: reporter,
		db:       db,
		logger:   logger,
	}
}

// TrackEvent records a single analytics event
func (s *AnalyticsService) TrackEvent(ctx context.Context, event *analytics.Event) error {
	s.logger.Debugw("Tracking event", "event_type", event.EventType, "user_id", event.UserID)
	return s.tracker.TrackEvent(ctx, event)
}

// TrackBatchEvents records multiple events
func (s *AnalyticsService) TrackBatchEvents(ctx context.Context, events []analytics.Event) error {
	s.logger.Debugw("Tracking batch events", "count", len(events))
	return s.tracker.TrackBatchEvents(ctx, events)
}

// GetPlatformMetrics retrieves platform-wide metrics
func (s *AnalyticsService) GetPlatformMetrics(ctx context.Context, days int) (*analytics.PlatformMetrics, error) {
	metrics, err := s.tracker.GetPlatformMetrics(ctx, days)
	if err != nil {
		s.logger.Errorw("Error getting platform metrics", "error", err)
		return nil, err
	}
	return metrics, nil
}

// GetUserMetrics retrieves metrics for a specific user
func (s *AnalyticsService) GetUserMetrics(ctx context.Context, userID string) (*analytics.UserEngagementMetrics, error) {
	metrics, err := s.tracker.GetUserEngagementMetrics(ctx, userID)
	if err != nil {
		s.logger.Errorw("Error getting user metrics", "error", err, "user_id", userID)
		return nil, err
	}
	return metrics, nil
}

// GetResourceMetrics retrieves metrics for a specific resource
func (s *AnalyticsService) GetResourceMetrics(ctx context.Context, resourceID string) (*analytics.ResourceAnalytics, error) {
	metrics, err := s.tracker.GetResourceAnalytics(ctx, resourceID)
	if err != nil {
		s.logger.Errorw("Error getting resource metrics", "error", err, "resource_id", resourceID)
		return nil, err
	}
	return metrics, nil
}

// GeneratePlatformReport generates a platform report
func (s *AnalyticsService) GeneratePlatformReport(ctx context.Context, period analytics.ReportPeriod) (*analytics.Report, error) {
	s.logger.Debugw("Generating platform report", "period", period)

	metrics, err := s.GetPlatformMetrics(ctx, 30)
	if err != nil {
		return nil, err
	}

	report := s.reporter.GeneratePlatformReport(metrics, period)

	// Save report to database
	if err := s.saveReport(ctx, report); err != nil {
		s.logger.Warnw("Error saving report", "error", err)
		// Don't fail - report is still valid
	}

	return report, nil
}

// GenerateUserReport generates a user report
func (s *AnalyticsService) GenerateUserReport(ctx context.Context, userID string, period analytics.ReportPeriod) (*analytics.Report, error) {
	s.logger.Debugw("Generating user report", "user_id", userID, "period", period)

	metrics, err := s.GetUserMetrics(ctx, userID)
	if err != nil {
		return nil, err
	}

	report := s.reporter.GenerateUserReport(userID, metrics, period)

	if err := s.saveReport(ctx, report); err != nil {
		s.logger.Warnw("Error saving report", "error", err)
	}

	return report, nil
}

// GenerateResourceReport generates a resource report
func (s *AnalyticsService) GenerateResourceReport(ctx context.Context, resourceID string, period analytics.ReportPeriod) (*analytics.Report, error) {
	s.logger.Debugw("Generating resource report", "resource_id", resourceID, "period", period)

	metrics, err := s.GetResourceMetrics(ctx, resourceID)
	if err != nil {
		return nil, err
	}

	report := s.reporter.GenerateResourceReport(metrics, period)

	if err := s.saveReport(ctx, report); err != nil {
		s.logger.Warnw("Error saving report", "error", err)
	}

	return report, nil
}

// GetInsights generates actionable insights
func (s *AnalyticsService) GetInsights(ctx context.Context) ([]analytics.AnalyticsInsight, error) {
	s.logger.Debugw("Generating insights")

	metrics, err := s.GetPlatformMetrics(ctx, 30)
	if err != nil {
		return nil, err
	}

	insights := s.reporter.GenerateInsights(metrics)
	return insights, nil
}

// DashboardSummary represents a dashboard overview
type DashboardSummary struct {
	Metrics  *analytics.PlatformMetrics `json:"metrics"`
	Insights []analytics.AnalyticsInsight `json:"insights"`
	Reports  map[string]interface{} `json:"reports"`
}

// GetDashboardSummary generates a dashboard overview
func (s *AnalyticsService) GetDashboardSummary(ctx context.Context) (*DashboardSummary, error) {
	s.logger.Debugw("Generating dashboard summary")

	metrics, err := s.GetPlatformMetrics(ctx, 30)
	if err != nil {
		return nil, fmt.Errorf("getting metrics: %w", err)
	}

	insights := s.reporter.GenerateInsights(metrics)

	summary := &DashboardSummary{
		Metrics:  metrics,
		Insights: insights,
		Reports: map[string]interface{}{
			"lastGenerated": 1704067200,
		},
	}

	return summary, nil
}

// saveReport saves a report to database
func (s *AnalyticsService) saveReport(ctx context.Context, report *analytics.Report) error {
	query := `
		INSERT INTO analytics_reports (id, title, type, period, metrics, generated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		ON CONFLICT (id) DO NOTHING
	`

	_, err := s.db.Exec(ctx, query,
		report.ID,
		report.Title,
		report.Type,
		string(report.Period),
		report.Metrics,
		report.GeneratedAt,
	)

	return err
}

// DeleteOldEvents removes analytics events older than retention period
func (s *AnalyticsService) DeleteOldEvents(ctx context.Context, retentionDays int) error {
	query := `DELETE FROM analytics_events WHERE timestamp < NOW() - INTERVAL '1 day' * $1`
	_, err := s.db.Exec(ctx, query, retentionDays)
	if err != nil {
		s.logger.Errorw("Error deleting old events", "error", err)
		return err
	}
	return nil
}
