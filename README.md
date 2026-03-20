# 🍽️ Restaurant Management — Production-Grade Go Backend

> **A personal note before we dive in**
>
> I started tinkering with Prometheus and Grafana one evening and ended up going down a rabbit hole — I figured if I
> was going to instrument something properly, I might as well build a project worth instrumenting. That's how this
> started: a full-fledged backend that treats observability as a first-class citizen, not an afterthought.
>
> Along the way I added **OpenTelemetry** for metrics, **testcontainers** for integration and end-to-end tests that
> spin up a real MongoDB instance, and **k6** for load testing — so every layer of the quality story is covered.
>
> Auth was an interesting design decision. **better-auth** is gaining a lot of traction, and I like it precisely
> because I build my frontends in **Next.js**. My preferred pattern is a two-hop architecture: the Next.js client
> talks to the Next.js server (API routes / Route Handlers), and *that* server calls this Go backend. The Go server
> never faces the browser directly, which means:
>
> - **No CORS configuration needed** — the internal hop is same-origin from the browser's perspective.
> - **Token caching lives in Next.js** — the Go server only needs to verify a forwarded JWT, keeping auth logic
>   lean on this side.
> - **Smaller attack surface** — the Go server is effectively an internal API.
>
> The only auth primitive this service implements is **JWT verification middleware**, which is exactly as much as it
> needs.

---

## Table of Contents

