# Roadmap

This document tracks the remaining development roadmap for the `currency`
library, together with a build and change log of the modernization pass that
brought the repository to its current state.

## Status overview

| Item | Status |
| --- | --- |
| Calculator (overflow-checked arithmetic) | ✅ done |
| CLDR data fetching | ✅ done |
| Code generator | ✅ done |
| Currency data generation | ✅ done |
| De-/serialization (JSON) | ✅ done |
| Documentation (package docs, AGENTS.md, README) | ✅ done |
| Locale data generation | ⏳ open |
| Locale-based formatting | ⏳ open |
| Kind-based displaying (standard, cash, accounting) | ⏳ partially prepared |

## Open items

### 1. Locale data generation

Extend the generator (`internal/cldrgen`) to emit locale tables alongside the
existing currency table: per-locale number symbols (decimal separator, group
separator, minus sign), the CLDR `standard` and `accounting` currency format
patterns, and per-locale currency symbols (e.g. `$`, `€`, `US$`, `CHF`).

Implementation notes from the groundwork investigation:

- The data lives in the CLDR `common/main/<locale>.xml` files under the
  `numbers` section (`symbols`, `currencyFormats`, `currencies`), which the
  already-used `golang.org/x/text/unicode/cldr` decoder parses when the
  section filter includes `numbers` (it already does).
- The generator should carry an explicit locale allowlist so output is
  deterministic regardless of how many locale files the `core.zip` contains.
- Restrict extraction to the `latn` numbering system for a first iteration.
- CLDR data is inherited along a locale's parent chain (`de-CH` → `de` →
  `root`, with exceptions such as `zh-Hant` → `root` listed in
  `supplementalData.xml` under `parentLocales`). Resolve at generation time
  and store per-locale diffs against the parent to keep tables small; store
  the parent per locale so the runtime walks the same chain.
- Parse the format patterns at generation time into positive/negative
  prefix/suffix affixes plus grouping sizes (`#,##0.00` → 3/3,
  `#,##,##0.00` → 3/2), leaving only the `¤` symbol placeholder for runtime
  substitution. Select `currencyFormatLength` elements without a `type`
  attribute (the `short` variants are compact notation and out of scope).
- Note for sandboxed environments: `unicode.org` may not be reachable. The
  same files can be fetched per-file from the GitHub mirror, e.g.
  `https://raw.githubusercontent.com/unicode-org/cldr/release-40/common/main/de.xml`,
  and zipped locally as `common/...` — the generator only needs
  `common/supplemental/supplementalData.xml` plus the allowlisted
  `common/main/*.xml` files.

### 2. Locale-based formatting

Add a runtime formatter on top of the generated locale tables:

- `Amount.Display(locale string) string` — format using the locale's
  standard pattern, e.g. `19,99 €` for `de` and `€19.99` for `en`.
- Locale lookup: normalize the tag (`de_CH`/`de-CH`), exact match first,
  then truncate subtags, then fall back to `root`.
- Symbol lookup walks the stored parent chain and falls back to the ISO
  code when no symbol is defined.
- Apply the CLDR currency-spacing rule in simplified form: insert a
  non-breaking space when a letter-final symbol abuts a digit (`CHF 12.34`).
- Out of scope initially: non-`latn` digit systems, per-currency pattern
  overrides, bidi isolation for RTL locales; document these limitations.

### 3. Kind-based displaying

The groundwork is already merged: `Kind` now carries a `format` style, so
`Accounting` is a distinct value from `Standard`, and `Cash` rounding
(e.g. 0.05 CHF, 0.50 DKK) works via `Amount.Round(Cash)`.

Remaining work, once locale formatting exists:

- `Amount.DisplayKind(locale string, k Kind) string` — `Cash` applies cash
  rounding and the cash scale before rendering; `Accounting` uses the
  accounting pattern (negative amounts in parentheses where the locale says
  so, e.g. `(€19.99)` for `en`).

## Build and change log

All changes below were made on the `claude/currency-repo-review-roadmap-d5k701`
branch during the 2026-08 review pass. Build/verify with `make build test vet
fmt-check`; regenerate tables with `make gen-fetch` (fetches CLDR
`core.zip`, then runs `go generate`).

