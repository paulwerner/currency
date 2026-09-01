// Locale table generation: emits the locales and localePatterns tables
// consumed by locale-based formatting.
//
// CLDR data is inherited along a locale's parent chain (de_CH -> de ->
// root, with exceptions such as zh_Hant -> root listed in
// supplementalData.xml under parentLocales). The x/text cldr decoder only
// resolves truncation-based inheritance, so this file walks the chain
// itself: it resolves every locale's effective data at generation time and
// stores per-locale diffs against the nearest allowlisted ancestor, which
// keeps the tables small while letting the runtime walk the same chain.
//
// Extraction is limited to the latn numbering system and to the default
// (non-"short") currency format length; compact notation, other digit
// systems, and per-currency pattern overrides are out of scope.
package main

import (
	"fmt"
	"log"
	"sort"
	"strings"

	"github.com/paulwerner/currency/internal/gen"
	"golang.org/x/text/unicode/cldr"
)

// localeAllowlist holds the locales emitted into the locales table, as
// CLDR identifiers (underscore-separated, matching the core.zip file
// names). The list is explicit so output is deterministic regardless of
// how many locale files the core.zip contains. Locales that only appear
// on parent chains (e.g. en_001, es_419, no) need to be present in the
// core.zip but are flattened into the diffs rather than emitted.
var localeAllowlist = []string{
	"root",
	"ar", "da", "de", "de_AT", "de_CH",
	"en", "en_AU", "en_CA", "en_GB", "en_IN", "en_NZ",
	"es", "es_MX", "fi", "fr", "fr_CA", "fr_CH",
	"id", "it", "ja", "ko", "nb", "nl", "pl",
	"pt", "pt_PT", "ru", "sv", "th", "tr",
	"zh", "zh_Hant",
}

// inheritedMarker is CLDR's explicit "same as parent" value. Elements
// carrying it are treated as absent.
const inheritedMarker = "↑↑↑"

// currencySymbolPlaceholder is the "¤" placeholder substituted with the
// currency symbol at format time.
const currencySymbolPlaceholder = "¤"

// rawNumbers holds the latn number-formatting data extracted from a single
// locale file (or, after resolution, the effective data of a locale).
// Empty strings mean "not defined here".
type rawNumbers struct {
	decimal, group, minus string
	standard, accounting  string            // unparsed currency format patterns
	symbols               map[string]string // ISO 4217 code -> symbol
}

