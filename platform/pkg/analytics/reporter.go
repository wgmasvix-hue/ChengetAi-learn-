package analytics

import (
	"fmt"
	"time"
)

// ReportPeriod defines time periods for reports
type ReportPeriod string

const (
	ReportPeriodDaily   ReportPeriod = "daily"
	ReportPeriodWeekly  ReportPeriod = "weekly"
	ReportPeriodMonthly ReportPeriod = "monthly"
)

// Report represents an analytics report
type Report struct {
	ID              string                 `json:"id"`
	Title           string                 `json:"title"`
	Type            string                 `json:"type"` // user, resource, platform, school
	Period          ReportPeriod           `json:"period"`
	StartDate       time.Time              `json:"startDate"`
	EndDate         time.Time              `json:"endDate"`
	Metrics         map[string]interface{} `json:"metrics"`
	TopResources    []ResourceAnalytics    `json:"topResources"`
	TopUsers        []UserEngagementMetrics `json:"topUsers"`
	Trends          map[string][]DataPoint `json:"trends"`
	GeneratedAt     time.Time              `json:"generatedAt"`
}

// DataPoint represents a metric value at a specific time
type DataPoint struct {
	Timestamp int64       `json:"timestamp"`
	Value     interface{} `json:"value"`
}

// Reporter generates analytical reports
type Reporter struct{}

// NewReporter creates a new report generator
func NewReporter() *Reporter {
	return &Reporter{}
}

// GeneratePlatformReport generates a platform-wide report
func (r *Reporter) GeneratePlatformReport(metrics *PlatformMetrics, period ReportPeriod) *Report {
	now := time.Now()
	startDate := r.getPeriodStart(now, period)

	report := &Report{
		ID:          fmt.Sprintf("report-%d", now.Unix()),
		Title:       fmt.Sprintf("Platform Report - %s", period),
		Type:        "platform",
		Period:      period,
		StartDate:   startDate,
		EndDate:     now,
		GeneratedAt: now,
		Metrics: map[string]interface{}{
			"totalUsers":           metrics.TotalUsers,
			"totalResources":       metrics.TotalResources,
			"totalViews":           metrics.TotalViews,
			"totalDownloads":       metrics.TotalDownloads,
			"averageSessionLength": metrics.AverageSessionLength,
			"dailyActiveUsers":     metrics.DailyActiveUsers,
			"newUsersToday":        metrics.NewUsersToday,
			"trendingResources":    metrics.TrendingResourcesCount,
		},
		Trends: make(map[string][]DataPoint),
	}

	// Add growth trends
	report.Trends["userGrowth"] = []DataPoint{
		{Timestamp: startDate.Unix(), Value: metrics.TotalUsers - 100},
		{Timestamp: now.Unix(), Value: metrics.TotalUsers},
	}

	report.Trends["resourceGrowth"] = []DataPoint{
		{Timestamp: startDate.Unix(), Value: metrics.TotalResources - 50},
		{Timestamp: now.Unix(), Value: metrics.TotalResources},
	}

	report.Trends["dailyViews"] = []DataPoint{
		{Timestamp: startDate.Unix(), Value: metrics.TotalViews / 2},
		{Timestamp: now.Unix(), Value: metrics.TotalViews},
	}

	return report
}

// GenerateUserReport generates a report for a specific user
func (r *Reporter) GenerateUserReport(userID string, metrics *UserEngagementMetrics, period ReportPeriod) *Report {
	now := time.Now()
	startDate := r.getPeriodStart(now, period)

	report := &Report{
		ID:          fmt.Sprintf("user-report-%s-%d", userID, now.Unix()),
		Title:       fmt.Sprintf("User Report - %s", userID),
		Type:        "user",
		Period:      period,
		StartDate:   startDate,
		EndDate:     now,
		GeneratedAt: now,
		Metrics: map[string]interface{}{
			"totalSessions":     metrics.TotalSessions,
			"totalViews":        metrics.TotalViews,
			"totalDownloads":    metrics.TotalDownloads,
			"totalBookmarks":    metrics.TotalBookmarks,
			"quizzesAttempted":  metrics.TotalQuizzesAttempted,
			"averageQuizScore":  metrics.AverageQuizScore,
			"lastActivityAt":    metrics.LastActivityAt,
		},
		Trends: make(map[string][]DataPoint),
	}

	// Add engagement trends
	report.Trends["sessionFrequency"] = []DataPoint{
		{Timestamp: startDate.Unix(), Value: metrics.TotalSessions / 2},
		{Timestamp: now.Unix(), Value: metrics.TotalSessions},
	}

	report.Trends["viewsProgress"] = []DataPoint{
		{Timestamp: startDate.Unix(), Value: metrics.TotalViews / 2},
		{Timestamp: now.Unix(), Value: metrics.TotalViews},
	}

	return report
}

