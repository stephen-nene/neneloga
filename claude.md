# The Logging System — A Complete Mental Map

Think of a logging system as a **pipeline** with four distinct stages. Every major logging company (Datadog, Grafana/Loki, Elastic, Splunk, Logtail) is just an opinionated implementation of this same pipeline. Once you see the stages clearly, you can design your own version.

---

## Stage 1 — Log Ingestion (How logs arrive at your logger)

This is the entry point. Logs come in two fundamental models: **push** (the service sends them to you) and **pull** (you go fetch them from the service).

---

### Push-based Ingestion

The service actively sends logs to your logger.

**HTTP / REST API**
The service makes a POST request to your logger's endpoint — e.g., `POST /ingest` with a JSON body. This is the simplest to implement and the most universal. Almost every service can make an HTTP call. The downside is that if your logger is down, logs are lost unless the service has a retry mechanism.

**gRPC**
Instead of HTTP/JSON, the service uses Protocol Buffers over HTTP/2. Much faster and more efficient for high-volume logs. Used internally by companies like Google. Requires a proto schema definition. Overkill for small systems but worth knowing it exists.

**UDP (Syslog)**
Fire-and-forget. The service sends a UDP packet and doesn't wait for any acknowledgment. Extremely fast, zero blocking of the main application. The tradeoff: logs can be silently lost if the network drops or your logger is busy. The `syslog` protocol (RFC 5424) is the standard format here. Still used heavily in Linux system logs.

**TCP (Syslog)**
Same syslog format but over TCP, so you get reliable delivery. The connection is maintained and the service knows if a log failed to deliver. Slower than UDP but you don't lose data. Most network devices (routers, firewalls) emit logs this way.

**Message Queue / Event Broker**
The service publishes logs to a broker — **Kafka**, **RabbitMQ**, **Redis Pub/Sub**, **AWS SQS**, **NATS**. Your logger is a consumer of that queue. This decouples the service from the logger entirely. If your logger goes down, logs sit in the queue waiting. If your service spikes, the queue absorbs the burst. This is how large-scale production systems do it. Kafka specifically is used by LinkedIn, Netflix, and Uber for exactly this kind of pipeline.

**WebSockets**
A persistent bidirectional connection between the service and the logger. The service streams logs in real-time over an open socket. Useful when you want low-latency delivery without the overhead of a new HTTP connection per log.

**Logging SDK / Library (In-process)**
Your service imports a library (like Winston, Pino, Bunyan in Node.js, or Loguru in Python) that you configure to ship logs directly to your logger's endpoint. The library handles batching, retries, and formatting internally. This is elegant because the service code just calls `logger.info("...")` and the SDK handles the rest. Most SaaS loggers (Logtail, Axiom) give you an SDK for this.

---

### Pull-based Ingestion

Your logger goes to the service and fetches logs.

**File Tailing**
Your logger watches a `.log` file on disk and reads new lines as they are appended — like `tail -f` but programmatic. This is how **Filebeat** (from Elastic) and **Promtail** (from Grafana/Loki) work. The service just writes to disk normally, and an agent sitting on the same machine picks up changes. Zero code change required in the service.

**Container Log Scraping**
Docker and Kubernetes expose container stdout/stderr as log streams. A collector running in the cluster (like Fluentd, Fluent Bit, or Promtail as a DaemonSet) reads these streams and ships them. Your service just prints to stdout — infrastructure handles everything else. This is the Kubernetes-native logging pattern.

**Polling an API endpoint**
Your logger periodically calls `GET /logs?since=<timestamp>` on the service's own API. The service exposes its own log endpoint, and your logger pulls from it on a schedule. Simple but introduces latency (you only get logs as often as you poll) and puts load on the service.

**Database Querying**
If the service writes logs directly to a database table, your logger can query that table. This is more of an internal pattern — not typical for distributed systems, but valid in monolithic setups.

---

### Hybrid — The Agent Pattern

An **agent** (Filebeat, Fluentd, Fluent Bit, Vector, Logstash) is a lightweight process that runs on the same host as your service. It tails local log files (pull from disk), then ships them to your central logger (push to network). This is the dominant pattern in production because:
- The service doesn't need to know about the logger
- Agents handle buffering and retries
- One agent can ship logs from multiple services on the same host

---

## Stage 2 — Log Processing (What happens before saving)

Raw logs are messy, inconsistent, and often contain information you don't want or are missing information you need. Processing is the transformation layer.

