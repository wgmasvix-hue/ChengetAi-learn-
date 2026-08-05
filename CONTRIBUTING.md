# Contributing to ChengetAi

We're thrilled you want to contribute to the digital knowledge infrastructure for Africa! This document provides guidelines for contributing to the ChengetAi Platform.

## Code of Conduct

Be respectful, inclusive, and collaborative. We're building for millions of African learners and educators.

## Before You Start

1. Read [VISION.md](docs/VISION.md) to understand our mission
2. Read [ARCHITECTURE.md](docs/ARCHITECTURE.md) to understand our system design
3. Read [KNOWLEDGE_ECONOMY.md](docs/KNOWLEDGE_ECONOMY.md) to understand creator rewards
4. Check [GitHub Issues](../../issues) to see what needs work
5. Join our discussions in the community forum

## Architecture Rules

**Always follow these rules:**

1. **DSpace is the Source of Truth**
   - Never duplicate repository functionality
   - Never store content in ChengetAi databases
   - Always retrieve knowledge through DSpace APIs
   - Extend and enhance, never replace

2. **Clean Architecture**
   - Separate concerns: domain, services, handlers, repositories
   - Dependency injection, not tight coupling
   - Hexagonal architecture patterns
   - Test-driven development

3. **Independently Deployable Services**
   - Each service runs in its own container
   - Services communicate via REST or NATS
   - No shared databases between services
   - Configuration via environment variables

4. **Scalability First**
   - Think in terms of millions of users
   - Design for horizontal scaling
   - Use caching strategically
   - Avoid monoliths and tight coupling

## Getting Started

### Development Environment

```bash
# Clone the repository
git clone https://github.com/wgmasvix-hue/ChengetAi-learn-.git
cd ChengetAi-learn-

# Copy environment template
cp .env.example .env

# Start the platform with Docker Compose
docker-compose up -d

# Verify all services are healthy
docker-compose ps

# View logs
docker-compose logs -f api-gateway
```

### First Contribution

1. **Pick an issue** from the [Issues](../../issues) page
2. **Comment on the issue** saying you'll work on it
3. **Create a feature branch**: `git checkout -b feature/your-feature`
4. **Follow the development guidelines** below
5. **Push your changes**: `git push origin feature/your-feature`
6. **Create a Pull Request** with a clear description

## Development Workflow

### File Organization

Each service follows this structure:

```
apps/SERVICE_NAME/
├── cmd/server/
│   └── main.go                 # Entry point
├── internal/
│   ├── domain/                 # Domain models
│   ├── services/               # Business logic
│   ├── repositories/           # Data access
│   ├── handlers/               # HTTP handlers
│   ├── middleware/             # Middleware
│   └── config/                 # Configuration
├── migrations/                 # Database migrations
├── tests/                      # Integration tests
├── Dockerfile                  # Container definition
├── Makefile                    # Build commands
├── README.md                   # Service documentation
└── go.mod                      # Go module

platform/
├── pkg/                        # Shared packages
├── internal/                   # Shared internals
├── configs/                    # Platform configs
└── database/                   # Shared migrations
```

### Code Standards

#### Go Code

```go
// Good: Clear, well-documented domain code
type UserService interface {
    // CreateUser creates a new user and assigns a wallet.
    // Returns ErrEmailTaken if email already exists.
    CreateUser(ctx context.Context, user *User) error
}

// Bad: Vague, missing documentation
type Service interface {
    DoThing(x interface{}) interface{}
}
```

#### Comments
- Use GoDoc comment style: `// FunctionName does X.`
- Explain the "why", not just the "what"
- Document error cases and side effects
- Update comments when you change code

#### Testing
- Minimum 80% coverage for business logic
- Write unit tests first (TDD preferred)
- Use table-driven tests for multiple scenarios
- Mock external dependencies
- Name tests clearly: `TestFunctionName_Scenario_Expected`

#### Error Handling
- Never ignore errors (no `_ = err`)
- Wrap errors with context: `fmt.Errorf("creating user: %w", err)`
- Use custom error types for domain errors
- Return early to avoid nesting

### Running Tests

