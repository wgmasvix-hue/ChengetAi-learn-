package analytics

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

// Tracker tracks platform analytics events and metrics
type Tracker struct {
	db     *pgxpool.Pool
	logger *zap.SugaredLogger
}

// Event represents an analytics event
type Event struct {
	EventType    string            `json:"eventType"`
	UserID       string            `json:"userId"`
	ResourceID   string            `json:"resourceId"`
	SchoolID     string            `json:"schoolId"`
	Action       string            `json:"action"`
	Metadata     map[string]string `json:"metadata"`
	Timestamp    int64             `json:"timestamp"`
	SessionID    string            `json:"sessionId"`
}

// EventType constants
const (
	EventTypeView          = "view"
	EventTypeDownload      = "download"
	EventTypeBookmark      = "bookmark"
	EventTypeQuizSubmit    = "quiz_submit"
	EventTypeRecommendationClick = "recommendation_click"
	EventTypeLogin         = "login"
	EventTypeLogout        = "logout"
	EventTypeResourceCreate = "resource_create"
	EventTypeWalletTransaction = "wallet_transaction"
)

// NewTracker creates a new analytics tracker
func NewTracker(db *pgxpool.Pool, logger *zap.SugaredLogger) *Tracker {
	return &Tracker{
		db:     db,
		logger: logger,
	}
}

// TrackEvent records an analytics event
func (t *Tracker) TrackEvent(ctx context.Context, event *Event) error {
	if event.EventType == "" {
		return fmt.Errorf("event type cannot be empty")
	}

	if event.Timestamp == 0 {
		event.Timestamp = time.Now().Unix()
	}

	query := `
		INSERT INTO analytics_events (event_type, user_id, resource_id, school_id, action, metadata, session_id, timestamp)
		VALUES ($1, $2, $3, $4, $5, $6, $7, to_timestamp($8))
	`

	_, err := t.db.Exec(ctx, query,
		event.EventType,
		event.UserID,
		event.ResourceID,
		event.SchoolID,
		event.Action,
		event.Metadata,
		event.SessionID,
		event.Timestamp,
	)

	if err != nil {
		t.logger.Errorw("Error tracking event", "error", err, "event_type", event.EventType)
		return fmt.Errorf("tracking event: %w", err)
	}

	return nil
}

// TrackBatchEvents records multiple events
func (t *Tracker) TrackBatchEvents(ctx context.Context, events []Event) error {
	if len(events) == 0 {
		return fmt.Errorf("events cannot be empty")
	}

	for _, event := range events {
		if err := t.TrackEvent(ctx, &event); err != nil {
			t.logger.Warnw("Failed to track event", "error", err)
			// Continue tracking other events
		}
	}

	return nil
}

// UserEngagementMetrics represents user engagement statistics
type UserEngagementMetrics struct {
	UserID            string
	TotalSessions     int
	TotalViews        int
	TotalDownloads    int
	TotalBookmarks    int
	TotalQuizzesAttempted int
	AverageQuizScore  float32
	LastActivityAt    time.Time
	CreatedAt         time.Time
}

// GetUserEngagementMetrics retrieves engagement metrics for a user
func (t *Tracker) GetUserEngagementMetrics(ctx context.Context, userID string) (*UserEngagementMetrics, error) {
	query := `
		SELECT
			user_id,
			COUNT(DISTINCT session_id) as total_sessions,
			COALESCE(SUM(CASE WHEN event_type = 'view' THEN 1 ELSE 0 END), 0) as total_views,
			COALESCE(SUM(CASE WHEN event_type = 'download' THEN 1 ELSE 0 END), 0) as total_downloads,
			COALESCE(SUM(CASE WHEN event_type = 'bookmark' THEN 1 ELSE 0 END), 0) as total_bookmarks,
			COALESCE(SUM(CASE WHEN event_type = 'quiz_submit' THEN 1 ELSE 0 END), 0) as total_quizzes,
			COALESCE(AVG(CAST(metadata->>'score' AS FLOAT)), 0) as avg_quiz_score,
			MAX(timestamp) as last_activity,
			MIN(timestamp) as first_activity
		FROM analytics_events
		WHERE user_id = $1
		GROUP BY user_id
	`

	var metrics UserEngagementMetrics
	var lastActivitySQL, firstActivitySQL *time.Time

	err := t.db.QueryRow(ctx, query, userID).Scan(
		&metrics.UserID,
		&metrics.TotalSessions,
		&metrics.TotalViews,
		&metrics.TotalDownloads,
		&metrics.TotalBookmarks,
		&metrics.TotalQuizzesAttempted,
		&metrics.AverageQuizScore,
		&lastActivitySQL,
		&firstActivitySQL,
	)

	if err != nil {
		return nil, fmt.Errorf("querying metrics: %w", err)
	}

	if lastActivitySQL != nil {
		metrics.LastActivityAt = *lastActivitySQL
	}
	if firstActivitySQL != nil {
		metrics.CreatedAt = *firstActivitySQL
	}

	return &metrics, nil
}

