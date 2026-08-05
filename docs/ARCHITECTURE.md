# ChengetAi Platform Architecture

## Overview

ChengetAi is built on seven independent layers that communicate through well-defined interfaces. This document details the architectural patterns, communication protocols, and implementation guidelines.

## Architecture Principles

### 1. Clean Architecture

Every service follows clean architecture layers:

```
Presentation Layer
    ↓
Use Cases / Business Logic
    ↓
Entities / Domain Logic
    ↓
Frameworks & Drivers (DB, API, etc)
```

### 2. Domain-Driven Design (DDD)

- Services are organized around business domains, not technical functions
- Bounded contexts define service boundaries
- Ubiquitous language ensures clarity

### 3. Hexagonal Architecture

Services have:
- **Ports**: Interfaces to external systems
- **Adapters**: Implementations of ports
- No dependencies on external frameworks in core domain

### 4. SOLID Principles

- **S**ingle Responsibility: One reason to change
- **O**pen/Closed: Open for extension, closed for modification
- **L**iskov Substitution: Subtypes must be substitutable
- **I**nterface Segregation: Many specific interfaces
- **D**ependency Inversion: Depend on abstractions

### 5. Repository Pattern

All data access through repositories:
```go
type UserRepository interface {
    GetByID(ctx context.Context, id string) (*User, error)
    Save(ctx context.Context, user *User) error
    Delete(ctx context.Context, id string) error
}
```

## Layer Architecture

### Layer 1: Infrastructure

#### Components
- **Container Runtime**: Docker
- **Orchestration**: Docker Compose (development), Kubernetes (production)
- **Reverse Proxy**: Caddy (automatic HTTPS)
- **Database**: PostgreSQL (primary data store)
- **Cache**: Redis (sessions, real-time data)
- **Object Storage**: MinIO (documents, media)
- **Search Engine**: Typesense (full-text & semantic search)
- **Message Bus**: NATS (async communication)
- **Monitoring**: Prometheus (metrics)
- **Visualization**: Grafana (dashboards)
- **Logging**: Loki (centralized logs)
- **Tracing**: OpenTelemetry (distributed tracing)

#### Design
- All components containerized
- Environment-based configuration
- Health checks for all services
- Graceful shutdown handling
- Load balancing ready

### Layer 2: Platform Services

These services are **shared** across all applications. Each service is independently deployable and follows this structure:

```
service-name/
├── cmd/
│   └── server/
│       └── main.go
├── internal/
│   ├── domain/
│   │   ├── entities.go
│   │   └── errors.go
│   ├── services/
│   │   └── user_service.go
│   ├── repositories/
│   │   └── user_repository.go
│   ├── handlers/
│   │   └── user_handler.go
│   ├── middleware/
│   │   └── auth_middleware.go
│   └── config/
│       └── config.go
├── migrations/
│   └── 001_create_users.sql
├── tests/
│   └── user_service_test.go
├── Dockerfile
├── Makefile
├── go.mod
├── .env.example
└── README.md
```

#### Core Services

**API Gateway** (Port 8000)
- Route requests to appropriate services
- Rate limiting
- Request/response logging
- JWT validation
- API versioning

**Authentication Service** (Port 8001)
- User registration
- Login / logout
- JWT token generation
- OAuth2 integration (Google, Apple)
- 2FA support
- Session management

**User Service** (Port 8002)
- User profiles
- User preferences
- User settings
- Role management
- Permission management

**School Service** (Port 8003)
- School management
- Institutional hierarchy
- Teacher assignments
- Student enrollment
- Curriculum mapping

**Wallet Service** (Port 8004)
- Creator wallets
- Balance tracking
- Transaction history
- Withdrawal management
- Revenue distribution engine

**Payment Service** (Port 8005)
- Payment processing
- Subscription management
- Invoice generation
- Refund handling
- PCI compliance

**Notification Service** (Port 8006)
- Email notifications
- SMS notifications
- Push notifications
- Notification templates
- Notification preferences

**Analytics Service** (Port 8007)
- Event tracking
- Learning analytics
- Usage metrics
- Custom dashboards
- Data export

**Search Service** (Port 8008)
- Full-text search (Typesense)
- Semantic search (embeddings)
- Search indexing
- Real-time updates
- Search analytics

### Layer 3: Discovery Layer

Located at **repo.dare.co.zw**, powered by **DSpace**

