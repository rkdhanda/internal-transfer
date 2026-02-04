.PHONY: help run test test-integration build migrate-up migrate-down docker-up docker-down docker-logs clean lint install-tools

# Default target
.DEFAULT_GOAL := help

## help: Display this help message
help:
	@echo "Available targets:"
	@grep -E '^## ' $(MAKEFILE_LIST) | sed 's/^## /  /'

## run: Run the application locally
run:
	@echo "Starting application..."
	go run server/main.go

docker-up:
	@echo "Starting Docker services..."
	docker-compose up -d
	@echo "Waiting for database to be ready..."
	@sleep 5
	@echo "Running migrations..."
	@docker-compose exec -T app sh -c "migrate -path ./migrations -database 'postgresql://postgres:postgres@postgres:5432/transfers?sslmode=disable' up" || true
	@echo "Services are up and running"
	@echo "API available at http://localhost:8080"

## docker-down: Stop all services and remove volumes
docker-down:
	@echo "Stopping Docker services..."
	docker-compose down -v
	@echo "Services stopped"

## migrate-up: Run database migrations up
migrate-up:
	@echo "Running migrations up..."
	migrate -path ./migrations -database "postgresql://postgres:postgres@localhost:5432/transfers?sslmode=disable" up
	@echo "Migrations completed"

## migrate-down: Run database migrations down
migrate-down:
	@echo "Running migrations down..."
	migrate -path ./migrations -database "postgresql://postgres:postgres@localhost:5432/transfers?sslmode=disable" down
	@echo "Migrations rolled back"

## migrate-create: Create a new migration (usage: make migrate-create NAME=migration_name)
migrate-create:
	@if [ -z "$(NAME)" ]; then \
		echo "Error: NAME is required. Usage: make migrate-create NAME=migration_name"; \
		exit 1; \
	fi
	@echo "Creating migration: $(NAME)"
	migrate create -ext sql -dir migrations -seq $(NAME)