// ResourceAnalytics represents statistics for a resource
type ResourceAnalytics struct {
	ResourceID        string
	Title             string
	TotalViews        int
	TotalDownloads    int
	TotalBookmarks    int
	UniqueViewers     int
	UniqueDownloaders int
	AverageEngagement float32
	Trending          bool
}

// GetResourceAnalytics retrieves analytics for a resource
func (t *Tracker) GetResourceAnalytics(ctx context.Context, resourceID string) (*ResourceAnalytics, error) {
	query := `
		SELECT
			r.id,
			r.title,
			r.views_count,
			r.downloads_count,
			r.bookmarks_count,
			COALESCE(COUNT(DISTINCT CASE WHEN ae.event_type = 'view' THEN ae.user_id END), 0) as unique_viewers,
			COALESCE(COUNT(DISTINCT CASE WHEN ae.event_type = 'download' THEN ae.user_id END), 0) as unique_downloaders
		FROM resources r
		LEFT JOIN analytics_events ae ON r.id = ae.resource_id
		WHERE r.id = $1
		GROUP BY r.id, r.title, r.views_count, r.downloads_count, r.bookmarks_count
	`

	var analytics ResourceAnalytics

	err := t.db.QueryRow(ctx, query, resourceID).Scan(
		&analytics.ResourceID,
		&analytics.Title,
		&analytics.TotalViews,
		&analytics.TotalDownloads,
		&analytics.TotalBookmarks,
		&analytics.UniqueViewers,
		&analytics.UniqueDownloaders,
	)

	if err != nil {
		return nil, fmt.Errorf("querying resource analytics: %w", err)
	}

	// Calculate engagement score
	if analytics.TotalViews > 0 {
		analytics.AverageEngagement = float32(analytics.TotalDownloads+analytics.TotalBookmarks) / float32(analytics.TotalViews) * 100
	}

	return &analytics, nil
}

// PlatformMetrics represents overall platform statistics
type PlatformMetrics struct {
	TotalUsers           int
	TotalResources       int
	TotalViews           int
	TotalDownloads       int
	AverageSessionLength int // seconds
	DailyActiveUsers     int
	NewUsersToday        int
	TrendingResourcesCount int
}

// GetPlatformMetrics retrieves platform-wide statistics
func (t *Tracker) GetPlatformMetrics(ctx context.Context, days int) (*PlatformMetrics, error) {
	if days == 0 {
		days = 30
	}

	query := `
		SELECT
			(SELECT COUNT(*) FROM users) as total_users,
			(SELECT COUNT(*) FROM resources WHERE status = 'published') as total_resources,
			(SELECT COUNT(*) FROM analytics_events WHERE event_type = 'view' AND timestamp > NOW() - INTERVAL '1 day' * $1) as total_views,
			(SELECT COUNT(*) FROM analytics_events WHERE event_type = 'download' AND timestamp > NOW() - INTERVAL '1 day' * $1) as total_downloads,
			COALESCE((SELECT AVG(session_length) FROM user_sessions), 0)::int as avg_session_length,
			(SELECT COUNT(DISTINCT user_id) FROM analytics_events WHERE timestamp > NOW() - INTERVAL '1 day') as dau,
			(SELECT COUNT(*) FROM users WHERE created_at > NOW() - INTERVAL '1 day') as new_users_today,
			(SELECT COUNT(*) FROM resources WHERE views_count > (SELECT AVG(views_count) FROM resources) * 1.5) as trending_count
	`

	var metrics PlatformMetrics

	err := t.db.QueryRow(ctx, query, days).Scan(
		&metrics.TotalUsers,
		&metrics.TotalResources,
		&metrics.TotalViews,
		&metrics.TotalDownloads,
		&metrics.AverageSessionLength,
		&metrics.DailyActiveUsers,
		&metrics.NewUsersToday,
		&metrics.TrendingResourcesCount,
	)

	if err != nil {
		return nil, fmt.Errorf("querying platform metrics: %w", err)
	}

	return &metrics, nil
}
