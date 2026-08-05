# Analytics Service

The Analytics Service provides comprehensive tracking, reporting, and insights for the ChengetAi platform.

## Features

- **Event Tracking**: Track user actions and platform events in real-time
- **Real-time Metrics**: Platform-wide, user, and resource metrics
- **Report Generation**: Automated daily/weekly/monthly reports
- **Dashboards**: Executive dashboards with KPI summaries
- **Insights**: AI-powered recommendations and alerts
- **Data Retention**: Configurable retention policies for analytics data
- **Performance Analytics**: Track resource usage and engagement

## Port

**8012**

## Endpoints

### Event Tracking

#### Track Single Event
```
POST /events/track

Request:
{
  "eventType": "view",
  "userId": "user-123",
  "resourceId": "resource-456",
  "action": "viewed_resource",
  "metadata": {"duration": "120"}
}

Response:
{
  "status": "event tracked"
}
```

#### Track Batch Events
```
POST /events/batch

Request:
[
  {
    "eventType": "view",
    "userId": "user-123",
    "resourceId": "resource-456"
  },
  {
    "eventType": "download",
    "userId": "user-789",
    "resourceId": "resource-456"
  }
]

Response:
{
  "status": "events tracked",
  "count": 2
}
```

### Metrics

#### Platform Metrics
```
GET /metrics/platform?days=30

Response:
{
  "totalUsers": 1500,
  "totalResources": 5000,
  "totalViews": 125000,
  "totalDownloads": 45000,
  "dailyActiveUsers": 350,
  "newUsersToday": 42,
  "trendingResources": 28
}
```

#### User Metrics
```
GET /metrics/user/{userId}

Response:
{
  "userId": "user-123",
  "totalSessions": 25,
  "totalViews": 150,
  "totalDownloads": 45,
  "totalBookmarks": 12,
  "averageQuizScore": 85.5
}
```

#### Resource Metrics
```
GET /metrics/resource/{resourceId}

Response:
{
  "resourceId": "resource-456",
  "title": "Biology 101",
  "totalViews": 2500,
  "totalDownloads": 800,
  "totalBookmarks": 200,
  "uniqueViewers": 1200,
  "engagementRate": 40.0
}
```

### Reports

#### Platform Report
```
GET /reports/platform?period=monthly

Response:
{
  "id": "report-123",
  "title": "Platform Report - monthly",
  "type": "platform",
  "period": "monthly",
  "metrics": {...},
  "trends": {...}
}
```

#### User Report
```
GET /reports/user/{userId}?period=weekly

Response:
{
  "id": "user-report-123",
  "title": "User Report - user-123",
  "type": "user",
  "metrics": {...}
}
```

#### Resource Report
```
GET /reports/resource/{resourceId}?period=monthly

Response:
{
  "id": "resource-report-123",
  "title": "Resource Report - Biology 101",
  "type": "resource",
  "metrics": {...}
}
```

### Insights & Dashboard

#### Platform Insights
```
GET /insights/platform

Response:
[
  {
    "type": "success",
    "title": "Strong User Growth",
    "description": "Added 42 new users today",
    "metric": "newUsersToday",
    "value": 42,
    "action": "Continue current strategy"
  }
]
```

#### Dashboard Summary
```
GET /dashboard/summary

Response:
{
  "metrics": {...},
  "insights": [...],
  "reports": {...}
}
```

## Environment Variables

```
SERVICE_NAME=analytics-service
SERVICE_PORT=8012
ENVIRONMENT=development

DB_HOST=localhost
DB_PORT=5432
DB_NAME=chengetai
DB_USER=chengetai
DB_PASSWORD=
DB_MAX_CONNS=25
DB_MIN_CONNS=5

ANALYTICS_RETENTION_DAYS=90
ANALYTICS_BATCH_SIZE=100
REPORT_GENERATION_TIME=00:00

LOG_LEVEL=info
```

## Database Tables

- **analytics_events**: Raw event log with 15+ indices
- **analytics_reports**: Generated reports snapshot
- **event_aggregations**: Pre-computed hourly/daily metrics

## Event Types

- `view`: Resource viewed
- `download`: Resource downloaded
- `bookmark`: Resource bookmarked
- `quiz_submit`: Quiz submission
- `recommendation_click`: Recommendation clicked
- `login`: User login
- `logout`: User logout
- `resource_create`: Resource created
- `wallet_transaction`: Wallet transaction

## Architecture

- **Tracker**: Records events to database
- **Reporter**: Generates reports from metrics
- **Service Layer**: Orchestrates operations
- **Handler Layer**: HTTP endpoints
- **Health Checks**: `/health` and `/ready` endpoints

## Metrics & KPIs

### Platform Metrics
- Daily/Weekly/Monthly Active Users (DAU/WAU/MAU)
- View/Download/Engagement rates
- New user acquisitions
- Resource trending

### User Metrics
- Session count and duration
- Resource views/downloads
- Quiz performance
- Learning progress

### Resource Metrics
- View count
- Download count
- Engagement rate (downloads + bookmarks / views)
- Unique viewers
- Trending status

## Report Periods

- **Daily**: Last 24 hours
- **Weekly**: Last 7 days
- **Monthly**: Last 30 days

## Development

```bash
# Build
make build

# Run
make run

# Test
make test

# Lint
make lint

# Docker
make docker-build
```

## Integration

- **All Services**: Send analytics events
- **Dashboard**: Consumes insights and reports
- **Recommendations**: Uses engagement metrics
- **Reports**: Scheduled generation and distribution