### Module and layout

- Go directive raised from 1.18 to 1.25.0 — the actual floor imposed by
  `golang.org/x/text` v0.41.0 (upgraded from v0.3.7); nothing in the
  library itself needs a newer toolchain, so consumers one Go release
  behind are not cut off.
- Library moved from `pkg/` to the module root, so
  `go get github.com/paulwerner/currency` imports the package directly, as
  the README quick start always claimed. The old import path
  `.../currency/pkg` is gone.
- The generator moved to `internal/cldrgen` (a regular `package main`
  command run via `go generate`), keeping its dependencies tracked in
  `go.mod`. Verified to reproduce the committed `tables.go` and `common.go`
  byte-for-byte from a CLDR 40 `core.zip`.
- Deprecated `io/ioutil` and `// +build` usages replaced; vet warning about
  a non-constant format string fixed; `internal/gen/writer.go` gofmt-ed.

### Bug fixes

- `round()` ignored the rounding increment entirely, making cash rounding
  impossible; rounding is now increment-aware and `Round(Standard)` is a
  no-op for amounts already in minor units (previously it rounded away the
  entire fractional part — e.g. USD 12.34 became 12.00).
- `CurrencyFromISO("XXX")` returned `(nil, nil)`, causing nil dereferences
  in any subsequent use; it now returns the canonical `XXX` value.
- `IsPositive()` returned `true` for zero; zero is now neither positive nor
  negative.
- `Alloc()` summed ratios without overflow checking.
- `Amount` stores its `Currency` by value; a nil currency can no longer be
  smuggled into an amount, and the zero `Amount` is a usable `XXX 0.00`.

### API completed (was `panic("not implemented")`)

- Comparisons: `Cmp`, `GreaterThan`, `GreaterThanOrEqual`, `LessThan`,
  `LessThanOrEqual` (currency-mismatch aware).
- `Abs`, `Neg` (overflow-checked, returning new amounts).
- `String()` → `"EUR 12.34"` (locale-independent, standard scale).
- JSON de-/serialization: `{"amount":1234,"currency":"EUR"}` with strict
  decoding (unknown fields and missing fields rejected).
- Renamed `ErrSplitNegative` → `ErrInvalidSplitNumber` (it also covers
  zero); all error messages share the `currency: ` prefix.

### Tests, CI, tooling

- New tests: currency-mismatch and overflow paths, split/alloc error and
  sum-preservation properties, kind-based rounding against real CLDR data
  (CHF/DKK/SEK/CAD), predicates, `Abs`/`Neg`, `String`, JSON round trips
  and failure modes, `XXX` normalization, malformed ISO codes,
  `internal/data` table lookups and `FixCase`.
- GitHub Actions CI: gofmt check, `go vet`, build, `go test -race`, and a
  `go mod tidy` cleanliness check.
- Makefile: `.PHONY` hygiene, `build`/`cover`/`vet`/`fmt`/`fmt-check`
  targets, parameterized `CLDR_VERSION`; `clean` no longer deletes the
  committed generated tables.
- `fetch-cldr.sh`: `set -eu`, required `CLDR_VERSION`, `curl -fL` with
  retries over https.

### Documentation

- Package documentation in `doc.go`; contributor/agent guide in
  `AGENTS.md`, referenced by `CLAUDE.md`; README rewritten to match the
  actual API; this roadmap document.

### Review follow-ups

Addressed from the first review round on the modernization PR:

- `pow()` no longer reports a false overflow caused by one redundant
  final squaring of the base (e.g. `pow(2, 62)` on 64-bit).
- `String()` no longer prints a double sign for `math.MinInt` in
  scale-0 currencies.
- `UnmarshalJSON` wraps out-of-range amounts (32-bit platforms) in
  `ErrInvalidJSON` like every other decode failure.
- `Currency` methods moved to value receivers and `Amount.Currency()`
  returns a value, closing the interior-pointer hole that allowed
  mutating an amount's currency in place.
- `NewAmount` no longer returns an always-nil error.
- Go directive lowered from 1.26 to the actual 1.25.0 floor (see above).
