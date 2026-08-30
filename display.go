package currency

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"
)

// Display formats the amount for the given locale using the locale's
// standard currency format pattern from CLDR, e.g. "€19.99" for "en" and
// "19,99 €" for "de".
//
// The locale is a BCP 47 tag matched case-insensitively; both "-" and "_"
// subtag separators are accepted ("de-CH", "de_ch"). Lookup tries the
// exact tag first, then truncates subtags ("de-CH-1996" → "de-CH" → "de"),
// and falls back to the root locale for tags not in the generated tables.
// The currency symbol is resolved along the locale's CLDR parent chain and
// falls back to the ISO code when no symbol is defined; where a letter of
// the symbol would abut a digit, a no-break space is inserted (CLDR's
// currency-spacing rule in simplified form), e.g. "CHF 12.34" for "en".
//
// Limitations: only the "latn" numbering system is supported, per-currency
// pattern overrides and CLDR's minimum-grouping-digits rule are not
// applied, and no bidi isolation marks are emitted for RTL locales. No
// likely-subtags mapping is applied either: truncation sends "zh-TW" to
// "zh" (Simplified), not to "zh-Hant" — pass the script subtag explicitly
// where it matters.
func (a *Amount) Display(locale string) string {
	return a.display(lookupLocale(locale), formatStandard)
}

// display renders the amount using the locale at the given table index
// and the pattern selected by the format style
func (a *Amount) display(idx int, style formatStyle) string {
	loc := flattenLocale(idx)
	p := loc.standard
	if style == formatAccounting {
		p = loc.accounting
	}

	number := a.formatDigits(p, loc.decimal, loc.group)
	prefix, suffix := p.posPrefix, p.posSuffix
	if a.value < 0 {
		prefix, suffix = p.negPrefix, p.negSuffix
	}
	symbol := localeSymbol(idx, a.currency.Code())
	prefix = expandAffix(prefix, symbol, loc.minus)
	suffix = expandAffix(suffix, symbol, loc.minus)

	var b strings.Builder
	b.WriteString(prefix)
	// simplified CLDR currency spacing: the number always starts and
	// ends with a digit, so only the affix side needs checking
	if r, _ := utf8.DecodeLastRuneInString(prefix); unicode.IsLetter(r) {
		b.WriteString("\u00a0")
	}
	b.WriteString(number)
	if r, _ := utf8.DecodeRuneInString(suffix); unicode.IsLetter(r) {
		b.WriteString("\u00a0")
	}
	b.WriteString(suffix)
	return b.String()
}

// formatDigits renders the absolute decimal value with the locale's
// separators: the integer digits grouped as the pattern prescribes, the
// fractional digits zero-padded to the currency's standard scale
func (a *Amount) formatDigits(p *localePattern, decimal, group string) string {
	scale, _ := Standard.Rounding(a.currency)
	exp := 1
	if scale > 0 {
		var ok bool
		if exp, ok = pow(10, scale); !ok {
			// unreachable: the generated roundings table caps scales
			// at 4, so pow(10, scale) fits even in 32-bit int. The
			// fallback renders minor units as whole units, corrupting
			// the magnitude — if scales ever come from another source,
			// it must not stay silent
			scale, exp = 0, 1
		}
	}
	whole, frac := a.value/exp, a.value%exp
	if frac < 0 {
		// |frac| < exp, safe to negate
		frac = -frac
	}
	// strconv renders the sign, which is carried by the pattern affixes
	// instead; trimming it also sidesteps negating loBound, which would
	// overflow in a scale-0 currency
	digits := strings.TrimPrefix(strconv.Itoa(whole), "-")

	var b strings.Builder
	b.WriteString(groupDigits(digits, group, p.primGroup, p.secGroup))
	if scale > 0 {
		b.WriteString(decimal)
		fmt.Fprintf(&b, "%0*d", scale, frac)
	}
	return b.String()
}

