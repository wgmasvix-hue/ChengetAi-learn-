# ChengetAi Platform

> Preserve African Knowledge | Democratise Education | Reward Knowledge Creators

ChengetAi is Africa's AI Knowledge Platform - a complete digital knowledge infrastructure serving millions of learners, thousands of institutions, and hundreds of millions of educational resources.

**Chengeta** (isiZulu) means: To Preserve, To Protect, To Keep Safe

## Vision

We are building an ecosystem, not an application. Everything starts with knowledge. Everything ends with knowledge.

### Core Principles

- **Knowledge should not disappear** - We preserve African educational resources forever
- **Teachers should own their knowledge** - Content creators maintain ownership and earn from their contributions
- **Students should discover trusted knowledge** - Learning happens through verified, quality resources
- **AI should learn from trusted knowledge** - All AI capabilities use Retrieval Augmented Generation
- **Creators should earn from their knowledge** - A transparent knowledge economy rewards contribution
- **Institutions should preserve knowledge forever** - DSpace is the source of truth for all educational content

## Architecture Overview

The ChengetAi Platform consists of **seven independent layers**:

### Layer 1: Infrastructure
- Ubuntu, Docker, Docker Compose
- PostgreSQL, Redis, MinIO, Typesense, NATS
- Monitoring: Prometheus, Grafana, Loki
- Tracing: OpenTelemetry

### Layer 2: Platform Services
Shared microservices used by all applications:
- API Gateway
- Authentication & Authorization
- Users, Schools, Organizations
- Wallet & Payments
- Notifications & Analytics
- Search & Audit Logs
- Settings & Feature Flags

### Layer 3: Discovery Layer
- **Source of Truth**: DSpace at repo.dare.co.zw
- Communities | Collections | Items | Metadata | Bitstreams
- REST API | OAI-PMH | SWORD
- Persistent Identifiers & Versioning

### Layer 4: Knowledge Layer (Go)
Intelligent knowledge processing:
- Harvest repository content
- Extract metadata & text
- OCR documents
- Generate embeddings & semantic indexes
- Build knowledge graphs
- Prepare content for AI

### Layer 5: Knowledge Economy Layer
The innovative core of ChengetAi:
- Teachers contribute to DSpace Communities (not the app)
- Every resource becomes: Permanent | Discoverable | Citable | Searchable | AI-Ready | Monetisable
- Automatic revenue distribution to creators
- Teachers become knowledge entrepreneurs

### Layer 6: Intelligence Layer
AI capabilities grounded in trusted knowledge:
- AI Tutor & Question Answering
- Quiz Generation & Summaries
- Essay Feedback & Lesson Planning
- Voice Tutor & Recommendations
- Translation & Study Planning
- **Always cite repository resources via RAG**

### Layer 7: Experience Layer
User-facing applications:
- ChengetAi Learn (Zimbabwe education)
- ChengetAi Teacher
- ChengetAi Library
- ChengetAi Research
- ChengetAi Parent
- ChengetAi Skills
- ChengetAi Admin

Each application is **only a UI**. Business logic lives in platform services.

## Technology Stack

| Component | Technology |
|-----------|-----------|
| Language | Go 1.25+ |
| Frontend | Flutter |
| Database | PostgreSQL |
| Caching | Redis |
| Object Storage | MinIO |
| Search | Typesense |
| Messaging | NATS |
| Authentication | JWT + OAuth2 |
| Containerisation | Docker |
| Orchestration | Docker Compose |
| CI/CD | GitHub Actions |
| Monitoring | Prometheus + Grafana |
| Logging | Loki |
| Tracing | OpenTelemetry |

## Architecture Rules

Every service follows:
- Clean Architecture
- Domain-Driven Design (DDD)
- SOLID Principles
- CQRS where appropriate
- Repository Pattern
- Dependency Injection
- Hexagonal Architecture
- Context everywhere
- Structured logging
- Graceful shutdown
- Health checks & metrics
- Feature flags
- Environment-based configuration

## Getting Started

```bash
# Clone the repository
git clone https://github.com/wgmasvix-hue/ChengetAi-learn-.git
cd ChengetAi-learn-

# Start the platform with Docker Compose
docker-compose up -d

# View the full architecture
cat docs/Architecture.md
```

## Repository Structure

```
chengetai/
├── apps/                          # Microservices
│   ├── api-gateway/              # API Gateway Service
│   ├── auth-service/             # Authentication & Authorization
│   ├── user-service/             # User Management
│   ├── school-service/           # School Management
│   ├── wallet-service/           # Knowledge Economy Wallet
│   ├── payment-service/          # Payment Processing
│   ├── notification-service/     # Notifications
│   ├── analytics-service/        # Analytics Engine
│   └── search-service/           # Search Service
├── platform/
│   ├── pkg/                      # Shared packages
│   ├── internal/                 # Platform internals
│   ├── configs/                  # Configuration files
│   └── database/                 # Database migrations
├── docker/                        # Docker configurations
├── deployments/                   # Deployment scripts
├── scripts/                       # Utility scripts
├── mobile/                        # Flutter mobile app
├── web/                          # Web applications
├── tests/                        # Integration tests
├── docs/                         # Documentation
├── docker-compose.yml            # Development environment
└── go.mod                        # Go module
```

## Documentation

- [Architecture.md](docs/Architecture.md) - Detailed architectural decisions
- [Vision.md](docs/Vision.md) - Long-term vision and roadmap
- [KnowledgeEconomy.md](docs/KnowledgeEconomy.md) - Knowledge economy design
- [DiscoveryLayer.md](docs/DiscoveryLayer.md) - DSpace integration
- [AI.md](docs/AI.md) - AI and RAG strategy
- [ContributorGuide.md](docs/ContributorGuide.md) - How to contribute
- [API.md](docs/API.md) - API specifications
- [Deployment.md](docs/Deployment.md) - Deployment guide
- [Docker.md](docs/Docker.md) - Docker setup
- [Security.md](docs/Security.md) - Security guidelines
- [Roadmap.md](docs/Roadmap.md) - Development roadmap

## Development Guidelines

1. **Never implement large features in one step** - Work incrementally
2. **Always compile after changes** - Verify code builds
3. **Always create tests** - Maintain quality
4. **Always document** - Knowledge is power
5. **Never break architecture** - DSpace is the source of truth
6. **Never duplicate logic** - Use shared packages
7. **When uncertain, choose scalability** - Think long-term

## The Golden Rule

> **DSpace is the Source of Truth**
> 
> Nothing duplicates repository functionality.
> Everything consumes repository knowledge.

## Key Metrics

ChengetAi must support:
- **Millions** of learners
- **Thousands** of institutions
- **Hundreds of millions** of educational resources
- **Real-time** AI-powered learning
- **Forever** knowledge preservation
- **Transparent** knowledge economy

## Support

This is built by Africans, for Africans. Every decision reinforces the vision of democratised education across the continent.

## License

TBD - Knowledge for Africa

---

**Built with ❤️ for Africa's Future**
