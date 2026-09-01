package currency

import (
	"sort"
	"strings"
	"testing"
)

// TestLocales_Invariants checks the structural properties of the generated
// locales table that the runtime lookup relies on
func TestLocales_Invariants(t *testing.T) {
	root := locales[0]
	if root.name != "root" || root.parent != 0 {
		t.Fatalf("expected root locale at index 0, got %q (parent %d)", root.name, root.parent)
	}
	if root.decimal == "" || root.group == "" || root.minus == "" {
		t.Errorf("expected root to define all number symbols, got %+v", root)
	}
	if root.standard == 0 || root.accounting == 0 {
		t.Errorf("expected root to define both format patterns, got %+v", root)
	}

	names := make([]string, 0, len(locales)-1)
	for _, l := range locales[1:] {
		names = append(names, l.name)
	}
	if !sort.StringsAreSorted(names) {
		t.Errorf("expected locale names to be sorted, got %v", names)
	}

	seen := map[string]bool{"root": true}
	for i, l := range locales {
		if i > 0 && seen[l.name] {
			t.Errorf("[%v]: duplicate locale name", l.name)
		}
		seen[l.name] = true
		if strings.Contains(l.name, "_") {
			t.Errorf("[%v]: expected canonical hyphenated name", l.name)
		}
		if int(l.parent) >= len(locales) {
			t.Errorf("[%v]: parent index %d out of range", l.name, l.parent)
		}
		if i > 0 && int(l.parent) == i {
			t.Errorf("[%v]: locale is its own parent", l.name)
		}
		// parent chains must terminate at root without cycles
		steps := 0
		for j := i; j != 0; j = int(locales[j].parent) {
			if steps++; steps > len(locales) {
				t.Errorf("[%v]: parent chain does not reach root", l.name)
				break
			}
		}
		if int(l.standard) > len(localePatterns) || int(l.accounting) > len(localePatterns) {
			t.Errorf("[%v]: pattern index out of range: %+v", l.name, l)
		}
		for k, s := range l.symbols {
			if k > 0 && l.symbols[k-1].code >= s.code {
				t.Errorf("[%v]: currency symbols not sorted at %q", l.name, s.code)
			}
			if s.symbol == "" {
				t.Errorf("[%v]: empty symbol for %q", l.name, s.code)
			}
			if _, err := CurrencyFromISO(s.code); err != nil {
				t.Errorf("[%v]: symbol for unknown currency %q", l.name, s.code)
			}
		}
	}
}

// TestLocalePatterns_Invariants checks that every generated pattern keeps
// exactly one currency placeholder per sign and has consistent grouping
func TestLocalePatterns_Invariants(t *testing.T) {
	const placeholder = "¤"
	for i, p := range localePatterns {
		if n := strings.Count(p.posPrefix+p.posSuffix, placeholder); n != 1 {
			t.Errorf("[%v]: want 1 positive currency placeholder, got %d", i, n)
		}
		if n := strings.Count(p.negPrefix+p.negSuffix, placeholder); n != 1 {
			t.Errorf("[%v]: want 1 negative currency placeholder, got %d", i, n)
		}
		if p.primGroup == 0 && p.secGroup != 0 {
			t.Errorf("[%v]: secondary group without primary group: %+v", i, p)
		}
	}
}

// resolveLocale walks the stored parent chain the same way the future
// runtime formatter will, overlaying each locale's diffs over its parent
func resolveLocale(t *testing.T, name string) (
	decimal, group, minus string,
	standard, accounting localePattern,
	symbols map[string]string,
) {
	t.Helper()
	idx := -1
	for i := range locales {
		if locales[i].name == name {
			idx = i
			break
		}
	}
	if idx < 0 {
		t.Fatalf("locale %q not in the locales table", name)
	}
	var chain []int
	for i := idx; ; i = int(locales[i].parent) {
		chain = append(chain, i)
		if i == 0 {
			break
		}
	}
	symbols = map[string]string{}
	for k := len(chain) - 1; k >= 0; k-- {
		l := &locales[chain[k]]
		if l.decimal != "" {
			decimal = l.decimal
		}
		if l.group != "" {
			group = l.group
		}
		if l.minus != "" {
			minus = l.minus
		}
		if l.standard != 0 {
			standard = localePatterns[l.standard-1]
		}
		if l.accounting != 0 {
			accounting = localePatterns[l.accounting-1]
		}
		for _, s := range l.symbols {
			symbols[s.code] = s.symbol
		}
	}
	return decimal, group, minus, standard, accounting, symbols
}