---

### Parsing

Converting raw unstructured text into structured data you can query.

**JSON parsing** — If logs are already JSON (`{"level":"error","msg":"...","ts":1234}`), you parse them directly. This is the ideal case.

**Grok patterns** — A named regex system used by Logstash. You define a pattern like `%{TIMESTAMP_ISO8601:timestamp} %{LOGLEVEL:level} %{GREEDYDATA:message}` and it extracts fields from a raw string. Has hundreds of built-in patterns for common formats (Apache access logs, nginx, syslog, etc.).

**Regex extraction** — Raw regex. More flexible than Grok, less readable. You define capture groups and map them to field names.

**Key-value parsing** — Logs in format `level=error service=api latency=120ms`. Split on spaces, then on `=`. Logfmt is a standardized version of this.

**Multiline handling** — Stack traces span multiple lines. You need to detect "this is a continuation of the previous log entry" (usually by checking if a line starts with a timestamp or not). If you don't handle this, each line of a stack trace becomes a separate log entry.

---

### Enrichment

Adding context that the service didn't include but is useful for debugging.

**Metadata injection** — Automatically adding `hostname`, `environment` (prod/staging/dev), `service name`, `region`, `deployment version`. The service shouldn't have to repeat this on every log line.

**Timestamp normalization** — Different services format timestamps differently. Normalize everything to ISO 8601 UTC. If a log has no timestamp, stamp it with the ingestion time.

**GeoIP lookup** — Given an IP address in the log, enrich it with country, city, ASN. Useful for access logs. Maxmind GeoLite2 is the standard free database for this.

**User-agent parsing** — `Mozilla/5.0 (Windows NT...)` → `{browser: "Chrome", os: "Windows", version: "120"}`. Libraries like `ua-parser` do this.

**Trace/Span ID correlation** — If your service uses distributed tracing (OpenTelemetry), you can attach the `trace_id` and `span_id` to every log. Now you can jump from a log entry to the full request trace.

**Lookup table enrichment** — Given a `user_id`, look up and attach `plan: "enterprise"`, `company: "Acme Corp"`. Requires a side data source.

---

### Filtering

Deciding what to keep and what to drop.

**Log level filtering** — Drop `DEBUG` and `TRACE` logs in production. Only keep `INFO`, `WARN`, `ERROR`. This alone can reduce log volume by 80%.

**Noise suppression** — Some logs repeat constantly and provide no value (health check endpoints, heartbeat messages). Drop them by matching patterns.

**Sampling** — Keep only 1 in 100 `INFO` logs from a very noisy service. Keep 100% of `ERROR` logs. Tail-based sampling (making the keep/drop decision after seeing the full trace) is more sophisticated.

**Deduplication** — If the same error fires 10,000 times in a minute, store it once with a count of 10,000 rather than 10,000 identical rows.

**Rate limiting** — Cap ingestion from a single noisy source so it doesn't flood your pipeline and cause backpressure.

---

### Transformation

Reshaping the data.

**Field renaming** — Normalize inconsistent field names across services. One service calls it `msg`, another `message`, another `text` — normalize to `message`.

**PII redaction / masking** — Automatically detect and mask credit card numbers, emails, passwords, national IDs before they ever hit storage. You match patterns like `\b\d{4}[- ]?\d{4}[- ]?\d{4}[- ]?\d{4}\b` and replace with `****`. This is legally required under GDPR and PCI-DSS.

**Log level normalization** — One service says `CRITICAL`, another says `FATAL`, another says `5`. Normalize everything to a standard scale: `DEBUG < INFO < WARN < ERROR < FATAL`.

**Format conversion** — Convert incoming syslog format into your internal JSON schema. Convert CSV logs into JSON.

---

### Routing

Sending different logs to different destinations based on content.

**By severity** — `ERROR` and `FATAL` logs go to an alerting system (PagerDuty, Slack). All logs go to cold storage.

**By type** — Audit logs (who did what) go to a tamper-proof compliance store. Security logs go to a SIEM (Security Information and Event Management) system. Application logs go to Elasticsearch.

**By service** — Logs from the payments service are stored separately with stricter retention and access controls.

**Fan-out** — One log entry gets sent to multiple destinations simultaneously (e.g., Elasticsearch for search AND S3 for long-term archival).

---

### Aggregation & Batching

