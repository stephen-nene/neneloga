# Contributing to neneloga

Thanks for your interest. Here's how to get started.

---

## Prerequisites

- Go 1.25+
- Rust + Cargo (2024 edition)
- Git

---

## Local setup

```bash
git clone https://github.com/<you>/neneloga.git
cd neneloga
```

**Go server:**

```bash
cd go
go mod download
go run main.go
```

**Rust server:**

```bash
cd rust
cargo build
cargo run
```

---

## Project structure

```
neneloga/
├── go/           # Go implementation (Gin)
├── rust/         # Rust implementation (Actix-web)
└── docs/         # Architecture and design notes
```

Each implementation lives independently. Changes to one do not require changes to the other unless you're adding a new shared endpoint.

---

## Making changes

1. Fork the repo and create a branch from `main`
2. Keep commits focused — one logical change per commit
3. Test your changes locally before opening a PR
4. Open a pull request with a short description of what changed and why

---

## Adding a new endpoint

When adding an endpoint:

- Add it to both the Go and Rust implementations if it's a core pipeline endpoint
- Update the endpoint table in `README.md`
- Keep handlers thin — business logic belongs in internal packages, not route handlers

---

## Code style

**Go** — follow standard `gofmt` formatting. Run `go vet ./...` before committing.

**Rust** — follow `rustfmt` defaults. Run `cargo clippy` before committing.

---

## Reporting issues

Open a GitHub issue with:

- What you expected to happen
- What actually happened
- Steps to reproduce
- Go or Rust version (`go version` / `cargo --version`)