#### Capabilities
- REST API: Retrieve collections, items, metadata
- OAI-PMH: Harvest metadata
- SWORD: Automated ingestion
- Versioning: Track resource history
- Persistent Identifiers: DOI, Handle

#### Integration
ChengetAi Knowledge Layer communicates with DSpace via:
```
Knowledge Service
    ↓
DSpace REST API (http://repo.dare.co.zw/rest)
    ↓
Collections → Items → Metadata → Bitstreams
```

#### Knowledge Structure
```
repo.dare.co.zw/handle/
├── communities/
│   └── Africa/
│       ├── Zimbabwe/
│       │   ├── Primary/
│       │   ├── Secondary/O_Level/
│       │   ├── Secondary/A_Level/
│       │   ├── TVET/
│       │   ├── Polytechnic/
│       │   └── University/
│       ├── Kenya/
│       ├── South Africa/
│       └── [Every African Country]
```

### Layer 4: Knowledge Layer (Go)

**Purpose**: Process repository content into AI-ready knowledge

**Responsibilities**:
1. **Harvest**: Query DSpace for new/updated resources
2. **Extract**: Parse documents, OCR images, extract text
3. **Enrich**: Add metadata, classify content
4. **Embed**: Generate vector embeddings
5. **Index**: Build search indexes
6. **Graph**: Create knowledge graphs
7. **Prepare**: Ready content for AI

**Services**:
- Knowledge Harvester: Poll DSpace for changes
- Document Processor: OCR, text extraction
- Embedding Service: Generate vector representations
- Search Indexer: Update Typesense indexes
- Knowledge Graph: Build semantic relationships
- Metadata Classifier: Curriculum, subject, language classification

**Implementation**:
```go
type KnowledgeService interface {
    // Harvest new resources from DSpace
    HarvestResources(ctx context.Context) error
    
    // Process a resource
    ProcessResource(ctx context.Context, resource *Resource) error
    
    // Generate embeddings
    GenerateEmbeddings(ctx context.Context, text string) ([]float32, error)
    
    // Search knowledge base
    Search(ctx context.Context, query string) ([]Result, error)
    
    // Get knowledge context for AI
    GetContext(ctx context.Context, query string) (*Context, error)
}
```

### Layer 5: Knowledge Economy Layer

**The innovation core** - Transforms teachers into entrepreneurs

#### Components

**Resource Tracking**:
- Views: Each student view
- Downloads: Each resource download
- Bookmarks: Saved for later
- AI Citations: Used by AI tutors
- Quiz Usage: Used in assessments

**Revenue Tracking**:
- Subscription revenue (from students/institutions)
- Per-resource fees (microeconomic)
- AI usage fees (from AI tutor subscriptions)
- Institutional licensing

**Revenue Distribution**:
```
Total Revenue
    ↓
Platform Fee (20%) → Operations, hosting, AI
    ↓
Creator Pool (75%) → Distributed by engagement
    ↓
Referral Rewards (5%) → Who brought the user
```

#### Wallet System

Each contributor has a wallet:
```go
type CreatorWallet struct {
    CreatorID      string
    Balance        decimal.Decimal
    Currency       string
    Pending        decimal.Decimal  // Unverified transactions
    TotalEarned    decimal.Decimal
    WithdrawalsYTD decimal.Decimal
}
```

#### Smart Distribution

Earnings are calculated per resource:
```
Views × $0.001 (configurable per institution)
+ Downloads × $0.01
+ Bookmarks × $0.005
+ AI Citations × $0.05
+ Quiz Instances × $0.10
```

Distributed monthly via:
- Direct bank transfer
- Mobile money (MTN, Econet)
- Digital wallet (Bitcoin, Stablecoin future)

### Layer 6: Intelligence Layer

**Core Principle**: Retrieval Augmented Generation (RAG)

AI never answers from memory. Every response:
1. Retrieves relevant knowledge from DSpace
2. Generates response based on retrieved knowledge
3. Cites all sources
4. Links to original resources

#### AI Capabilities

**AI Tutor**
```
Student Question
    ↓
Knowledge Retrieval (Vector DB)
    ↓
LLM Processing
    ↓
Citation Mapping
    ↓
Response Generation
    ↓
Resource Links
```

**Quiz Generation**
- Extract concepts from resources
- Generate questions at various difficulty levels
- Validate questions for clarity and accuracy
- Track usage and effectiveness

