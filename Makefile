.PHONY: test-integration test-e2e testsum-integration testsum-e2e load-test \
        run build dev \
        up down logs \
        observability-up observability-down observability-logs

# ─── Server ──────────────────────────────────────────────────────────────────

## Build the Go binary
build:
	go build -o bin/restaurant-management ./cmd/...

## Run the server directly with `go run`
run:
	go run ./cmd/...

## Alias: same as `run` (useful muscle-memory target)
dev: run

# ─── Docker / Infrastructure ─────────────────────────────────────────────────

## Start all infra services (MongoDB, Prometheus, Grafana) in the background
up:
	docker compose up -d

## Stop all infra services
down:
	docker compose down

## Tail logs for all infra containers
logs:
	docker compose logs -f

# ─── Observability ───────────────────────────────────────────────────────────

## Start ONLY the observability stack (Prometheus + Grafana)
observability-up:
	docker compose up -d prometheus grafana

## Stop ONLY the observability stack
observability-down:
	docker compose stop prometheus grafana

## Tail logs for Prometheus and Grafana only
observability-logs:
	docker compose logs -f prometheus grafana

# ─── Tests ───────────────────────────────────────────────────────────────────

test-integration:
	go test -v ./internal/tests/integration/...

test-e2e:
	go test -v ./internal/tests/e2e/...

testsum-integration:
	gotestsum --format testname -- -v ./internal/tests/integration/...

testsum-e2e:
	gotestsum --format testname -- -v ./internal/tests/e2e/...

# ─── Load Test ───────────────────────────────────────────────────────────────

load-test:
	k6 run scripts/load_test.js

# ─── Database Seeding ────────────────────────────────────────────────────────

## Seed MongoDB using mongoimport (requires mongoimport in PATH and MongoDB running)
seed:
	bash scripts/seed.sh