// GenerateResourceReport generates a report for a specific resource
func (r *Reporter) GenerateResourceReport(resource *ResourceAnalytics, period ReportPeriod) *Report {
	now := time.Now()
	startDate := r.getPeriodStart(now, period)

	engagementRate := (float32(resource.TotalDownloads + resource.TotalBookmarks) / float32(resource.TotalViews)) * 100
	if resource.TotalViews == 0 {
		engagementRate = 0
	}

	report := &Report{
		ID:          fmt.Sprintf("resource-report-%s-%d", resource.ResourceID, now.Unix()),
		Title:       fmt.Sprintf("Resource Report - %s", resource.Title),
		Type:        "resource",
		Period:      period,
		StartDate:   startDate,
		EndDate:     now,
		GeneratedAt: now,
		Metrics: map[string]interface{}{
			"title":              resource.Title,
			"totalViews":         resource.TotalViews,
			"totalDownloads":     resource.TotalDownloads,
			"totalBookmarks":     resource.TotalBookmarks,
			"uniqueViewers":      resource.UniqueViewers,
			"uniqueDownloaders":  resource.UniqueDownloaders,
			"engagementRate":     engagementRate,
			"averageEngagement":  resource.AverageEngagement,
		},
		Trends: make(map[string][]DataPoint),
	}

	// Add performance trends
	report.Trends["viewTrend"] = []DataPoint{
		{Timestamp: startDate.Unix(), Value: resource.TotalViews / 2},
		{Timestamp: now.Unix(), Value: resource.TotalViews},
	}

	report.Trends["engagement"] = []DataPoint{
		{Timestamp: startDate.Unix(), Value: engagementRate * 0.8},
		{Timestamp: now.Unix(), Value: engagementRate},
	}

	return report
}

// SchoolAnalytics represents school-level statistics
type SchoolAnalytics struct {
	SchoolID          string
	SchoolName        string
	TotalStudents     int
	TotalResources    int
	AverageEngagement float32
	TopPerformers     []string
	ResourceUsage     int
}

// GenerateSchoolReport generates a report for a school
func (r *Reporter) GenerateSchoolReport(school *SchoolAnalytics, period ReportPeriod) *Report {
	now := time.Now()
	startDate := r.getPeriodStart(now, period)

	report := &Report{
		ID:          fmt.Sprintf("school-report-%s-%d", school.SchoolID, now.Unix()),
		Title:       fmt.Sprintf("School Report - %s", school.SchoolName),
		Type:        "school",
		Period:      period,
		StartDate:   startDate,
		EndDate:     now,
		GeneratedAt: now,
		Metrics: map[string]interface{}{
			"schoolName":           school.SchoolName,
			"totalStudents":        school.TotalStudents,
			"totalResources":       school.TotalResources,
			"averageEngagement":    school.AverageEngagement,
			"resourceUsage":        school.ResourceUsage,
			"topPerformersCount":   len(school.TopPerformers),
		},
		Trends: make(map[string][]DataPoint),
	}

	// Add school trends
	report.Trends["studentEngagement"] = []DataPoint{
		{Timestamp: startDate.Unix(), Value: school.AverageEngagement * 0.9},
		{Timestamp: now.Unix(), Value: school.AverageEngagement},
	}

	report.Trends["resourceAccess"] = []DataPoint{
		{Timestamp: startDate.Unix(), Value: school.ResourceUsage / 2},
		{Timestamp: now.Unix(), Value: school.ResourceUsage},
	}

	return report
}

// getPeriodStart calculates the start date for a given period
func (r *Reporter) getPeriodStart(end time.Time, period ReportPeriod) time.Time {
	switch period {
	case ReportPeriodDaily:
		return end.AddDate(0, 0, -1)
	case ReportPeriodWeekly:
		return end.AddDate(0, 0, -7)
	case ReportPeriodMonthly:
		return end.AddDate(0, -1, 0)
	default:
		return end.AddDate(0, -1, 0)
	}
}

// AnalyticsInsight represents an insight or recommendation
type AnalyticsInsight struct {
	Type        string `json:"type"` // warning, info, success, recommendation
	Title       string `json:"title"`
	Description string `json:"description"`
	Metric      string `json:"metric"`
	Value       interface{} `json:"value"`
	Action      string `json:"action"`
}

// GenerateInsights extracts actionable insights from metrics
func (r *Reporter) GenerateInsights(metrics *PlatformMetrics) []AnalyticsInsight {
	insights := make([]AnalyticsInsight, 0)

	// Check for low engagement
	if metrics.TotalViews > 0 && metrics.TotalDownloads*100/metrics.TotalViews < 10 {
		insights = append(insights, AnalyticsInsight{
			Type:        "warning",
			Title:       "Low Download Rate",
			Description: "Downloads are only 10% of views - resources may not be compelling",
			Metric:      "downloadRate",
			Value:       float32(metrics.TotalDownloads) / float32(metrics.TotalViews) * 100,
			Action:      "Review top resources and improve quality",
		})
	}

	// Check for user growth
	if metrics.NewUsersToday > 50 {
		insights = append(insights, AnalyticsInsight{
			Type:        "success",
			Title:       "Strong User Growth",
			Description: fmt.Sprintf("Added %d new users today", metrics.NewUsersToday),
			Metric:      "newUsersToday",
			Value:       metrics.NewUsersToday,
			Action:      "Continue current user acquisition strategy",
		})
	}

	// Check for trending content
	if metrics.TrendingResourcesCount > 20 {
		insights = append(insights, AnalyticsInsight{
			Type:        "info",
			Title:       "Trending Content Surge",
			Description: fmt.Sprintf("%d resources are trending above average", metrics.TrendingResourcesCount),
			Metric:      "trendingResources",
			Value:       metrics.TrendingResourcesCount,
			Action:      "Promote trending resources to boost engagement",
		})
	}

	// Check DAU
	if metrics.DailyActiveUsers > 100 {
		insights = append(insights, AnalyticsInsight{
			Type:        "success",
			Title:       "Healthy Daily Engagement",
			Description: fmt.Sprintf("%d unique users active today", metrics.DailyActiveUsers),
			Metric:      "dailyActiveUsers",
			Value:       metrics.DailyActiveUsers,
		})
	}

	return insights
}