func (b *builder) genLocales(w *gen.Writer, db *cldr.CLDR) {
	overrides := parentOverrides(db.Supplemental())
	parent := func(loc string) string {
		if loc == "root" {
			return ""
		}
		if p, ok := overrides[loc]; ok {
			return p
		}
		if i := strings.LastIndex(loc, "_"); i >= 0 {
			return loc[:i]
		}
		return "root"
	}

	raw := map[string]*rawNumbers{}
	extract := func(loc string) *rawNumbers {
		if r, ok := raw[loc]; ok {
			return r
		}
		ldml := db.RawLDML(loc)
		if ldml == nil {
			log.Fatalf("error locale %q not found in core.zip", loc)
		}
		r := extractNumbers(loc, ldml)
		raw[loc] = r
		return r
	}

	resolved := map[string]*rawNumbers{}
	resolving := map[string]bool{}
	var resolve func(loc string) *rawNumbers
	resolve = func(loc string) *rawNumbers {
		if r, ok := resolved[loc]; ok {
			return r
		}
		if resolving[loc] {
			log.Fatalf("error parent chain of %q contains a cycle", loc)
		}
		resolving[loc] = true
		r := &rawNumbers{symbols: map[string]string{}}
		if loc != "root" {
			p := resolve(parent(loc))
			*r = *p
			r.symbols = make(map[string]string, len(p.symbols))
			for c, s := range p.symbols {
				r.symbols[c] = s
			}
		}
		x := extract(loc)
		overlayString(&r.decimal, x.decimal)
		overlayString(&r.group, x.group)
		overlayString(&r.minus, x.minus)
		overlayString(&r.standard, x.standard)
		overlayString(&r.accounting, x.accounting)
		for c, s := range x.symbols {
			r.symbols[c] = s
		}
		if loc == "root" {
			if r.decimal == "" || r.group == "" || r.minus == "" || r.standard == "" {
				log.Fatalf("error root locale data is incomplete: %+v", r)
			}
			if r.accounting == "" {
				// root's accounting format is an alias to the standard one
				r.accounting = r.standard
			}
		}
		resolved[loc] = r
		return r
	}

	// The emitted table starts with root at index 0 (so a zero parent
	// field means "root") followed by the remaining locales sorted by
	// their canonical name, ready for binary search.
	type outLocale struct {
		id, name string
	}
	out := []outLocale{}
	seen := map[string]bool{}
	for _, id := range localeAllowlist {
		if seen[id] {
			log.Fatalf("error duplicate locale %q in allowlist", id)
		}
		seen[id] = true
		out = append(out, outLocale{id, strings.ReplaceAll(id, "_", "-")})
	}
	sort.Slice(out, func(i, j int) bool {
		if (out[i].id == "root") != (out[j].id == "root") {
			return out[i].id == "root"
		}
		return out[i].name < out[j].name
	})
	index := map[string]int{}
	for i, l := range out {
		index[l.id] = i
	}

	// storedParent reports the nearest ancestor that is itself emitted;
	// intermediate locales (en_001, es_419, no) are flattened into diffs.
	storedParent := func(loc string) string {
		for p := parent(loc); ; p = parent(p) {
			if seen[p] {
				return p
			}
			if p == "" {
				log.Fatalf("error parent chain of %q did not reach root", loc)
			}
		}
	}

	var patterns []localePattern
	patternIndex := map[localePattern]int{}
	// patternRef returns the 1-based index of the parsed pattern,
	// interning it on first use (0 is reserved for "inherited").
	patternRef := func(s string) uint16 {
		p := parsePattern(s)
		i, ok := patternIndex[p]
		if !ok {
			i = len(patterns)
			patternIndex[p] = i
			patterns = append(patterns, p)
		}
		return uint16(i + 1)
	}

	codes := append([]string(nil), constants...)
	sort.Strings(codes)

	locales := make([]localeData, 0, len(out))
	for _, l := range out {
		r := resolve(l.id)
		d := localeData{name: l.name}
		var p *rawNumbers
		if l.id == "root" {
			p = &rawNumbers{} // diff root against nothing
		} else {
			pid := storedParent(l.id)
			d.parent = uint16(index[pid])
			p = resolve(pid)
		}
		d.decimal = diffString(r.decimal, p.decimal)
		d.group = diffString(r.group, p.group)
		d.minus = diffString(r.minus, p.minus)
		if r.standard != p.standard {
			d.standard = patternRef(r.standard)
		}
		if r.accounting != p.accounting {
			d.accounting = patternRef(r.accounting)
		}
		for _, code := range codes {
			if s := r.symbols[code]; s != "" && s != p.symbols[code] {
				d.symbols = append(d.symbols, currencySymbol{code: code, symbol: s})
			}
		}
		locales = append(locales, d)
	}

	writeLocaleTables(w, locales, patterns)
}

// parentOverrides returns the parent-chain exceptions from
// supplementalData.xml (parentLocales), e.g. zh_Hant -> root.
func parentOverrides(supp *cldr.SupplementalData) map[string]string {
	m := map[string]string{}
	if supp.ParentLocales == nil {
		return m
	}
	for _, p := range supp.ParentLocales.ParentLocale {
		for _, loc := range strings.Fields(p.Locales) {
			m[loc] = p.Parent
		}
	}
	return m
}

// overlayString sets *dst to src unless src is undefined
func overlayString(dst *string, src string) {
	if src != "" {
		*dst = src
	}
}

// diffString returns v if it differs from the parent's value, "" otherwise
func diffString(v, parent string) string {
	if v == parent {
		return ""
	}
	return v
}

// value returns the character data of a CLDR leaf element, treating the
// explicit inheritance marker as absent. It rejects elements it should
// never encounter (aliases) to fail loudly on unexpected data
func value(loc string, c *cldr.Common) string {
	if c.Alias != nil {
		log.Fatalf("error %s: unexpected alias element", loc)
	}
	if c.Data() == inheritedMarker {
		return ""
	}
	return c.Data()
}

