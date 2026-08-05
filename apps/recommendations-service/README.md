# Recommendations Service

The Recommendations Service delivers personalized learning recommendations using collaborative filtering, content-based filtering, and hybrid approaches.

## Features

- **Hybrid Recommendations**: Combines collaborative and content-based filtering
- **Personalization**: Tailored to individual learning styles and preferences
- **Diversity Metrics**: Ensures variety in recommendations
- **Feedback Loop**: Learns from user engagement patterns
- **Performance Tracking**: Measures recommendation effectiveness

## Port

**8006**

## Endpoints

### Recommendations

#### Get Recommendations
```
GET /recommendations/{userId}?topK=10&strategy=hybrid

Query Parameters:
- topK: Number of recommendations (default: 10)
- strategy: collaborative|content_based|hybrid (default: hybrid)

Response:
{
  "userId": "user-123",
  "recommendations": [
    {
      "resourceId": "resource-1",
      "title": "Biology Fundamentals",
      "score": 0.92,
      "reason": "Similar to resources you've viewed",
      "language": "en"
    }
  ],
  "strategy": "hybrid",
  "generatedAt": 1704067200
}
```

#### Generate Recommendations
```
POST /recommendations/generate

Request:
{
  "userId": "user-123",
  "topK": 10,
  "strategy": "hybrid"
}

Response:
{
  "userId": "user-123",
  "recommendations": [...],
  "strategy": "hybrid"
}
```

#### Get Metrics
```
GET /metrics/{userId}

Response:
{
  "coverage": 0.85,
  "diversity": 0.72,
  "novelty": 0.68,
  "serendipity": 0.45,
  "personalization": 0.92
}
```

#### Submit Feedback
```
POST /feedback

Request:
{
  "userId": "user-123",
  "resourceId": "resource-1",
  "helpful": true,
  "reason": "Exactly what I was looking for"
}

Response:
{
  "status": "feedback recorded"
}
```

## Environment Variables

```
SERVICE_NAME=recommendations-service
SERVICE_PORT=8006
ENVIRONMENT=development

DB_HOST=localhost
DB_PORT=5432
DB_NAME=chengetai
DB_USER=chengetai
DB_PASSWORD=
DB_MAX_CONNS=25
DB_MIN_CONNS=5

DEFAULT_TOP_K=10
DEFAULT_MIN_SIMILARITY=0.5
DEFAULT_STRATEGY=hybrid

LOG_LEVEL=info
```

## Recommendation Strategies

### Collaborative Filtering
- Finds similar users based on behavior patterns
- Recommends resources liked by similar users
- Strength: Discovers new content, good for new topics
- Weakness: Cold start problem for new users

### Content-Based Filtering
- Analyzes resource embeddings and metadata
- Recommends similar resources to ones user viewed
- Strength: No cold start, interpretable
- Weakness: Limited novelty, echo chamber effect

### Hybrid Approach
- Combines both strategies (60% content, 40% collaborative)
- Balances novelty with relevance
- Recommended for general use

## Metrics

- **Coverage**: Percentage of catalog represented in recommendations
- **Diversity**: Variety of resource types in recommendations
- **Novelty**: How new recommendations are vs. user history
- **Serendipity**: Unexpectedly good recommendations
- **Personalization**: How tailored to individual user

## Database Tables

- **recommendation_history**: Tracks generated recommendations
- **recommendation_feedback**: User feedback on recommendations
- Supporting indices for efficient filtering

## Architecture

- **Recommendation Engine**: Computes personalized suggestions
- **Hybrid Strategy**: Blends collaborative and content approaches
- **Metrics Engine**: Measures effectiveness
- **Feedback Loop**: Incorporates user signals
- **Health Checks**: `/health` and `/ready` endpoints

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

- **Search Service**: Retrieves resource embeddings and metadata
- **AI Service**: Uses recommendations in learning paths
- **Auth Service**: Validates user identity
- **Quiz Service**: May recommend quizzes based on content
