# Quiz Service

The Quiz Service handles educational assessment through intelligent quiz generation, delivery, and scoring.

## Features

- **Quiz Generation**: Automatically generate questions from educational content
- **Question Types**: Multiple choice, true/false, short answer, essay
- **Difficulty Levels**: Easy, medium, hard with adaptive distribution
- **Auto-Scoring**: Immediate feedback with detailed scoring
- **Learning Analytics**: Track progress with statistics and performance metrics
- **Multi-language**: Support for assessments in multiple languages

## Port

**8005**

## Endpoints

### Quiz Management

#### Generate Quiz
```
POST /quizzes/generate

Request:
{
  "resourceId": "resource-123",
  "content": "Full text content...",
  "questionCount": 10,
  "difficultyLevel": "medium",
  "questionTypes": ["multiple_choice", "true_false"],
  "language": "en"
}

Response:
{
  "id": "quiz-123",
  "resourceId": "resource-123",
  "questions": [...],
  "estimatedTime": 20,
  "passingScore": 70
}
```

#### Get Quiz
```
GET /quizzes/{quizId}

Response:
{
  "id": "quiz-123",
  "title": "Biology Quiz",
  "questions": [...],
  "estimatedTime": 20
}
```

#### Submit Response
```
POST /quizzes/{quizId}/submit

Request:
{
  "userId": "user-123",
  "responses": [
    {
      "questionId": "q-1",
      "answerText": "a"
    }
  ]
}

Response:
{
  "id": "response-123",
  "score": 80,
  "percentage": 80.0,
  "timeSpent": 1200
}
```

#### Get User Quizzes
```
GET /quizzes/user/{userId}?limit=20&offset=0

Response:
[
  {
    "id": "response-123",
    "quizId": "quiz-123",
    "score": 80,
    "percentage": 80.0,
    "submittedAt": 1704067200
  }
]
```

#### Get Response Details
```
GET /responses/{responseId}

Response:
{
  "id": "response-123",
  "quizId": "quiz-123",
  "userId": "user-123",
  "responses": [...],
  "score": 80,
  "percentage": 80.0
}
```

## Environment Variables

```
SERVICE_NAME=quiz-service
SERVICE_PORT=8005
ENVIRONMENT=development

DB_HOST=localhost
DB_PORT=5432
DB_NAME=chengetai
DB_USER=chengetai
DB_PASSWORD=
DB_MAX_CONNS=25
DB_MIN_CONNS=5

DEFAULT_QUESTION_COUNT=10
MAX_QUESTION_COUNT=100
DEFAULT_DIFFICULTY=medium

LOG_LEVEL=info
```

## Database Tables

- **quizzes**: Quiz definitions with questions stored as JSON
- **quiz_responses**: User responses with scoring
- Supporting indices for efficient queries

## Architecture

- **Generator**: Creates questions using LLM
- **Repositories**: Data access for quizzes and responses
- **Scoring Engine**: Automated assessment and grading
- **Analytics**: Performance tracking and statistics
- **Health Checks**: `/health` and `/ready` endpoints

## Question Types

1. **Multiple Choice** - 4 options with single correct answer
2. **True/False** - Binary choice with explanation
3. **Short Answer** - Text response with semantic matching
4. **Essay** - Open-ended assessment (manual review)

## Scoring

- Questions carry point values based on difficulty
- Immediate feedback after submission
- Passing score: 70% (configurable)
- Performance analytics and progress tracking

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

- **AI Service**: Generates questions using Claude
- **Auth Service**: Validates user identity
- **Search Service**: Links quizzes to resources
