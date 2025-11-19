.PHONY: help setup up down restart test clean build

help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

setup: ## Initial setup - copy .env.example to .env if needed
	@if [ ! -f .env ]; then \
		echo "Creating .env from .env.example..."; \
		cp .env.example .env; \
		echo "Please edit .env file with your configuration"; \
	else \
		echo ".env file already exists"; \
	fi

up: ## Start all services (setup runs automatically)
	docker-compose up -d

build: ## Start all services (setup runs automatically)
	docker-compose up -d --build

down: ## Stop all services
	docker-compose down

restart: ## Restart all services
	docker-compose restart

test: ## Run integration tests
	python3 tests/integration_test.py

clean: ## Stop services and remove volumes
	docker-compose down -v

logs: ## Show logs from all services
	docker-compose logs -f

ps: ## Show status of all services
	docker-compose ps

.PHONY: swagger_generate
swagger_generate:
	./scripts/generate_from_swagger.sh

.PHONY: grpc_generate
grpc_generate:
	protoc -I api/proto api/proto/*.proto \
	  --go_out=parking/internal/grpc/gen \
	  --go_opt=paths=source_relative \
	  --go_opt=Mparking.proto=github.com/h4x4d/parking_net/parking/internal/grpc/gen \
	  --go-grpc_out=parking/internal/grpc/gen \
	  --go-grpc_opt=paths=source_relative \
	  --go-grpc_opt=Mparking.proto=github.com/h4x4d/parking_net/parking/internal/grpc/gen

	protoc -I api/proto api/proto/*.proto \
	  --go_out=booking/internal/grpc/gen \
	  --go_opt=paths=source_relative \
	  --go_opt=Mparking.proto=github.com/h4x4d/parking_net/booking/internal/grpc/gen \
	  --go-grpc_out=booking/internal/grpc/gen \
	  --go-grpc_opt=paths=source_relative \
	  --go-grpc_opt=Mparking.proto=github.com/h4x4d/parking_net/booking/internal/grpc/gen

.PHONY: codegen
codegen: grpc_generate swagger_generate