**Batching** — Instead of writing every single log to disk/DB one at a time (extremely slow), buffer them in memory and write in batches of 500 or every 5 seconds, whichever comes first.

**Windowed aggregation** — Count how many 4xx errors occurred in each 1-minute window. Store the count, not every individual 4xx. Reduces storage while preserving the signal.

**Metric extraction** — Parse `latency=120ms` out of a log and feed it into a time-series metrics store. Logs become the source of truth for metrics.

---

## Stage 3 — Log Storage (Where logs live)

Each storage type has different tradeoffs on query speed, cost, scalability, and how long you can retain data.

---

### File System (.log files)

The most basic. Append-only text files on disk. Log rotation (e.g., with `logrotate` on Linux) creates a new file daily and compresses old ones. Simple, cheap, no dependencies. Bad for querying across many files, across machines, or for large volumes.

**Structured variants** — Instead of plain text, write JSON Lines (one JSON object per line, `.jsonl`). Still a flat file but machine-parseable.

**Log rotation strategies:**
- By time: new file every hour/day
- By size: new file every 100MB
- Compressed archives: old files gzipped automatically
- Retention policy: delete files older than 30 days

---

### Relational Databases (SQL)

PostgreSQL, MySQL, SQLite. You define a `logs` table with columns for `id`, `timestamp`, `level`, `service`, `message`, `metadata` (JSONB in Postgres). You get SQL queries, joins, aggregations. Postgres JSONB in particular lets you index and query inside the JSON metadata field.

Good for small-to-medium log volumes. Struggles with billions of rows and write-heavy workloads without careful indexing and partitioning.

**TimescaleDB** is a Postgres extension that adds time-series optimizations — automatic chunking of data by time, fast range queries, compression of old chunks. Purpose-built for this kind of append-heavy, time-range-query workload.

---

### Document / Search Stores

**Elasticsearch / OpenSearch** — The most widely used log storage in the industry. Logs are stored as JSON documents. You get full-text search, field-level filtering, aggregations, and near-real-time indexing. Kibana (or OpenSearch Dashboards) sits on top as the UI. The entire ELK stack (Elasticsearch, Logstash, Kibana) is a complete logging pipeline. Very powerful but resource-hungry (RAM, disk).

**MongoDB** — Flexible document storage. You can query by any field, store nested objects. Not optimized specifically for logs but works, especially if you're already running MongoDB.

---

### Columnar / Analytics Databases

Designed for reading many rows across a few columns very fast — exactly the pattern of log analytics.

**ClickHouse** — An open-source columnar database. Ingests millions of rows per second. Extremely fast aggregation queries ("count errors per service per hour for the last 7 days"). Used internally by Cloudflare, Uber, and several logging SaaS products. This is what serious logging systems use for the hot storage layer.

**Apache Parquet on S3** — Store logs as Parquet files (a columnar format) in object storage, then query them with Athena (AWS), BigQuery (GCP), or DuckDB locally. Very cheap. Queries are slower but cost almost nothing per month.

---

### Time-Series Databases

**InfluxDB** — Optimized for timestamped numeric data. Less suited to raw log messages, better suited to metrics derived from logs (error rate, latency percentiles). Think of it as where you send the numeric signals extracted from your logs.

**Prometheus** — Specifically a metrics system, not a log system. But worth knowing: logs and metrics are different things. Prometheus scrapes metrics endpoints; it doesn't store log text.

---

### Purpose-built Log Databases

**Grafana Loki** — Stores logs indexed only by labels (like `service=api`, `env=prod`), not the full text. This makes it very cheap to run. Full-text search happens at query time (slower). Pairs with Promtail for ingestion and Grafana for visualization. The philosophy: "logs are like metrics, just with a text payload."

**Quickwit** — A newer open-source log search engine built on top of object storage (S3). Indexes are stored on S3 cheaply, search nodes are stateless. Designed to be cloud-native and cost-effective.

**VictoriaLogs** — From the VictoriaMetrics team. Extremely low resource usage. Accepts the same ingestion format as Elasticsearch. Good for self-hosted setups where you want Elasticsearch-like querying without the Elasticsearch RAM bill.

---

### Object Storage (Cold / Archive Tier)

**AWS S3, Google Cloud Storage, Azure Blob, Backblaze B2, Cloudflare R2** — Dirt cheap. $0.02/GB/month. You compress and ship old logs here. You don't query directly — you either download and process them, or use a query layer on top (Athena, Trino, DuckDB). This is the archival tier. "We keep 90 days hot in ClickHouse, and 7 years cold in S3."

