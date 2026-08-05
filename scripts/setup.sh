#!/bin/bash
set -e

echo "🚀 ChengetAi Platform Setup"
echo "================================"
echo ""

# Check prerequisites
echo "Checking prerequisites..."

if ! command -v docker &> /dev/null; then
    echo "❌ Docker is not installed"
    echo "Please install Docker from https://docker.com"
    exit 1
fi

if ! command -v docker-compose &> /dev/null; then
    echo "❌ Docker Compose is not installed"
    echo "Please install Docker Compose from https://docs.docker.com/compose/install/"
    exit 1
fi

if ! command -v go &> /dev/null; then
    echo "❌ Go is not installed"
    echo "Please install Go from https://golang.org/doc/install"
    exit 1
fi

echo "✓ All prerequisites installed"
echo ""

# Create .env file if it doesn't exist
if [ ! -f .env ]; then
    echo "Creating .env file..."
    cp .env.example .env
    echo "✓ .env file created"
    echo "⚠️  Please update .env with your configuration"
else
    echo "✓ .env file already exists"
fi

echo ""

# Download Go dependencies
echo "Downloading Go dependencies..."
go mod download
go mod tidy
echo "✓ Dependencies downloaded"

echo ""

# Build Docker images
echo "Building Docker images (this may take a few minutes)..."
docker-compose build
echo "✓ Docker images built"

echo ""
echo "================================"
echo "✓ Setup complete!"
echo ""
echo "Next steps:"
echo "1. Review and update .env if needed"
echo "2. Run: make dev"
echo "3. Visit http://localhost:8000/health"
echo ""
echo "Documentation:"
echo "- Vision:     docs/VISION.md"
echo "- Architecture: docs/ARCHITECTURE.md"
echo "- Economy:    docs/KNOWLEDGE_ECONOMY.md"
echo "- Contributing: CONTRIBUTING.md"
echo ""
