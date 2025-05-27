# Read-It-Later Application Makefile

.PHONY: help dev prod build test clean extension setup logs

# Default target
help:
	@echo "Available commands:"
	@echo "  make setup     - Initial project setup"
	@echo "  make dev       - Start development environment"
	@echo "  make prod      - Start production environment"
	@echo "  make build     - Build all components"
	@echo "  make test      - Run test suite"
	@echo "  make extension - Build browser extension"
	@echo "  make logs      - Show application logs"
	@echo "  make clean     - Clean up containers and volumes"
	@echo "  make reset     - Reset everything (clean + rebuild)"

# Initial project setup
setup:
	@echo "Setting up Read-It-Later application..."
	@if [ ! -f .env ]; then cp .env.example .env; echo "Created .env file from template"; fi
	@mkdir -p data/storage logs
	@chmod +x scripts/*.sh 2>/dev/null || true
	@echo "Setup complete! Run 'make dev' to start development environment"

# Development environment
dev: setup
	@echo "Starting development environment..."
	docker-compose -f docker-compose.dev.yml up --build

# Development environment in background
dev-bg: setup
	@echo "Starting development environment in background..."
	docker-compose -f docker-compose.dev.yml up --build -d

# Production environment
prod: setup
	@echo "Starting production environment..."
	docker-compose up --build -d

# Build all components
build:
	@echo "Building all components..."
	@echo "Building backend..."
	cd backend && go build -o bin/server ./cmd/server
	@echo "Building frontend..."
	cd frontend && npm run build
	@echo "Building extension..."
	cd extension && npm run build

# Run test suite
test:
	@echo "Running test suite..."
	@echo "Testing backend..."
	cd backend && go test ./...
	@echo "Testing frontend..."
	cd frontend && npm test
	@echo "Testing extension..."
	cd extension && npm test

# Build browser extension
extension:
	@echo "Building browser extension..."
	cd extension && npm install && npm run build
	@echo "Extension built in extension/dist/"

# Show application logs
logs:
	docker-compose logs -f

# Show logs for specific service
logs-backend:
	docker-compose logs -f backend

logs-frontend:
	docker-compose logs -f frontend

logs-db:
	docker-compose logs -f postgres

# Clean up containers and volumes
clean:
	@echo "Cleaning up containers and volumes..."
	docker-compose -f docker-compose.dev.yml down -v --remove-orphans
	docker-compose down -v --remove-orphans
	docker system prune -f

# Reset everything
reset: clean
	@echo "Resetting everything..."
	docker-compose -f docker-compose.dev.yml build --no-cache
	docker-compose build --no-cache

# Database operations
db-migrate:
	@echo "Running database migrations..."
	cd backend && go run cmd/migrate/main.go

db-seed:
	@echo "Seeding database..."
	cd backend && go run cmd/seed/main.go

db-reset:
	@echo "Resetting database..."
	docker-compose exec postgres psql -U postgres -d readitlater -c "DROP SCHEMA public CASCADE; CREATE SCHEMA public;"
	$(MAKE) db-migrate

# Development helpers
install-deps:
	@echo "Installing dependencies..."
	cd backend && go mod download
	cd frontend && npm install
	cd extension && npm install

format:
	@echo "Formatting code..."
	cd backend && go fmt ./...
	cd frontend && npm run format
	cd extension && npm run format

lint:
	@echo "Linting code..."
	cd backend && golangci-lint run
	cd frontend && npm run lint
	cd extension && npm run lint

# Security checks
security-check:
	@echo "Running security checks..."
	cd backend && gosec ./...
	cd frontend && npm audit
	cd extension && npm audit

# Backup and restore
backup:
	@echo "Creating backup..."
	docker-compose exec postgres pg_dump -U postgres readitlater > backup_$(shell date +%Y%m%d_%H%M%S).sql

restore:
	@echo "Restoring from backup..."
	@read -p "Enter backup file name: " backup_file; \
	docker-compose exec -T postgres psql -U postgres readitlater < $$backup_file

# Performance monitoring
monitor:
	@echo "Starting performance monitoring..."
	docker stats

# Update dependencies
update-deps:
	@echo "Updating dependencies..."
	cd backend && go get -u ./...
	cd frontend && npm update
	cd extension && npm update 