// groupDigits inserts the group separator into an unsigned digit string.
// prim is the size of the rightmost digit group and sec the size of every
// group to its left ("#,##0.00" is 3/3, the Indian "#,##,##0.00" is 3/2);
// a prim of 0 means the pattern does not group digits
func groupDigits(digits, sep string, prim, sec uint8) string {
	p := int(prim)
	if p == 0 || len(digits) <= p || sep == "" {
		return digits
	}
	s := int(sec)
	if s == 0 {
		s = p
	}
	var b strings.Builder
	head := (len(digits)-p-1)%s + 1
	b.WriteString(digits[:head])
	for i := head; i < len(digits)-p; i += s {
		b.WriteString(sep)
		b.WriteString(digits[i : i+s])
	}
	b.WriteString(sep)
	b.WriteString(digits[len(digits)-p:])
	return b.String()
}

// expandAffix substitutes the affix placeholders: "-" with the locale's
// minus sign, then "¤" with the resolved currency symbol, in that order
// so that a symbol containing "-" is never re-substituted
func expandAffix(affix, symbol, minus string) string {
	affix = strings.ReplaceAll(affix, "-", minus)
	return strings.ReplaceAll(affix, "¤", symbol)
}

// flatLocale is one locale's number-formatting data with the parent
// chain's inheritance already applied, so every field is populated
type flatLocale struct {
	decimal, group, minus string
	standard, accounting  *localePattern
}

// flattenLocale overlays the locale at the given table index over its
// parent chain. The root locale defines every field (an invariant checked
// by the table tests), so the result is fully populated
func flattenLocale(idx int) flatLocale {
	var f flatLocale
	for i := idx; ; i = int(locales[i].parent) {
		l := &locales[i]
		if f.decimal == "" {
			f.decimal = l.decimal
		}
		if f.group == "" {
			f.group = l.group
		}
		if f.minus == "" {
			f.minus = l.minus
		}
		if f.standard == nil && l.standard != 0 {
			f.standard = &localePatterns[l.standard-1]
		}
		if f.accounting == nil && l.accounting != 0 {
			f.accounting = &localePatterns[l.accounting-1]
		}
		if i == 0 {
			return f
		}
	}
}

// localeSymbol resolves the symbol a locale displays the given ISO 4217
// code with by walking the locale's parent chain; it falls back to the
// ISO code itself when no locale in the chain defines a symbol
func localeSymbol(idx int, code string) string {
	for i := idx; ; i = int(locales[i].parent) {
		syms := locales[i].symbols
		j := sort.Search(len(syms), func(k int) bool { return syms[k].code >= code })
		if j < len(syms) && syms[j].code == code {
			return syms[j].symbol
		}
		if i == 0 {
			return code
		}
	}
}

// lookupLocale maps a BCP 47 tag to an index into the locales table:
// exact match on the canonicalized tag first, then repeatedly truncating
// the last subtag, then the root locale at index 0
func lookupLocale(tag string) int {
	name := canonicalTag(tag)
	for name != "" {
		if i := findLocale(name); i >= 0 {
			return i
		}
		j := strings.LastIndexByte(name, '-')
		if j < 0 {
			break
		}
		name = name[:j]
	}
	return 0
}

// findLocale returns the table index of the locale with the given
// canonical name, or -1. The table is sorted by name after the root
// entry at index 0
func findLocale(name string) int {
	if name == "root" {
		return 0
	}
	rest := locales[1:]
	i := sort.Search(len(rest), func(k int) bool { return rest[k].name >= name })
	if i < len(rest) && rest[i].name == name {
		return i + 1
	}
	return -1
}

// canonicalTag normalizes a locale tag to the form used by the locales
// table: subtags separated by "-", the language subtag lower case, script
// subtags in title case and region subtags upper case, so "de_ch" and
// "DE-CH" both become "de-CH". Unrecognized subtag shapes are kept as-is;
// they never match and are truncated away by the lookup
func canonicalTag(tag string) string {
	parts := strings.FieldsFunc(tag, func(r rune) bool {
		return r == '-' || r == '_'
	})
	for i, p := range parts {
		switch {
		case i == 0:
			parts[i] = strings.ToLower(p)
		case len(p) == 4 && isASCIIAlpha(p):
			parts[i] = strings.ToUpper(p[:1]) + strings.ToLower(p[1:])
		case len(p) == 2 && isASCIIAlpha(p):
			parts[i] = strings.ToUpper(p)
		}
	}
	return strings.Join(parts, "-")
}

func isASCIIAlpha(s string) bool {
	for i := 0; i < len(s); i++ {
		c := s[i] | 0x20 // fold to lower case
		if c < 'a' || c > 'z' {
			return false
		}
	}
	return true
}
