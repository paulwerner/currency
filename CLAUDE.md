# CLAUDE.md

See [AGENTS.md](AGENTS.md) for repository guidance: project layout,
build/test/generate commands, code conventions, and roadmap discipline.

Quick reminders:

- `tables.go` and `common.go` are generated — edit `internal/cldrgen` and
  run `make gen` instead of touching them.
- Before pushing: `make build test vet fmt-check`.
- Open feature work is tracked in [docs/ROADMAP.md](docs/ROADMAP.md);
  update it in the same PR that completes or re-scopes an item.