// extractNumbers pulls the latn number symbols, the default-length
// standard and accounting currency format patterns, and the currency
// symbols out of a single, unresolved locale file. Elements with an alt
// attribute (narrow/variant forms) are skipped
func extractNumbers(loc string, ldml *cldr.LDML) *rawNumbers {
	r := &rawNumbers{symbols: map[string]string{}}
	if ldml.Numbers == nil {
		return r
	}
	// setOnce rejects repeated non-alt elements for the same field, which
	// would otherwise silently collapse to whichever comes last in file
	// order (e.g. two draft variants after a CLDR version bump)
	seen := map[string]bool{}
	setOnce := func(field string, dst *string, v string) {
		if seen[field] {
			log.Fatalf("error %s: duplicate %s element", loc, field)
		}
		seen[field] = true
		*dst = v
	}
	for _, s := range ldml.Numbers.Symbols {
		if s.NumberSystem != "latn" {
			continue
		}
		for _, e := range s.Decimal {
			if e.Alt == "" {
				setOnce("decimal", &r.decimal, value(loc, &e.Common))
			}
		}
		for _, e := range s.Group {
			if e.Alt == "" {
				setOnce("group", &r.group, value(loc, &e.Common))
			}
		}
		for _, e := range s.MinusSign {
			if e.Alt == "" {
				setOnce("minusSign", &r.minus, value(loc, &e.Common))
			}
		}
	}
	for _, cf := range ldml.Numbers.CurrencyFormats {
		if cf.NumberSystem != "latn" {
			continue
		}
		for _, cfl := range cf.CurrencyFormatLength {
			if cfl.Type != "" {
				continue // "short" is compact notation; out of scope
			}
			for _, f := range cfl.CurrencyFormat {
				if f.Alias != nil {
					continue // e.g. root aliases accounting to standard
				}
				var dst *string
				switch f.Type {
				case "standard":
					dst = &r.standard
				case "accounting":
					dst = &r.accounting
				default:
					continue
				}
				for _, p := range f.Pattern {
					if p.Alt == "" {
						setOnce(f.Type+" pattern", dst, value(loc, &p.Common))
					}
				}
			}
		}
	}
	if ldml.Numbers.Currencies != nil {
		for _, c := range ldml.Numbers.Currencies.Currency {
			var sym string
			for _, s := range c.Symbol {
				if s.Alt == "" {
					setOnce("symbol for "+c.Type, &sym, value(loc, s))
				}
			}
			if sym != "" {
				r.symbols[c.Type] = sym
			}
		}
	}
	return r
}

// parsePattern parses a CLDR currency format pattern such as
// "¤#,##0.00;(¤#,##0.00)" into its affixes and grouping sizes. The number
// of fraction digits in the pattern is ignored: at format time the
// currency's own scale wins (see Kind.Rounding). If the pattern has no
// explicit negative subpattern, the negative affixes default to the
// positive ones with a minus sign prepended, per TR35
func parsePattern(s string) localePattern {
	sub := splitPattern(s)
	pre, suf, prim, sec := parseSubpattern(s, sub[0])
	p := localePattern{
		posPrefix: pre, posSuffix: suf,
		primGroup: prim, secGroup: sec,
	}
	if len(sub) == 1 {
		p.negPrefix, p.negSuffix = "-"+pre, suf
	} else {
		// per TR35 only the affixes of the negative subpattern matter;
		// its number section is ignored
		p.negPrefix, p.negSuffix, _, _ = parseSubpattern(s, sub[1])
	}
	for _, affixes := range [][2]string{
		{p.posPrefix, p.posSuffix},
		{p.negPrefix, p.negSuffix},
	} {
		if n := strings.Count(affixes[0]+affixes[1], currencySymbolPlaceholder); n != 1 {
			log.Fatalf("error pattern %q: want 1 currency placeholder per subpattern, got %d", s, n)
		}
	}
	return p
}

// splitPattern splits a pattern into its positive and optional negative
// subpattern at an unquoted ";"
func splitPattern(s string) []string {
	inQuote := false
	for i, r := range s {
		switch {
		case r == '\'':
			inQuote = !inQuote
		case r == ';' && !inQuote:
			return []string{s[:i], s[i+1:]}
		}
	}
	return []string{s}
}

// parseSubpattern splits one subpattern into the literal prefix and suffix
// around the number section and derives the grouping sizes from it.
// Quoting is resolved in the affixes (see unquote)
func parseSubpattern(pat, sub string) (prefix, suffix string, prim, sec uint8) {
	const numberChars = "#0,."
	start, end := -1, -1
	inQuote := false
	for i, r := range sub {
		if r == '\'' {
			inQuote = !inQuote
			continue
		}
		if inQuote {
			continue
		}
		if strings.ContainsRune(numberChars, r) {
			if start < 0 {
				start = i
			}
			end = i + 1
		} else if r >= '1' && r <= '9' || strings.ContainsRune("E@*%‰", r) {
			// explicit digits, significant digits, exponents, the padding
			// escape, and percent/per-mille scaling are TR35 pattern
			// syntax this parser does not implement
			log.Fatalf("error pattern %q: unsupported pattern syntax character %q", pat, r)
		}
	}
	if start < 0 {
		log.Fatalf("error pattern %q: subpattern %q has no number section", pat, sub)
	}
	number := sub[start:end]
	if strings.Trim(number, numberChars) != "" {
		log.Fatalf("error pattern %q: number section %q is not contiguous", pat, number)
	}
	prim, sec = groupSizes(pat, number)
	return unquote(pat, sub[:start]), unquote(pat, sub[end:]), prim, sec
}

