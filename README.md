# example-banking

**Reference application**, not a library: an event-sourced **banking account** aggregate, HTTP API, and wiring that shows how **goeventsource** + **pgx** fit together in a real (small) service. Copy patterns from here when you build your own domain; do not import this module as a dependency in production unless you intend to fork it.

## What to copy


| Piece              | Location                            | Takeaway                                              |
| ------------------ | ----------------------------------- | ----------------------------------------------------- |
| Aggregate + events | Package `banking` (module root) | `BaseRoot`, `goeventsource.Record`, domain invariants        |
| Root JSON codec    | [codec.go](codec.go)              | `RootEncodeDecoder` for snapshots / API               |
| Application API    | [internal/](internal/)            | Service layer over `banking.Repo`                     |
| HTTP + wiring      | [cmd/](cmd/)                      | `pgxtest` pool, `pgx` repository options, listen |


## Install (to run or study)

```bash
git clone <your-fork-or-mirror>
cd example-banking
go mod download
```

This module depends on:

- `github.com/goeventsource/goeventsource` — core
- `github.com/goeventsource/pgx` — PostgreSQL + [pgxtest](../pgx/README.md)
- `github.com/google/uuid`, `github.com/Rhymond/go-money` — domain helpers

## Run the HTTP server

[cmd/main.go](cmd/main.go) uses `pgxtest.NewDemoPool` (Docker-backed). With Docker running:

```bash
go run ./cmd
```

The process listens on `:80` (adjust in `main` if needed). For production you would replace `pgxtest` with your own pool and migrations.

## Wiring recap (matches `cmd/main.go`)

1. **Factory** — `func(uuid.UUID, goeventsource.Version) *banking.Account` rebuilds the root from the event stream.
2. **Pool** — `pgxtest.NewDemoPool(ctx)` for demos; use `pgxpool.New` + your DSN in production.
3. **Repository config** — `pgxtest.NewRepositoryConfig(db, factory)` then set `cfg.StoreConfig.Codec` to a `DomainEventEncodeDecoder` map for every account event type.
4. **Snapshotter** — build with `pgxtest.NewSnapshotter` + `goeventsource.SnapshotterWriteStrategyByVersionStep` (example uses every 3 versions).
5. **Repository** — append `pgx.WithSnapshotterOpt(snap)` to `cfg.Opts`, then `pgxtest.NewRepository(cfg)`.
6. **Service** — `internal.NewService(repo, exchanger)`; the sample passes `nil` exchanger (no currency integration in the demo).
7. **HTTP** — `cmd/internal.New(svc, codec)` exposes handlers using the same root codec.

For production-shaped SQL wiring without test helpers, follow [../pgx/README.md](../pgx/README.md) and call `pgx.NewStore` / `pgx.NewRepository` directly.

## Layout

```
.
├── account.go, events, codec …   # domain (package banking)
├── internal/                       # use cases
└── cmd/
    ├── main.go                     # composition root
    └── internal/                   # HTTP
```

## Tests

```bash
go test ./...
```

## Local `replace` directives

[go.mod](go.mod) may point `goeventsource` and `pgx` at sibling paths for development—remove after modules are public; see [../README.md](../README.md).