**Summaries**
- Identify key concepts
- Generate multi-level summaries (1-min, 5-min, detailed)
- Extract key questions and answers
- Link to source resources

**Essay Feedback**
- Analyze structure and clarity
- Suggest improvements
- Check factual accuracy against resources
- Recommend relevant resources

**Recommendations**
- Based on learning patterns
- Considering curriculum progression
- Suggesting next topics
- Personalizing to learning style

#### Implementation Pattern

```go
type AIService interface {
    // Answer a question with citations
    Answer(ctx context.Context, question string) (*Answer, error)
    
    // Generate quiz from resources
    GenerateQuiz(ctx context.Context, resourceIDs []string) (*Quiz, error)
    
    // Summarize resources
    Summarize(ctx context.Context, resourceIDs []string) (*Summary, error)
    
    // Get recommendations
    GetRecommendations(ctx context.Context, studentID string) ([]*Resource, error)
}
```

### Layer 7: Experience Layer

Applications that use the platform. Each is **only a UI**.

#### ChengetAi Learn (MVP)
- Student dashboard
- AI tutor interface
- Resource discovery
- Quiz interface
- Learning progress tracking
- Offline content sync

#### Future Applications
- ChengetAi Teacher: Teacher dashboard, content management
- ChengetAi Library: Advanced search, browsing, curation
- ChengetAi Research: Academic collaboration tools
- ChengetAi Parent: Family engagement, progress tracking
- ChengetAi Skills: Professional development
- ChengetAi Admin: Institutional administration

## Communication Patterns

### Synchronous (HTTP/REST)

Services communicate via REST for immediate responses:
```
API Gateway
    ↓ (HTTP)
Platform Services
    ↓ (HTTP)
External APIs (DSpace, Payment Providers)
```

### Asynchronous (NATS)

For eventual consistency and decoupling:
```
Event Producer
    ↓ (NATS)
Event Bus
    ↓ (subscriptions)
Event Consumers
    ↓
Process & Update
```

Events:
- `user.created` → Send welcome email, create wallet
- `resource.created` → Harvest, process, index
- `quiz.completed` → Track usage, update analytics
- `payment.received` → Distribute revenue, update wallets

## Data Flow

### Resource Discovery & Processing

```
Teacher creates content in DSpace
    ↓
Knowledge Harvester polls DSpace REST API
    ↓
New resource detected
    ↓
Document Processor
    ├── Extract text
    ├── Run OCR
    └── Classify metadata
    ↓
Embedding Generator → Vector Database
    ↓
Typesense Indexer → Search Index
    ↓
Knowledge Graph Builder
    ↓
Ready for Learning & AI
```

### Student Learning

```
Student asks question
    ↓
API Gateway → AI Service
    ↓
Vector DB Search (semantic)
    ↓
Knowledge Context Retrieved
    ↓
LLM Generates Response
    ↓
Citation Mapping
    ↓
Analytics Event
    ↓
Response delivered with resource links
```

### Revenue Distribution

```
Student views resource
    ↓
Event: resource.viewed
    ↓
Analytics Service records
    ↓
Monthly batch process:
    ├── Calculate engagement metrics
    ├── Determine revenue share
    └── Update wallets
    ↓
Creator Wallet updated
    ↓
Withdrawal email sent
```

## Scalability Patterns

### Horizontal Scaling
- Each service independently scalable
- Stateless design (except search/cache)
- Load balancing via Caddy (dev) / Kubernetes (prod)

### Vertical Scaling
- PostgreSQL connection pooling
- Redis for hot data
- Typesense for search scale
- NATS for message throughput

### Caching Strategy
```
Request
    ↓
API Gateway (check Redis cache)
    ↓ miss
Service (compute result)
    ↓
Cache (Redis) + Database
    ↓
Response
```

Cache keys:
- `user:{id}` - User data
- `curriculum:{country}:{level}` - Curriculum definitions
- `search:{query}:{page}` - Search results
- `embeddings:{doc_id}` - Document embeddings

TTL strategy:
- User data: 1 hour
- Curriculum: 24 hours
- Search: 1 hour
- Embeddings: Never (vector DB is authoritative)

## Deployment Architecture

### Development (Docker Compose)