// groupSizes derives the primary and secondary integer grouping sizes from
// a pattern's number section: "#,##0.00" is 3/3, "#,##,##0.00" is 3/2 and
// "0.00" is 0/0 (no grouping)
func groupSizes(pat, number string) (prim, sec uint8) {
	integer := number
	if i := strings.IndexByte(integer, '.'); i >= 0 {
		integer = integer[:i]
	}
	groups := strings.Split(integer, ",")
	if len(groups) == 1 {
		return 0, 0
	}
	for _, g := range groups[1:] {
		if g == "" {
			log.Fatalf("error pattern %q: empty digit group", pat)
		}
	}
	prim = uint8(len(groups[len(groups)-1]))
	sec = prim
	if len(groups) > 2 {
		sec = uint8(len(groups[len(groups)-2]))
	}
	return prim, sec
}

// unquote resolves pattern quoting in an affix: text between single quotes
// is literal, and a doubled single quote is a literal apostrophe. A quoted
// minus sign or currency sign is rejected: the stored affixes treat "-"
// and "¤" as placeholders, so they cannot express "literal, do not
// substitute" for those characters (no CLDR 40 pattern in the allowlist
// quotes them)
func unquote(pat, s string) string {
	var b strings.Builder
	inQuote := false
	rs := []rune(s)
	for i := 0; i < len(rs); i++ {
		if rs[i] == '\'' {
			if i+1 < len(rs) && rs[i+1] == '\'' {
				b.WriteRune('\'')
				i++
				continue
			}
			inQuote = !inQuote
			continue
		}
		if inQuote && (rs[i] == '-' || string(rs[i]) == currencySymbolPlaceholder) {
			log.Fatalf("error pattern %q: quoted %q in affix %q would lose its literal meaning", pat, string(rs[i]), s)
		}
		b.WriteRune(rs[i])
	}
	if inQuote {
		log.Fatalf("error pattern %q: unbalanced quote in affix %q", pat, s)
	}
	return b.String()
}

func writeLocaleTables(w *gen.Writer, locales []localeData, patterns []localePattern) {
	w.WriteComment(`
	localePatterns holds the parsed currency format patterns referenced by
	the locales table, of type localePattern, defined in gen_common.go.
	Entries are referenced by 1-based index; 0 means "inherited".`)
	fmt.Fprintln(w, "var localePatterns = [...]localePattern{")
	for _, p := range patterns {
		fmt.Fprintf(w, "\t{posPrefix: %q, posSuffix: %q, negPrefix: %q, negSuffix: %q, primGroup: %d, secGroup: %d},\n",
			p.posPrefix, p.posSuffix, p.negPrefix, p.negSuffix, p.primGroup, p.secGroup)
	}
	fmt.Fprintln(w, "}")

	w.WriteComment(`
	locales holds the latn number-formatting data of each supported locale
	as a sparse overlay over its parent locale, of type localeData, defined
	in gen_common.go. The root locale is at index 0; the remaining entries
	are sorted by name. Currency symbols are restricted to the currencies
	listed in the currency constants above.`)
	fmt.Fprintln(w, "var locales = [...]localeData{")
	for _, l := range locales {
		fmt.Fprintf(w, "\t{name: %q", l.name)
		if l.parent != 0 {
			fmt.Fprintf(w, ", parent: %d", l.parent)
		}
		for _, f := range []struct{ name, v string }{
			{"decimal", l.decimal}, {"group", l.group}, {"minus", l.minus},
		} {
			if f.v != "" {
				fmt.Fprintf(w, ", %s: %q", f.name, f.v)
			}
		}
		if l.standard != 0 {
			fmt.Fprintf(w, ", standard: %d", l.standard)
		}
		if l.accounting != 0 {
			fmt.Fprintf(w, ", accounting: %d", l.accounting)
		}
		if len(l.symbols) > 0 {
			fmt.Fprintf(w, ", symbols: []currencySymbol{")
			for i, s := range l.symbols {
				if i > 0 {
					fmt.Fprint(w, ", ")
				}
				fmt.Fprintf(w, "{%q, %q}", s.code, s.symbol)
			}
			fmt.Fprint(w, "}")
		}
		fmt.Fprintln(w, "},")
	}
	fmt.Fprintln(w, "}")
}
