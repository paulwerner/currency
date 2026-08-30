---
name: steward
description: Repo-specific conventions for watching and driving pull requests in this repository to green.
---

# PR stewardship for paulwerner/currency

Green CI is the merge standard here: a PR is not done until the `test`
check (`.github/workflows/ci.yml`) passes on its current head and there is
no merge conflict. Never wait out a red check — fix it, or root-cause it
and explain on the PR.

## Validating before every push

Reproduce what CI runs, locally, in this order:

```sh
make build test vet fmt-check
go mod tidy && git diff --exit-code go.mod go.sum
```

A CI failure in the `Check formatting` step means `gofmt -w .`; in
`Verify go.mod is tidy` it means committing the result of `go mod tidy`.

## Repo-specific rules

- `tables.go` and `common.go` are generated (`// Code generated. DO NOT
  EDIT.`). Never patch them directly — change `internal/cldrgen` and run
  `make gen` (needs `./core.zip`; `make gen-fetch` downloads it). If
  `unicode.org` is blocked, fetch the CLDR files from
  `https://raw.githubusercontent.com/unicode-org/cldr/release-<N>/common/...`
  and zip them locally with the same `common/...` layout.
- After regenerating, `git diff tables.go common.go` must be empty unless
  the CLDR version or generator changed on purpose.
- Keep fixes minimal and overflow-checked; amounts are immutable and the
  runtime library stays dependency-free (x/text is generator-only).
- Tests are table-driven and sit next to the code; assert sentinel errors
  with `errors.Is`.
- The CI test job runs with `-race`; a failure only under CI is likely a
  race — run `go test -race ./...` locally, not plain `go test`.
- Update `docs/ROADMAP.md` (status table, change log) in the same PR that
  completes or re-scopes a roadmap item.
