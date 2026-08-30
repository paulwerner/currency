# AGENTS.md

Guidance for AI agents and human contributors working on this repository.

## What this is

`github.com/paulwerner/currency` is a dependency-light Go library for
handling monetary values: overflow-checked integer arithmetic in currency
minor units, ISO 4217 currency metadata, CLDR-derived rounding rules
(including cash rounding increments such as 0.05 CHF), and locale-based
formatting (`Amount.Display`). Kind-based displaying is planned; see
`docs/ROADMAP.md` for open items and design notes before starting work
on them.

## Layout

| Path | Purpose |
| --- | --- |
| `amount.go` | `Amount`: arithmetic, comparisons, split/alloc, rounding, JSON |
| `calc.go` | Overflow-checked integer primitives (`add`, `mul`, `round`, ...) |
| `currency.go` | `Currency`, `Kind` (Standard/Cash/Accounting), ISO parsing |
| `display.go` | Locale lookup and locale-based formatting (`Amount.Display`) |
| `common.go` | **Generated** from `internal/cldrgen/gen_common.go` — do not edit |
| `tables.go` | **Generated** CLDR currency and locale tables — do not edit |
| `doc.go` | Package docs and the `go:generate` directive |
| `internal/data/` | Compact string-table container used by the generated tables |
| `internal/gen/` | Go-code writer used by the generator |
| `internal/cldrgen/` | The CLDR table generator (`package main`) |

## Build, test, verify

```sh
make build      # go build ./...
make test       # go test -v -race ./...
make vet        # go vet ./...
make fmt-check  # fails if any file is not gofmt-formatted
make cover      # coverage summary
```

CI (`.github/workflows/ci.yml`) runs fmt-check, vet, build, race tests, and
a `go mod tidy` cleanliness check. Run all of these locally before pushing.

Green CI is the merge standard for this repository: every pull request must
pass the `test` check on its current head before it is merged, and a red
check is always worked immediately — fixed, or root-caused and explained on
the PR — never waited out. Repository admins should mirror this in GitHub
branch protection for `main` (Settings → Branches → require the `test`
status check).

## Code generation

`tables.go` and `common.go` are generated from Unicode CLDR data:

```sh
make gen-fetch   # downloads core.zip (CLDR_VERSION=40) and runs go generate
make gen         # runs go generate against an existing ./core.zip
```

- The generator is `internal/cldrgen`, executed from the repository root
  (paths are relative to the module root; it reads `./core.zip`).
- `core.zip` is gitignored; the generated `.go` files are committed.
- If `unicode.org` is unreachable (some sandboxes block it), fetch the
  needed files from `https://raw.githubusercontent.com/unicode-org/cldr/release-<N>/common/...`
  and zip them locally with the same `common/main/*.xml` and
  `common/supplemental/supplementalData.xml` layout. A partial zip needs
  `supplementalData.xml` plus every locale in the allowlist in
  `internal/cldrgen/gen_locales.go` and the locales on their parent
  chains (`root`, `en_001`, `es_419`, `no`); the generator fails with the
  name of any file that is missing.
- After regenerating, `git diff tables.go common.go` should be empty unless
  you changed the CLDR version or the generator on purpose.

## Conventions

- Never edit files whose header says `// Code generated. DO NOT EDIT.`
  Change `internal/cldrgen` and regenerate instead.
- Monetary values are `int` minor units. Every arithmetic primitive in
  `calc.go` returns `(result, ok)`; public `Amount` methods translate
  `!ok` into `ErrInvalidOperation`. Keep new operations overflow-checked
  and pure (amounts are immutable — return new `*Amount` values).
- Exported errors are sentinel values prefixed with `currency: ` and are
  compared with `errors.Is` in tests.
- Table-driven tests in `_test.go` files next to the code; use the
  existing style (`[%v]: want ..., got ...`).
- The module targets the Go version in `go.mod`; don't introduce
  dependencies beyond `golang.org/x/text` (generator-only) without a very
  good reason — the runtime library is dependency-free by design.

## Roadmap discipline

Open feature work (locale data generation, locale-based formatting,
kind-based display) is specified in `docs/ROADMAP.md`, including design
notes agreed during review. Read it first; update it (status table and,
when relevant, the change log) whenever you complete or re-scope an item.