func TestLocales_ResolvedData(t *testing.T) {
	// ¤ is the currency symbol placeholder; \u00a0 is a no-break
	// space, \u202f a narrow no-break space
	rootStd := localePattern{
		posPrefix: "¤\u00a0", negPrefix: "-¤\u00a0",
		primGroup: 3, secGroup: 3,
	}
	enStd := localePattern{
		posPrefix: "¤", negPrefix: "-¤",
		primGroup: 3, secGroup: 3,
	}
	enAcct := localePattern{
		posPrefix: "¤", negPrefix: "(¤", negSuffix: ")",
		primGroup: 3, secGroup: 3,
	}
	deStd := localePattern{
		posSuffix: "\u00a0¤", negPrefix: "-", negSuffix: "\u00a0¤",
		primGroup: 3, secGroup: 3,
	}
	frAcct := localePattern{
		posSuffix: "\u00a0¤", negPrefix: "(", negSuffix: "\u00a0¤)",
		primGroup: 3, secGroup: 3,
	}

	for _, tc := range []struct {
		locale               string
		decimal, group       string
		minus                string
		standard, accounting localePattern
	}{
		{"root", ".", ",", "-", rootStd, rootStd},
		{"en", ".", ",", "-", enStd, enAcct},
		// en-GB inherits everything relevant from en (via en-001)
		{"en-GB", ".", ",", "-", enStd, enAcct},
		// en-IN uses Indian 3/2 digit grouping
		{"en-IN", ".", ",", "-", localePattern{
			posPrefix: "¤", negPrefix: "-¤",
			primGroup: 3, secGroup: 2,
		}, enAcct},
		{"de", ",", ".", "-", deStd, deStd},
		// de-AT overrides the group separator and the pattern, keeps ","
		{"de-AT", ",", "\u00a0", "-", rootStd, deStd},
		// de-CH uses "." and apostrophe grouping and its own pattern
		{"de-CH", ".", "’", "-", localePattern{
			posPrefix: "¤\u00a0", negPrefix: "¤-",
			primGroup: 3, secGroup: 3,
		}, deStd},
		// fr wraps negative accounting amounts in parentheses
		{"fr", ",", "\u202f", "-", deStd, frAcct},
		// fr-CH stores no diffs at all and fully inherits from fr
		{"fr-CH", ",", "\u202f", "-", deStd, frAcct},
		// sv uses U+2212 as minus sign and no-break-space grouping
		{"sv", ",", "\u00a0", "−", deStd, deStd},
		// nb inherits its symbols from "no" (flattened at generation time)
		{"nb", ",", "\u00a0", "−", localePattern{
			posPrefix: "¤\u00a0", negPrefix: "¤\u00a0-",
			primGroup: 3, secGroup: 3,
		}, localePattern{
			posPrefix: "¤\u00a0", negPrefix: "(¤\u00a0", negSuffix: ")",
			primGroup: 3, secGroup: 3,
		}},
		{"ja", ".", ",", "-", enStd, enAcct},
		{"zh-Hant", ".", ",", "-", enStd, enAcct},
	} {
		decimal, group, minus, standard, accounting, _ := resolveLocale(t, tc.locale)
		if decimal != tc.decimal {
			t.Errorf("[%v]: want decimal %q, got %q", tc.locale, tc.decimal, decimal)
		}
		if group != tc.group {
			t.Errorf("[%v]: want group %q, got %q", tc.locale, tc.group, group)
		}
		if minus != tc.minus {
			t.Errorf("[%v]: want minus %q, got %q", tc.locale, tc.minus, minus)
		}
		if standard != tc.standard {
			t.Errorf("[%v]: want standard pattern %+v, got %+v", tc.locale, tc.standard, standard)
		}
		if accounting != tc.accounting {
			t.Errorf("[%v]: want accounting pattern %+v, got %+v", tc.locale, tc.accounting, accounting)
		}
	}
}

func TestLocales_ParentChain(t *testing.T) {
	for _, tc := range []struct {
		locale string
		parent string
	}{
		{"de-AT", "de"},
		{"de-CH", "de"},
		{"en-GB", "en"},
		{"es-MX", "es"},
		{"fr-CA", "fr"},
		{"pt-PT", "pt"},
		// parentLocales exceptions from supplementalData.xml
		{"zh-Hant", "root"},
		{"nb", "root"}, // via "no", which is not emitted
	} {
		found := false
		for _, l := range locales {
			if l.name == tc.locale {
				if got := locales[l.parent].name; got != tc.parent {
					t.Errorf("[%v]: want parent %q, got %q", tc.locale, tc.parent, got)
				}
				found = true
			}
		}
		if !found {
			t.Errorf("[%v]: locale not in the locales table", tc.locale)
		}
	}
}

func TestLocales_CurrencySymbols(t *testing.T) {
	for _, tc := range []struct {
		locale string
		code   string
		want   string // "" means no symbol; the runtime falls back to the ISO code
	}{
		{"root", "USD", "US$"},
		{"root", "EUR", "€"},
		{"root", "CHF", ""},
		{"en", "USD", "$"},
		{"en", "JPY", "¥"},
		{"en", "GBP", "£"},
		{"en", "CHF", ""},
		// en-GB restores the international symbols from en-001
		{"en-GB", "USD", "US$"},
		{"en-GB", "JPY", "JP¥"},
		{"en-GB", "GBP", "£"},
		{"en-AU", "AUD", "$"},
		{"en-AU", "USD", "USD"},
		{"de", "USD", "$"},
		{"de", "EUR", "€"},
		// de-CH resets the EUR symbol back to the ISO code
		{"de-CH", "EUR", "EUR"},
		{"fr", "USD", "$US"},
		{"es-MX", "MXN", "$"},
		{"sv", "SEK", "kr"},
		{"nb", "NOK", "kr"},
		{"ru", "RUB", "₽"},
		{"tr", "TRY", "₺"},
		{"zh-Hant", "TWD", "$"},
	} {
		_, _, _, _, _, symbols := resolveLocale(t, tc.locale)
		if got := symbols[tc.code]; got != tc.want {
			t.Errorf("[%v %v]: want symbol %q, got %q", tc.locale, tc.code, tc.want, got)
		}
	}
}