1. [Tech Stack](#tech-stack)
2. [Project Structure](#project-structure)
3. [Architecture Overview](#architecture-overview)
4. [API Reference](#api-reference)
5. [Observability](#observability)
6. [Testing](#testing)
7. [Load Testing](#load-testing)
8. [CI / CD](#ci--cd)
9. [Getting Started](#getting-started)
10. [Configuration](#configuration)
11. [Database Seeding](#database-seeding)
12. [Makefile Reference](#makefile-reference)

---

## Tech Stack

| Layer | Technology |
|---|---|
| Language | Go 1.25 |
| HTTP Framework | [Gin](https://github.com/gin-gonic/gin) |
| Database | MongoDB (via [mongo-driver v2](https://github.com/mongodb/mongo-go-driver)) |
| Config | [koanf](https://github.com/knadh/koanf) (env-based) |
| Auth | JWT verification ([golang-jwt/jwt v5](https://github.com/golang-jwt/jwt)) |
| Metrics | OpenTelemetry → Prometheus exporter |
| Tracing middleware | [otelgin](https://github.com/open-telemetry/opentelemetry-go-contrib) |
| Dashboards | Grafana (provisioned automatically) |
| Integration tests | [testcontainers-go](https://github.com/testcontainers/testcontainers-go) + MongoDB module |
| E2E tests | testcontainers-go + `net/http/httptest` |
| Load testing | [k6](https://k6.io) |
| CI / CD | GitHub Actions |
| Infra | Docker Compose |

---

## Project Structure

```
restaurant-management/
├── cmd/
│   └── main.go                  # Entry point — wires config, DB, server, telemetry
├── internal/
│   ├── config/                  # Env-based config via koanf
│   ├── database/                # MongoDB connection helper
│   ├── errs/                    # Shared error types
│   ├── handler/                 # HTTP handlers (food, menu, table, order)
│   ├── logger/                  # Structured logger setup
│   ├── middleware/               # Auth (JWT) + Logger + Recovery
│   ├── models/                  # Domain models (Food, Menu, Table, Order)
│   ├── observability/           # OpenTelemetry initialisation
│   ├── repository/              # MongoDB repository layer
│   ├── routes/                  # Route registration + Prometheus /metrics endpoint
│   ├── server/                  # HTTP server with graceful shutdown
│   └── tests/
│       ├── integration/         # Integration tests (testcontainers)
│       └── e2e/                 # E2E tests (testcontainers + httptest)
├── observability/
│   ├── prometheus.yml           # Prometheus scrape config
│   └── provisioning/
│       ├── datasources/         # Grafana datasource (Prometheus)
│       └── dashboards/          # Grafana dashboard JSON
├── scripts/
│   ├── load_test.js             # k6 load test script
│   ├── seed_data.json           # Seed data (menus, foods, tables)
│   └── seed.sh                  # CLI import script (mongoimport)
├── .github/workflows/
│   ├── ci.yml                   # CI: integration + e2e tests on PRs to master
│   └── cd.yml                   # CD: Docker build + push on merge to master
├── docker-compose.yml           # MongoDB + Prometheus + Grafana
├── Makefile                     # Developer convenience targets
└── go.mod
```

---

## Architecture Overview

```
┌──────────────────────────────────────────────────────────┐
│  Next.js Client (Browser)                                │
│  — calls Next.js API routes only, never Go directly      │
└────────────────────┬─────────────────────────────────────┘
                     │ Server-to-server (internal hop)
                     ▼
┌──────────────────────────────────────────────────────────┐
│  Next.js Server  (API Routes / Route Handlers)           │
│  — manages sessions via better-auth                      │
│  — forwards a signed JWT to the Go backend               │
└────────────────────┬─────────────────────────────────────┘
                     │ HTTP  (JWT in Authorization header)
                     ▼
┌──────────────────────────────────────────────────────────┐
│  Go Backend  :8080                                       │
│  Gin router                                              │
│  ├── /metrics        Prometheus scrape endpoint          │
│  ├── /ping           Health check                        │
│  └── /api/v1/                                            │
│       ├── menus      (public)                            │
│       ├── foods      (public)                            │
│       ├── tables     (JWT required)                      │
│       └── orders     (JWT required)                      │
└──────────┬───────────────────────────────────────────────┘
           │
           ▼
┌──────────────────┐    ┌────────────────────────────────┐
│  MongoDB :27017  │    │  Prometheus :9090              │
│  restaurant db   │    │  scrapes /metrics every 5s     │
└──────────────────┘    └────────────────┬───────────────┘
                                         │
                                         ▼
                             ┌───────────────────────┐
                             │  Grafana :3000         │
                             │  auto-provisioned      │
                             │  datasource + dashboard│
                             └───────────────────────┘
```

---

## API Reference

All routes are prefixed with `/api/v1`. Routes marked 🔒 require a `Authorization: Bearer <jwt>` header.

### Menus

| Method | Path | Description |
|---|---|---|
| `GET` | `/api/v1/menus` | List all menus |

### Foods

| Method | Path | Description |
|---|---|---|
| `GET` | `/api/v1/foods` | List all foods |
| `GET` | `/api/v1/foods/:foodId` | Get food by ID |
| `GET` | `/api/v1/foods/menu/:menuId` | Get foods by menu |

### Tables 🔒

| Method | Path | Description |
|---|---|---|
| `GET` | `/api/v1/tables` | List all tables |
| `GET` | `/api/v1/tables/:tableId` | Get table by ID |
| `POST` | `/api/v1/tables/book-table/:tableId` | Book a table |
| `POST` | `/api/v1/tables/release-table/:tableId` | Release a booked table |

### Orders 🔒

| Method | Path | Description |
|---|---|---|
| `GET` | `/api/v1/orders/:orderId` | Get order by ID |
| `GET` | `/api/v1/orders/user/:userID` | Get all orders for a user |
| `POST` | `/api/v1/orders/create-order` | Create a new order |

### Observability

| Method | Path | Description |
|---|---|---|
| `GET` | `/metrics` | Prometheus metrics |
| `GET` | `/ping` | Health check |

---

## Observability

The observability stack is fully **self-provisioning** via Docker Compose.

### How it works

1. **OpenTelemetry** is initialised at startup (`internal/observability/telemetry.go`).
2. The OTel Prometheus exporter bridges OTel metrics into the Prometheus registry.
3. **otelgin** middleware instruments every HTTP request (duration, status codes, etc.).
4. Prometheus scrapes `/metrics` every **5 seconds**.
5. Grafana auto-loads the Prometheus datasource and a pre-built dashboard on first boot.

### Start the observability stack

```bash
make observability-up   # Prometheus + Grafana only
# OR
make up                  # MongoDB + Prometheus + Grafana together
```

| Service | URL | Credentials |
|---|---|---|
| Grafana | <http://localhost:3000> | admin / admin |
| Prometheus | <http://localhost:9090> | — |

---

## Testing

Tests use **testcontainers-go** to spin up a real `mongo:7` container for each test suite — no mocking of the database layer.

### Integration Tests

Test the repository layer directly against a live MongoDB instance.

```bash
make test-integration        # go test -v
make testsum-integration     # gotestsum (prettier output)
```

### End-to-End Tests

Spin up the full server stack (`gin` router + real MongoDB via testcontainers) and exercise the HTTP API using `net/http/httptest`.

```bash
make test-e2e            # go test -v
make testsum-e2e         # gotestsum (prettier output)
```

> **Note:** Docker must be running locally for testcontainers to work.

---

## Load Testing

Uses [k6](https://k6.io). The default script ramps to **20 virtual users** over 30 s, holds for 1 min, then ramps down.

```bash
# Install k6: https://k6.io/docs/get-started/installation/
make load-test
```

Watch Grafana in real time while the load test runs to see request rates, latencies, and error rates.

---

## CI / CD

### CI (`.github/workflows/ci.yml`)

Triggers on **every pull request to `master`** and on direct pushes.

```
Checkout → Setup Go → Install gotestsum → Integration Tests → E2E Tests
```

No infrastructure secrets required — testcontainers handles MongoDB.

### CD (`.github/workflows/cd.yml`)

Triggers on **merge to `master`**.

```
Build Docker image → Push to Docker Hub (SHA + latest tags)
  → Deploy job (placeholder — wire in your own deploy command)
```

**Required repository secrets:**

| Secret | Description |
|---|---|
| `DOCKERHUB_USERNAME` | Your Docker Hub username |
| `DOCKERHUB_TOKEN` | Docker Hub access token |

---

## Getting Started

### Prerequisites

- Go 1.22+
- Docker + Docker Compose
- `mongoimport` + `jq` *(for database seeding only)*
- k6 *(for load testing only)*

### 1. Clone

```bash
git clone https://github.com/distroaryan/restaurant-management.git
cd restaurant-management
```

### 2. Start infrastructure

```bash
make up          # starts MongoDB, Prometheus, Grafana
```

### 3. Run the server

```bash
make run
# Server listening on :8080
```

### 4. (Optional) Seed the database

```bash
make seed        # imports menus, foods, and tables from scripts/seed_data.json
```

### 5. Hit the API

```bash
curl http://localhost:8080/ping
curl http://localhost:8080/api/v1/menus
```

---

## Configuration

All configuration is read from environment variables prefixed with `APP_`.

| Variable | Default | Description |
|---|---|---|
| `APP_ENV` | `development` | `development` or `production` |
| `APP_PORT` | `8080` | HTTP listen port |
| `APP_MONGO_URI` | `mongodb://127.0.0.1:27017` | MongoDB connection string |
| `APP_DB_NAME` | `restaurant` | Database name |
| `APP_JWT_SECRET` | `secret` | ⚠️ Change in production |

Example:

```bash
APP_PORT=9000 APP_MONGO_URI="mongodb://user:pass@host:27017" APP_JWT_SECRET="supersecret" make run
```

---

## Database Seeding

```bash
# Default (localhost:27017, db: restaurant)
bash scripts/seed.sh

# Custom connection
bash scripts/seed.sh --uri "mongodb://user:pass@host:27017" --db mydb

# Drop existing collections before importing
bash scripts/seed.sh --drop
```

The seed file (`scripts/seed_data.json`) ships with **10 menus**, **50 food items**, and **30 tables**.

---

## Makefile Reference

```
make run / dev           Start the server with go run
make build               Compile binary → bin/restaurant-management

make up                  Start all Docker services
make down                Stop all Docker services
make logs                Tail all container logs

make observability-up    Start Prometheus + Grafana only
make observability-down  Stop Prometheus + Grafana
make observability-logs  Tail Prometheus + Grafana logs

make test-integration    Run integration tests (go test -v)
make test-e2e            Run E2E tests (go test -v)
make testsum-integration Run integration tests (gotestsum)
make testsum-e2e         Run E2E tests (gotestsum)

make load-test           Run k6 load test
make seed                Seed MongoDB via mongoimport
```

---

## License

MIT
