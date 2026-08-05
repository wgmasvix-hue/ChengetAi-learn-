.PHONY: build run test clean lint

SERVICE_NAME=api-gateway
SERVICE_PATH=./apps/$(SERVICE_NAME)
BINARY=$(SERVICE_PATH)/bin/$(SERVICE_NAME)

build:
	@echo "Building $(SERVICE_NAME)..."
	@mkdir -p $(SERVICE_PATH)/bin
	@cd $(SERVICE_PATH) && go build -o bin/$(SERVICE_NAME) ./cmd/server

run:
	@echo "Running $(SERVICE_NAME)..."
	@$(BINARY)

test:
	@echo "Testing $(SERVICE_NAME)..."
	@cd $(SERVICE_PATH) && go test -v -cover ./...

lint:
	@echo "Linting $(SERVICE_NAME)..."
	@cd $(SERVICE_PATH) && go vet ./...
	@cd $(SERVICE_PATH) && go fmt ./...

clean:
	@echo "Cleaning $(SERVICE_NAME)..."
	@rm -rf $(SERVICE_PATH)/bin
	@cd $(SERVICE_PATH) && go clean

install-deps:
	@echo "Installing dependencies for $(SERVICE_NAME)..."
	@cd $(SERVICE_PATH) && go mod download
	@cd $(SERVICE_PATH) && go mod tidy

docker-build:
	@echo "Building Docker image for $(SERVICE_NAME)..."
	@docker build -t chengetai-$(SERVICE_NAME):latest \
		--build-arg SERVICE_NAME=$(SERVICE_NAME) \
		-f docker/Dockerfile.service .

docker-run:
	@echo "Running Docker container for $(SERVICE_NAME)..."
	@docker run --rm \
		-e SERVICE_NAME=$(SERVICE_NAME) \
		-e SERVICE_PORT=8000 \
		-p 8000:8000 \
		chengetai-$(SERVICE_NAME):latest
