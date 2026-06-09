# neneloga

A logging and observability pipeline service — similar in spirit to Prometheus or Datadog — built to explore log ingestion, processing, storage, and retrieval at scale. Implemented in both **Go** and **Rust** to demonstrate cross-language systems programming.

---

## What it does

neneloga is a self-hosted log pipeline. Applications and agents push log events to it, the pipeline normalizes and enriches them, stores them efficiently, and exposes APIs for search, tailing, and analytics.

```
Producers
  → Ingestion API
  → Parse / Validate / Enrich
  → Buffer / Batch
  → Storage
  → Index
  → Query / Tail / Export API
```

---

## Core modules

| Module | Responsibility |
|---|---|
| Ingestor | Accepts logs over HTTP, gRPC, TCP, or from agents |
| Parser | Converts raw text/JSON/syslog into structured records |
| Validator | Enforces schema, size limits, encoding |
| Enricher | Attaches host, env, region, trace IDs |
| Processor | Redaction, deduplication, sampling, routing |
| Buffer | Batches events before writing to storage |
| Storage writer | Chunks, compresses, and persists log data |
| Indexer | Builds lookup structures for fast retrieval |
| Query engine | Serves filtered, paginated, and streaming queries |

---

## API (planned endpoints)

| Method | Path | Description |
|---|---|---|
| `POST` | `/ingest` | Push a single log event |
| `POST` | `/batch` | Push a batch of log events |
| `GET` | `/logs` | Query logs with filters |
| `GET` | `/tail` | Live stream of incoming logs |
| `GET` | `/search` | Full-text search across log messages |
| `GET` | `/trace/:id/logs` | Fetch all logs for a trace ID |
| `GET` | `/stats` | Log volume and pipeline health |
| `GET` | `/export` | Download logs as JSON or CSV |
| `GET` | `/health` | Service health check |

---

## Implementations

Two server implementations exist side by side to demonstrate the same system design in different languages and ecosystems.

### Go (`/go`)

Built with [Gin](https://github.com/gin-gonic/gin).

```bash
cd go
go run main.go
```

### Rust (`/rust`)

Built with [Actix-web](https://actix.rs/).

```bash
cd rust
cargo run
```

Both run on port `8080`.

---

## Storage model

neneloga uses a tiered storage approach:

- **Hot** — recent logs in a fast queryable store (e.g. ClickHouse or SQLite for local dev)
- **Warm** — compressed chunk files on disk, indexed by time and service
- **Cold** — object storage archive (S3-compatible) for long-term retention

For local development, logs are stored as compressed JSON Lines chunks with a SQLite metadata index.

---

## Log record schema

Every log event is normalized to this structure:

```json
{
  "timestamp": "2026-06-08T12:00:00Z",
  "ingested_at": "2026-06-08T12:00:01Z",
  "service": "auth-api",
  "level": "ERROR",
  "message": "login failed for user",
  "tenant_id": "tenant_123",
  "trace_id": "abc123",
  "host": "node-7",
  "env": "production",
  "raw": "<original log line>"
}
```

---

## Requirements

- Go 1.25+
- Rust + Cargo (2024 edition)

---

## Project structure

```
neneloga/
├── go/           # Go implementation
│   ├── server/   # Route handlers
│   └── main.go
├── rust/         # Rust implementation
│   └── src/
│       └── main.rs
└── docs/         # Architecture notes and design decisions
```

---

## Contributing

See [CONTRIBUTING.md](./CONTRIBUTING.md).
