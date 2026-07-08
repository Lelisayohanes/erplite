# ──────────────────────────────────────────────────────────────
# ERPLite – Development Automation
# ──────────────────────────────────────────────────────────────

APP_NAME  := erplite
BUILD_DIR := bin
CMD_DIR   := ./cmd/server
BINARY    := $(BUILD_DIR)/$(APP_NAME)

.DEFAULT_GOAL := help

## build: Compile the application binary
.PHONY: build
build:
	@echo "▸ Building $(APP_NAME)..."
	@mkdir -p $(BUILD_DIR)
	go build -o $(BINARY) $(CMD_DIR)

## test: Run all automated tests
.PHONY: test
test:
	@echo "▸ Running tests..."
	go test -v -race -count=1 ./...

## bdd: Run Gherkin BDD feature tests with godog
.PHONY: bdd
bdd:
	@echo "▸ Running BDD features..."
	go test -v ./features/...

## lint: Validate code quality with go vet
.PHONY: lint
lint:
	@echo "▸ Linting..."
	go vet ./...

## migrate: Apply all pending database migrations (direction=up)
.PHONY: migrate
migrate:
	@echo "▸ Running migrations..."
	go run $(CMD_DIR) migrate up

## migrate-down: Roll back all database migrations
.PHONY: migrate-down
migrate-down:
	@echo "▸ Rolling back migrations..."
	go run $(CMD_DIR) migrate down

## run: Start the application with go run
.PHONY: run
run:
	go run $(CMD_DIR)

## clean: Remove build artefacts
.PHONY: clean
clean:
	@echo "▸ Cleaning..."
	@rm -rf $(BUILD_DIR)

## help: Display this help message
.PHONY: help
help:
	@echo "Usage: make <target>"
	@echo ""
	@sed -n 's/^## //p' $(MAKEFILE_LIST) | column -t -s ':' | sed 's/^/  /'