```bash
# Test a single service
cd apps/api-gateway
go test -v -cover ./...

# Test with coverage report
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out

# Run integration tests
cd ../../tests
go test -v ./...
```

### Building

```bash
# Build a service
cd apps/SERVICE_NAME
make build

# Build all services with Docker
docker-compose build

# Run a service locally
make run
```

### Database Migrations

```bash
# Create a new migration
# File: platform/database/migrations/XXX_description.sql

-- migrations/001_create_users_table.sql
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email VARCHAR(255) UNIQUE NOT NULL,
    name VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_users_email ON users(email);
```

Migrations are automatically applied on service startup.

## Commit Messages

Follow [Conventional Commits](https://www.conventionalcommits.org/):

```
type(scope): subject

description

Closes #123
```

Types: `feat`, `fix`, `docs`, `style`, `refactor`, `perf`, `test`, `chore`

Examples:
```
feat(auth): implement JWT token refresh
fix(user-service): correct email validation regex
docs(architecture): add service communication diagram
```

## Pull Request Process

1. **Update documentation** if you change behavior
2. **Update tests** - add tests for new features
3. **Verify all checks pass**: `make lint test`
4. **Write a clear PR description** explaining:
   - What problem does this solve?
   - How does it solve it?
   - Are there any breaking changes?
   - Testing instructions

### PR Template

```markdown
## Description
Brief description of what this PR does.

## Related Issue
Closes #123

## Changes
- Change 1
- Change 2

## Testing
How to test these changes:
1. Start the platform
2. Navigate to X
3. Verify Y happens

## Checklist
- [ ] Tests pass
- [ ] Documentation updated
- [ ] No breaking changes
- [ ] Follows architecture rules
```

## Architecture Review Checklist

Before submitting a PR, ensure:

- [ ] **Respects Clean Architecture**
  - [ ] Domain logic isolated from frameworks
  - [ ] Dependencies point inward
  - [ ] Entities don't depend on use cases

- [ ] **Follows SOLID Principles**
  - [ ] Single responsibility
  - [ ] Open for extension, closed for modification
  - [ ] Liskov substitution possible
  - [ ] Interface segregation
  - [ ] Depends on abstractions, not concretions

- [ ] **Maintains Scalability**
  - [ ] No N+1 queries
  - [ ] Efficient algorithms
  - [ ] Proper indexing
  - [ ] Caching where appropriate

- [ ] **Preserves DSpace as Source of Truth**
  - [ ] No repository functionality duplication
  - [ ] Uses DSpace APIs properly
  - [ ] Doesn't bypass repository
  - [ ] Respects versioning and identifiers

- [ ] **Production Ready**
  - [ ] Error handling complete
  - [ ] Logging sufficient
  - [ ] Metrics exposed
  - [ ] Health checks implemented

## Performance Targets

Code contributions must meet these targets:

| Operation | Target | P95 | P99 |
|-----------|--------|-----|-----|
| User login | < 100ms | < 150ms | < 200ms |
| Search | < 200ms | < 300ms | < 500ms |
| AI answer | < 2s | < 3s | < 5s |
| Quiz gen | < 5s | < 7s | < 10s |

Measure performance: `go test -bench ./... -benchmem`

## Documentation

- Update README.md if you change user-facing behavior
- Document new environment variables in .env.example
- Update API documentation in docs/
- Add architecture decisions to docs/ARCHITECTURE.md

## Questions?

- Open an issue with the `question` label
- Join our community discussions
- Ask in PRs if you need clarification

## Recognition

Contributors are recognized in:
- [CONTRIBUTORS.md](CONTRIBUTORS.md) - Hall of fame
- Release notes - For significant contributions
- Knowledge economy rewards - In the future, contributors will earn from platform usage

## Code of Conduct Summary

- Be respectful and inclusive
- Assume good intent
- Welcome diverse perspectives
- Focus on ideas, not individuals
- Avoid harassment and discrimination

---

**Thank you for building the future of African education!** 🚀

Every contribution brings us closer to our mission: Preserve African Knowledge, Democratise Education, Reward Knowledge Creators.
