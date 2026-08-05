# AI Service

The AI Service provides intelligent, personalized educational features powered by Claude API and Retrieval Augmented Generation (RAG).

## Features

- **AI Tutor**: Context-aware question answering powered by Claude
- **RAG (Retrieval Augmented Generation)**: Grounds AI responses in knowledge base
- **Semantic Search**: Finds relevant resources using embeddings
- **Learning Paths**: Personalized learning recommendations
- **Concept Explanations**: Detailed explanations tailored to difficulty level

## Port

**8004**

## Endpoints

### AI Tutor

#### Ask Question
```
POST /tutor/ask

Request:
{
  "userId": "user-123",
  "query": "What is photosynthesis?",
  "language": "en"
}

Response:
{
  "answer": "Photosynthesis is the process...",
  "sources": ["resource-1", "resource-2"],
  "followUpQuestions": ["Can you explain..."],
  "confidence": 0.85
}
```

#### Explain Concept
```
POST /tutor/explain

Request:
{
  "concept": "Photosynthesis",
  "detailLevel": "standard",
  "language": "en"
}

Response:
{
  "answer": "Detailed explanation...",
  "sources": ["resource-1", "resource-3"]
}
```

#### Get Learning Path
```
GET /tutor/learning-path/{userId}

Query Parameters:
- goal: Learning objective
- currentLevel: beginner|intermediate|advanced
- timeAvailable: Hours per week

Response:
{
  "goal": "Master Biology",
  "milestones": [...],
  "estimatedTime": 12,
  "resources": [...]
}
```

### RAG (Retrieval Augmented Generation)

#### Retrieve Context
```
POST /rag/retrieve

Request:
{
  "query": "photosynthesis",
  "embeddingVector": [...],
  "topK": 5,
  "minSimilarity": 0.5,
  "filters": {"language": "en"}
}

Response:
{
  "resources": [...],
  "totalFound": 10,
  "queryTime": 145
}
```

#### Semantic Search
```
POST /rag/search

Request:
{
  "query": "How photosynthesis works",
  "topK": 5,
  "filters": {"language": "en"}
}

Response:
{
  "resources": [...],
  "totalFound": 8
}
```

## Environment Variables

```
SERVICE_NAME=ai-service
SERVICE_PORT=8004
ENVIRONMENT=development

DB_HOST=localhost
DB_PORT=5432
DB_NAME=chengetai
DB_USER=chengetai
DB_PASSWORD=
DB_MAX_CONNS=25
DB_MIN_CONNS=5

CLAUDE_API_KEY=sk-...
CLAUDE_MODEL=claude-opus-5
CLAUDE_MAX_TOKENS=4096

RAG_MAX_CONTEXT_LENGTH=8192
RAG_MIN_SIMILARITY=0.5
RAG_TOP_K=5

DEFAULT_QUESTION_COUNT=10
DEFAULT_QUESTION_TYPES=multiple_choice,true_false,short_answer

LOG_LEVEL=info
```

## Architecture

- **Claude Integration**: Uses Anthropic Claude API for LLM inference
- **RAG Engine**: Retrieves relevant resources from vector database
- **Semantic Search**: Searches using embedding similarity
- **Learning Paths**: Generates personalized learning plans
- **Health Checks**: `/health` and `/ready` endpoints

## Dependencies

- PostgreSQL with pgvector extension (for embeddings)
- Claude API key (from Anthropic)
- Resource embeddings in database (from Search Service)

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

## Integration with Other Services

- **Search Service**: Retrieves indexed resources
- **Quiz Service**: Can generate quizzes for learning paths
- **Recommendations Service**: Complements personalized suggestions
- **Auth Service**: Validates user identity
