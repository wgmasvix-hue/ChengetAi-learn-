.PHONY: build run test clean lint format

SERVICE_NAME=school-service
SERVICE_PATH=./apps/$(SERVICE_NAME)
BINARY=$(SERVICE_PATH)/bin/$(SERVICE_NAME)

build:
	@echo "Building $(SERVICE_NAME)..."
	@mkdir -p $(SERVICE_PATH)/bin
	@cd .. && go build -o $(SERVICE_NAME)/bin/$(SERVICE_NAME) ./apps/$(SERVICE_NAME)/cmd/server

run:
	@echo "Running $(SERVICE_NAME)..."
	@$(BINARY)

test:
	@echo "Testing $(SERVICE_NAME)..."
	@go test -v -cover ./...

lint:
	@echo "Linting $(SERVICE_NAME)..."
	@go vet ./...
	@go fmt ./...

clean:
	@echo "Cleaning $(SERVICE_NAME)..."
	@rm -rf bin
	@go clean

docker-build:
	@echo "Building Docker image for $(SERVICE_NAME)..."
	@docker build -t chengetai-$(SERVICE_NAME):latest \
		--build-arg SERVICE_NAME=$(SERVICE_NAME) \
		-f ../../docker/Dockerfile.service ../../