```yaml
services:
  postgres:      # Database
  redis:         # Cache
  minio:         # Object storage
  typesense:     # Search
  nats:          # Messaging
  api-gateway:   # API Gateway
  auth-service:  # Auth
  user-service:  # Users
  # ... other services
  caddy:         # Reverse proxy
```

### Production (Kubernetes)

```
Kubernetes Cluster
    ├── Namespaces
    │   ├── platform (services)
    │   └── monitoring (prometheus, grafana, loki)
    ├── StatefulSets
    │   ├── PostgreSQL
    │   ├── Redis
    │   ├── Typesense
    │   └── NATS
    ├── Deployments
    │   ├── API Gateway
    │   ├── Services (9x)
    │   └── Knowledge Layer
    ├── Services (K8s)
    │   └── Expose endpoints
    └── ConfigMaps & Secrets
        ├── Database credentials
        ├── API keys
        └── Feature flags
```

## Security Architecture

### Authentication
- JWT tokens (access + refresh)
- OAuth2 for social login
- 2FA for sensitive operations
- Session tokens in Redis

### Authorization
- Role-Based Access Control (RBAC)
- Resource-Based Access Control (RBAC)
- Policy-based access
- Audit logging of all access

### Data Protection
- TLS 1.3 for all communication
- Encryption at rest for sensitive data
- PCI-DSS for payment data
- GDPR compliance for EU users

### API Security
- Rate limiting (per user, per IP)
- CORS configuration
- CSRF protection
- API key rotation
- Request signing (future)

## Monitoring & Observability

### Metrics (Prometheus)
- Request latency (p50, p95, p99)
- Error rates
- Service availability
- Database connection pools
- Cache hit rates
- Queue depths (NATS)

### Logging (Loki)
- Structured JSON logs
- Log levels: DEBUG, INFO, WARN, ERROR
- Correlation IDs for request tracing
- Sensitive data redaction

### Tracing (OpenTelemetry)
- Distributed tracing
- Service-to-service spans
- Database query tracing
- Cache operation tracing

### Dashboards (Grafana)
- System health
- Service performance
- Learning analytics
- Revenue analytics
- Error tracking

## Testing Strategy

### Unit Tests
- Business logic only
- Mock all dependencies
- 80%+ coverage target
- Run before commit

### Integration Tests
- Service with database
- Service with external APIs (DSpace mock)
- NATS message flows
- Run in CI/CD

### Load Tests
- Simulating millions of students
- Concurrent requests per service
- Database query performance
- Cache effectiveness

### End-to-End Tests
- Full user journeys
- Cross-service workflows
- Offline scenarios
- Run nightly

## Configuration Management

### Environment Variables
```bash
# Database
DB_HOST=localhost
DB_PORT=5432
DB_NAME=chengetai
DB_USER=chengetai
DB_PASSWORD=...

# Services
NATS_URL=nats://nats:4222
REDIS_URL=redis://redis:6379
MINIO_URL=http://minio:9000

# Feature flags
FEATURE_RAG_ENABLED=true
FEATURE_VOICE_TUTOR=false

# API Keys
DSPACE_API_KEY=...
PAYMENT_API_KEY=...
```

### Feature Flags
Stored in database, checked at runtime:
```go
if featureFlags.IsEnabled("quiz_generation", userID) {
    // Show quiz generation UI
}
```

## Performance Targets

| Operation | Target | P95 | P99 |
|-----------|--------|-----|-----|
| User login | < 100ms | < 150ms | < 200ms |
| Resource search | < 200ms | < 300ms | < 500ms |
| AI answer | < 2s | < 3s | < 5s |
| Quiz generation | < 5s | < 7s | < 10s |
| Payment process | < 1s | < 2s | < 5s |
| Wallet update | < 500ms | < 750ms | < 1s |

## The Golden Rule

**DSpace is the Source of Truth**

1. Do NOT duplicate repository functionality
2. Do NOT store content in ChengetAi databases
3. Do NOT replace DSpace with custom solutions
4. Do CONSUME knowledge through DSpace APIs
5. Do EXTEND with AI processing
6. Do ENHANCE with learning interfaces

Every architectural decision must respect this principle.

---

**This architecture supports:**
- ✅ Millions of concurrent users
- ✅ Hundreds of millions of resources
- ✅ Real-time AI-powered learning
- ✅ Forever knowledge preservation
- ✅ Transparent knowledge economy
- ✅ Institutional compliance
- ✅ Geographic distribution
