-include .env
export

help:
	@echo ''
	@echo 'Usage: make [TARGET] [EXTRA_ARGUMENTS]'
	@echo 'Targets:'
	@echo '  make dev           - development (docker-compose up)'
	@echo '  make build        - build container'
	@echo '  make production   - docker production build'
	@echo '  make clean        - remove docker volumes/images'
	@echo '  make run          - run application (go run main.go)'
	@echo '  make test         - run unit tests'
	@echo '  make test-coverage - run tests with coverage report'
	@echo '  make migrate-up   - apply all pending SQL migrations'
	@echo '  make migrate-down - rollback the last SQL migration'
	@echo '  make migrate-version - print current migration version'
	@echo '  make migrate-create NAME=xxx - scaffold a new migration pair'
	@echo ''

MIGRATE_DB_URL = postgres://$(MASTER_DB_USER):$(MASTER_DB_PASSWORD)@$(MASTER_DB_HOST):$(MASTER_DB_PORT)/$(MASTER_DB_NAME)?sslmode=$(or $(MASTER_SSL_MODE),disable)
MIGRATE_PATH = internal/adapters/database/migrations/sql

dev:
	if [ ! -f .env ]; then cp .env.example .env; fi;
	docker-compose -f .docker/docker-compose-dev.yml down
	docker-compose -f .docker/docker-compose-dev.yml up

build:
	docker-compose -f .docker/docker-compose-prod.yml build
	docker-compose -f .docker/docker-compose-dev.yml down build

production:
	docker-compose -f .docker/docker-compose-prod.yml up -d --build

clean:
	docker-compose -f .docker/docker-compose-prod.yml down -v
	docker-compose -f .docker/docker-compose-dev.yml down -v

run:
	go run main.go

test:
	go test ./tests/unit/... ./internal/... ./pkg/...

test-coverage:
	go test -coverprofile=coverage.out ./tests/unit/... ./internal/... ./pkg/...
	go tool cover -func=coverage.out

# test-coverage-check runs tests with coverage and fails if total coverage is below 70%.
# Requires a Go toolchain that supports merged coverage (see tests/README.md).
test-coverage-check:
	go test -coverprofile=coverage.out ./tests/unit/... ./internal/... ./pkg/...
	@go tool cover -func=coverage.out | grep '^total:' | awk '{gsub(/%/,""); if ($$3 < 70) { print "Coverage " $$3 "% is below 70%"; exit 1 } }'
	go tool cover -func=coverage.out

# migrate-* targets require the golang-migrate CLI (see internal/adapters/database/migrations/sql/README.md).
# The app applies pending migrations automatically on startup (fail-fast);
# these targets are for manual rollback and scaffolding new migrations.
migrate-up:
	migrate -database "$(MIGRATE_DB_URL)" -path $(MIGRATE_PATH) up

migrate-down:
	migrate -database "$(MIGRATE_DB_URL)" -path $(MIGRATE_PATH) down 1

migrate-version:
	migrate -database "$(MIGRATE_DB_URL)" -path $(MIGRATE_PATH) version

migrate-create:
	@if [ -z "$(NAME)" ]; then echo "Usage: make migrate-create NAME=add_something"; exit 1; fi
	migrate create -ext sql -dir $(MIGRATE_PATH) -seq $(NAME)
