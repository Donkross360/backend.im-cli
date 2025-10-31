.PHONY: help mock-api cli build test clean

help:
	@echo "Backend.im CLI - Makefile Commands"
	@echo ""
	@echo "  make mock-api      - Run mock API server"
	@echo "  make cli           - Build CLI tool"
	@echo "  make build         - Build everything"
	@echo "  make test          - Test mock API"
	@echo "  make clean         - Clean build artifacts"
	@echo "  make docker-build  - Build Docker images"
	@echo "  make docker-up     - Start Docker services"
	@echo "  make docker-down   - Stop Docker services"

mock-api:
	@echo "🚀 Starting mock API server..."
	@cd mock-api && docker-compose up mock-api

cli:
	@echo "🔨 Building CLI..."
	@cd cli && docker build -t backend-im-cli .
	@docker run --rm -v $(PWD)/cli:/app backend-im-cli go build -o backend-im ./cmd/backend-im

build: mock-api cli

test:
	@echo "🧪 Testing mock API..."
	@curl -X POST http://localhost:8080/api/generate \
		-H "Content-Type: application/json" \
		-d '{"prompt": "Create a REST API"}' || echo "Mock API not running. Start with: make docker-up"

clean:
	@echo "🧹 Cleaning..."
	@rm -f cli/backend-im
	@docker-compose down -v 2>/dev/null || true

docker-build:
	@echo "🐳 Building Docker images..."
	@docker-compose build

docker-up:
	@echo "🚀 Starting Docker services..."
	@docker-compose up -d mock-api
	@echo "✅ Mock API running on http://localhost:8080"
	@echo "📡 WebSocket: ws://localhost:8080/ws"

docker-down:
	@echo "🛑 Stopping Docker services..."
	@docker-compose down

docker-logs:
	@docker-compose logs -f mock-api

# Run CLI commands via Docker
cli-help:
	@docker-compose run --rm cli --help

cli-auth:
	@docker-compose run --rm cli auth

cli-generate:
	@docker-compose run --rm cli generate "Create a REST API"

cli-deploy:
	@docker-compose run --rm cli deploy --project my-project

# Shortcuts for common deploy operations
deploy:
	@if [ -z "$(PROJECT)" ]; then \
		echo "Usage: make deploy PROJECT=<project-id> [DIR=<directory>]"; \
		exit 1; \
	fi; \
	DIR=$${DIR:-.}; \
	docker-compose run --rm -v $$(pwd)/$$DIR:/workspace/$$DIR cli deploy $$PROJECT --dir /workspace/$$DIR --watch

deploy-no-watch:
	@if [ -z "$(PROJECT)" ]; then \
		echo "Usage: make deploy-no-watch PROJECT=<project-id> [DIR=<directory>]"; \
		exit 1; \
	fi; \
	DIR=$${DIR:-.}; \
	docker-compose run --rm -v $$(pwd)/$$DIR:/workspace/$$DIR cli deploy $$PROJECT --dir /workspace/$$DIR

