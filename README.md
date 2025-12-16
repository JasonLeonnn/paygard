# Jalytics

![Go](https://img.shields.io/badge/Go-1.24+-00ADD8?logo=go)
![PostgreSQL](https://img.shields.io/badge/PostgreSQL-15+-336791?logo=postgresql)
![Prometheus](https://img.shields.io/badge/Prometheus-Metrics-E6522C?logo=prometheus)
![Grafana](https://img.shields.io/badge/Grafana-Dashboard-F46800?logo=grafana)
![License](https://img.shields.io/badge/license-MIT-green)

**Jalytics** is a Go-based backend system for **financial transaction monitoring and anomaly detection**, designed like a production-ready fintech system. It features real-time metrics, alerts, and visual dashboards with Prometheus and Grafana.

## 💡 Motivation

I wanted to understand how production backends handle real-time monitoring and alerting, while creating a scalable Go backend system.
Financial fraud detection seemed like a concrete domain to explore statistical anomaly detection with meaningful observability.

---

## 🚀 Features

### Core

* Create financial transactions via REST API
* Detect anomalous transactions using statistical baselines
* Persist alerts for suspicious activity

### Observability

* Prometheus metrics for:
  * HTTP latency
  * Database query latency
  * Transaction volume
  * Anomaly count
* `/metrics` endpoint for scraping
* Grafana dashboard for visualization
![Grafana Overview](docs/img/grafana-overview.png)
* Alert rules for anomaly rate and latency (Prometheus Alerting)

### Reliability & Safety

* Token bucket rate limiting middleware
* Context-aware DB queries
* Graceful shutdown with signal handling

---

## 🧠 Architecture Overview

```
┌──────────┐    HTTP     ┌────────────┐    calls    ┌──────────────┐
│  Client  │ ─────────▶ │  Handlers  │ ─────────▶ │   Services   │
└──────────┘             └────────────┘             └──────────────┘
                                                              │
                                                              ▼
                                                      ┌──────────────┐
                                                      │   Database   │
                                                      └──────────────┘
```

Cross-cutting concerns:

* **Middleware**: rate limiting, metrics
* **Metrics**: Prometheus histograms & counters
* **Alerts**: persisted anomaly signals

---

## 📡 API Endpoints

### Create Transaction

`POST /transactions`

```json
{
  "amount": 10000,
  "category": "food",
  "merchant": "Starbucks",
  "transaction_date": "2025-12-14T12:00:00Z"
}
```

Response:

```json
{
  "id": 42,
  "anomaly": true
}
```

---

### List Alerts

`GET /alerts?limit=50`

Returns latest detected anomalies, including severity:

[
  {
    "id": "b63e6a61-0bca-4f27-9e9f-7d1b02b2b3b5",
    "transaction_id": "e1f32b2c-5a6f-4f93-8e90-9d5a1e8f4d20",
    "alert_message": "Amount exceeds 3σ baseline",
    "severity": "high",
    "created_at": "2025-12-15T12:34:56Z"
  }
]

---

### Health Check

`GET /health`

Returns `200 OK` if service is running.

---

### Metrics

`GET /metrics`

Exposes Prometheus metrics:

* `transactions_total{category}` -> Total number transactions per category.
* `anomalies_total` -> Total detected anomalies.
* `anomalies_by_severity_total{severity}` -> Total anomalies by severity.
* `http_request_duration_seconds` -> HTTP latency.
* `db_query_duration_seconds` -> Database query latency.

These metrics can be visualized in **Grafana dashboards**.

---

## 🎯 Anomaly Detection & Alert Severity

Jalytics models transaction risk using **statistical baselines**:

- **Per-category baselines**:
  - A background job recomputes rolling baselines for each category over a configurable window (`BASELINE_WINDOW_DAYS`).
  - Baselines store the average amount and standard deviation for recent transactions.

- **Z-score based detection**:
  - Each new transaction is evaluated against its category baseline with a z-score.
  - Thresholds map to **severity levels**:
    - `medium` – moderately unusual amounts
    - `high` – significantly above normal (e.g. > 3σ)
    - `critical` – extreme outliers (e.g. > 5σ)

- **Low-variance fallback**:
  - When a category is very stable (stddev ≈ 0), Jalytics falls back to a rule like:
    - flag if `amount ≥ 2x` the stable baseline and deviation is materially large.

Alerts are persisted with **severity** so downstream systems and dashboards can prioritize the riskiest events first.

---

## 🗄 Database

* PostgreSQL
* Indexed on:

  * transaction timestamps
  * alert creation time
  * foreign keys

All schema changes are handled via migrations.

---

## 🏃 Running Locally

### Prerequisites

* Go 1.24+
* PostgreSQL
* Prometheus & Grafana
* Docker (optional)

### Environment & Configuration

Configuration is driven by environment variables:

```bash
DATABASE_URL=postgres://user:pass@localhost:5432/jalytics?sslmode=disable
PORT=8000
LOG_LEVEL=info               # optional (debug, info, warn, error)
BASELINE_WINDOW_DAYS=30      # optional, rolling window for baselines
ANOMALY_STRICT_MODE=false    # optional, enables more sensitive anomaly tiers
```

Values can be provided directly in your shell or via a `.env` file loaded by the server.

### Start Server

```bash
go run cmd/api/main.go
```

### Testing

#### Compile
```bash
go build ./...
```

#### Run static tests
```bash
go test ./...
```

#### Static analysis
```bash
go vet ./...
```
Prometheus scraping and Grafana dashboards have been **manually validated** by:
- sending sample transactions
- verifying anomaly alerts and their severities are emitted
- confirming metrics series update in real time.

---

## 🔒 Rate Limiting

The API uses a token-bucket rate limiter:

* 5 requests/second
* burst up to 10

This prevents abuse and simulates real production constraints.

---

## Grafana Dashboard

Custom dashboard UID: **jalytics_dashboard**

Panels:
- Transaction per category
- Anomalies over time
- HTTP request latency
- Database query latency
- (Optional) Transaction per minute & top merchants

Add custom thresholds for color coding to visualize critical metrics.

---

## 🧩 Design & Trade-offs

Jalytics is intentionally small but **designed like a production backend**:

- **Stack choices**
  - `chi` router for lightweight, composable HTTP middleware.
  - `pgx`/`pgxpool` for efficient PostgreSQL access and connection pooling.
  - Prometheus client for **pull-based** metrics and Grafana dashboards.

- **Anomaly detection strategy**
  - Uses a **simple, explainable statistical model** (rolling baselines + z-score) instead of a black-box ML model.
  - Pros: easy to reason about, debuggable, and tunable via config.
  - Trade-off: not as powerful as fully featured ML, but ideal for a focused backend demo.

- **Where logic lives**
  - Baselines are computed via SQL aggregations in PostgreSQL for **data locality and performance**.
  - The Go service layer owns the anomaly evaluation and alert creation, keeping business rules in code.

- **Safety & simplicity**
  - Global rate limiting to protect the service under load; could be extended to per-tenant/IP in a multi-tenant setting.
  - Config is env-driven with validation so misconfiguration is caught on startup rather than at runtime.

---

## 🚧 Future Enhancements

- Per-merchant anomaly baselines (not just category-wide)
- Webhook push notifications to external systems for critical alerts 
- Cursor-based pagination to handle high volume queries
- Transaction reversal API with idempotency keys for refund handling

---

## 📄 License

MIT License