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
| Locale data generation | ✅ done |
| Locale-based formatting | ✅ done |
| Kind-based displaying (standard, cash, accounting) | ⏳ partially prepared |

## Open items

### 1. Kind-based displaying

The groundwork is merged: `Kind` carries a `format` style, so `Accounting`
is a distinct value from `Standard`, `Cash` rounding (e.g. 0.05 CHF,
0.50 DKK) works via `Amount.Round(Cash)`, and the locale formatter's
internal `display` method (see the change log) already takes the
`formatStyle` selecting between the standard and accounting patterns.

Remaining work:

- `Amount.DisplayKind(locale string, k Kind) string` — `Cash` applies cash
  rounding and the cash scale before rendering; `Accounting` uses the
  accounting pattern (negative amounts in parentheses where the locale says
  so, e.g. `(€19.99)` for `en`).

## Build and change log

Build/verify with `make build test vet fmt-check`; regenerate tables with
`make gen-fetch` (fetches CLDR `core.zip`, then runs `go generate`).

### Locale-based formatting (`claude/roadmap-item-pr-locals-74ruxf`, 2026-08)

- `Amount.Display(locale string) string` in the new `display.go` renders
  an amount with the locale's standard CLDR currency pattern, separators,
  and symbols: `19,99 €` for `de`, `€19.99` for `en`, `EUR 1’234.56` for
  `de-CH`, `$12,34,567.89` (3/2 grouping) for `en-IN`.
- Locale lookup canonicalizes the tag (case folding, `_` vs `-`, script
  subtags title-cased), tries an exact match via binary search over the
  sorted table, truncates subtags (`zh-Hant-TW` → `zh-Hant`, `en-US` →
  `en`), and falls back to `root` — so unknown locales still render.
- Field resolution flattens the stored parent chain at call time
  (`flattenLocale`); currency symbols are resolved along the same chain
  (`localeSymbol`) and fall back to the ISO code (`CHF 12.34` in `en`).
- The CLDR currency-spacing rule is applied in simplified form: a
  no-break space is inserted where a letter of the substituted symbol
  would abut a digit; nothing is inserted across an intervening minus
  sign (`EUR-12.34` in `de-CH`). The sign is carried entirely by the
  pattern affixes, whose `-` placeholder becomes the locale's minus sign
  (U+2212 in `sv`/`fi`/`nb`), so `math.MinInt` in a scale-0 currency
  formats without overflow.
- The internal `display` method takes the `formatStyle`, so kind-based
  displaying (open item 1) only needs to select `formatAccounting` and
  apply cash rounding on top.
- Documented limitations (Display doc comment): `latn` numbering system
  only, no per-currency pattern overrides, CLDR's minimum-grouping-digits
  rule is not applied, and no bidi isolation marks are emitted for RTL
  locales.
- Tests: rendering across locales and currencies (grouping variants,
  symbol fallback and resets, RTL data, `XXX` placeholder), tag
  normalization and fallback chains, `MinInt`/`MaxInt` bounds on 32- and
  64-bit platforms, and unit tests for the digit-grouping helper.

### Locale data generation (`claude/roadmap-item-pr-eiypo2`, 2026-08)

- The generator gained `internal/cldrgen/gen_locales.go`, which emits two
  new tables into `tables.go`: `localePatterns` (currency format patterns
  parsed into affixes plus grouping sizes, deduplicated) and `locales`
  (per-locale number symbols, pattern references, and currency symbols).
  The shared types (`localeData`, `localePattern`, `currencySymbol`) live
  in `gen_common.go` and are copied into `common.go`.
- Locale coverage is an explicit 33-entry allowlist (root plus major
  locales for the currencies exposed as package constants), restricted to
  the `latn` numbering system and the default-length `standard` and
  `accounting` patterns; compact ("short") notation and non-latn digit
  systems remain out of scope.
- Inheritance is resolved at generation time along the CLDR parent chain,
  honoring the `parentLocales` exceptions in `supplementalData.xml`
  (`zh-Hant` → `root`, `nb` → `no`) that the x/text decoder ignores.
  Locales that only appear on parent chains (`en-001`, `es-419`, `no`) are
  flattened into diffs against the nearest emitted ancestor. Each emitted
  locale stores only its diffs plus its parent index, so the runtime
  formatter can walk the same chain; the root locale sits at index 0 and
  defines every field.
- CLDR's `↑↑↑` "same as parent" markers are treated as absent, `alt`
  variants (narrow/variant symbols) are skipped, and currency symbols are
  limited to the currencies exposed as package constants. Quirks in the
  pinned CLDR 40 data (e.g. Russian's draft `XXXX` symbol for `XXX`) are
  reproduced faithfully rather than patched by hand.
- The currency table portion of `tables.go` is byte-for-byte unchanged.
- New tests in `tables_test.go` cover table invariants (root completeness,
  sorted names, parent chains terminating at root, one `¤` placeholder per
  sign) and resolved spot checks against known CLDR 40 values (de-CH
  separators and `¤-` negative pattern, en-IN 3/2 grouping, sv/nb U+2212
  minus, en-001 flattening, symbol resets such as de-CH showing EUR as
  `EUR`).

### Modernization pass (`claude/currency-repo-review-roadmap-d5k701`, 2026-08)

All changes below were made during the 2026-08 review pass.

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