---

### Cache / Buffer Layer

**Redis** — Not a permanent store. Used as a high-speed in-memory buffer between ingestion and the real storage. Your logger writes incoming logs to Redis first (fast), then a background worker drains Redis into the real DB in batches (durable). Absorbs traffic spikes without losing logs.

---

### Hot / Warm / Cold Tiering

Large systems don't use one storage for everything. They tier by age:

- **Hot tier** (last 7 days) — Fast storage (ClickHouse, Elasticsearch). Full query capability. Expensive.
- **Warm tier** (last 90 days) — Slower but cheaper (compressed ClickHouse, smaller Elasticsearch). Still queryable.
- **Cold tier** (1 year+) — Object storage (S3 + Parquet). Very cheap. Queried rarely and slowly.

Data is automatically migrated between tiers by TTL policies.

---

## Stage 4 — Log Retrieval (How you get logs out)

---

### REST API

`GET /logs?level=error&service=api&from=2026-06-01&to=2026-06-08&limit=100`

The most universal interface. You expose a query endpoint. Supports filtering by any indexed field. Returns paginated JSON results. Cursor-based pagination (using the last log's ID as a cursor) is better than offset-based for large datasets because offsets become slow as pages get deep.

---

### Full-Text Search

If your storage supports it (Elasticsearch, ClickHouse, Loki), you can search for any substring in the log message. `GET /logs?q="connection refused"`. This is how you find unknown errors without knowing the exact field values.

**Query languages:**
- **Lucene / KQL (Kibana Query Language)** — `level:error AND service:api AND message:"timeout"`
- **LogQL (Loki)** — `{service="api"} |= "error" | json | level = "error"`
- **SQL** — `SELECT * FROM logs WHERE level = 'error' AND message LIKE '%timeout%'`
- **ClickHouse SQL** — Standard SQL plus time-series functions

---

### Real-Time Tailing

"Show me logs as they come in, live."

**Server-Sent Events (SSE)** — The server pushes new log entries to the client over a long-lived HTTP connection. One-directional, simple to implement, works over standard HTTP. The client opens `GET /logs/stream` and keeps the connection open.

**WebSockets** — Bidirectional. The client can send filter commands over the same connection while receiving log entries. More complex to implement but more interactive.

**Long Polling** — The client asks "any new logs since X?", the server holds the request open until something arrives, then responds. Old-school but works everywhere.

---

### Aggregations and Analytics Queries

Beyond finding individual log lines — summarizing them.

- Count of errors per hour over the last 24 hours
- Top 10 error messages by frequency
- Average response latency per endpoint
- P95 / P99 latency percentiles
- Number of unique users who hit errors today

These are typically SQL aggregation queries or pre-computed using materialized views / streaming aggregation at ingestion time.

---

### Alerting Queries (Continuous Retrieval)

Your logger runs queries on a schedule: "every 1 minute, count errors in the last 5 minutes. If > 100, fire an alert." This is not a user-facing retrieval — it's an automated background process. Outputs go to Slack, email, PagerDuty, or a webhook.

---

### Export / Download

Users want to download logs for offline analysis.

- **CSV export** — good for spreadsheet analysis
- **JSON Lines export** — good for programmatic processing
- **Raw .log file download** — good for passing to tools that expect log files
- **Streaming download** for large result sets (don't load 10 million rows into memory — stream them to disk)

---

### Dashboards / Visualization

Retrieval through a UI rather than an API.

- **Grafana** — connects to Loki, ClickHouse, Elasticsearch, Prometheus, TimescaleDB — renders time-series charts, log tables, histograms
- **Kibana** — Elasticsearch-specific. Discover tab for log search. Dashboard tab for charts.
- **Custom UI** — you build your own frontend that talks to your REST API

---

## The Full Mental Model

```
[Service A] ─┐
[Service B] ─┤── INGESTION ──► PROCESSING ──► STORAGE ──► RETRIEVAL
[Service C] ─┘    (push/pull)   (parse/enrich  (hot/warm/   (API/search/
                                 filter/route)   cold tiers)  stream/dashboards)
```

Every logging product — Datadog, Splunk, Grafana Cloud, Logtail, Axiom, New Relic — is a specific set of choices at each of these four stages. Now you know the entire decision space they're choosing from.
