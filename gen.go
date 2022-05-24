package main

//go:generate go run gen.go gen_common.go

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	cw "github.com/paulwerner/gocodewriter"
	"github.com/paulwerner/gomoney/internal/tag"
	"golang.org/x/text/unicode/cldr"
)

var (
	outputFile = flag.String("out", "pkg/tables.go",
		"the file to which the tables should be written")
)

func main() {
	flag.Parse()

	zip := openCLDRCoreZip()
	d := &cldr.Decoder{}
	d.SetDirFilter("supplemental", "main")
	d.SetSectionFilter("numbers")
	db, err := d.DecodeZip(zip)
	if err != nil {
		log.Fatalf("DecodeZip: %v", err)
	}

	cw.Repackage("gen_common.go", "./pkg/common.go", "money")
	w := cw.NewWriter()
	defer w.WriteGoFile(*outputFile, "money")

	fmt.Fprintln(w, `import "github.com/paulwerner/gomoney/internal/tag"`)

	b := builder{}
	b.genCurrencies(w, db.Supplemental())
}

var constants = []string{
	// Undefined and testing.
	"XXX", "XTS",
	// G11 currencies https://en.wikipedia.org/wiki/G10_currencies.
	"USD", "EUR", "JPY", "GBP", "CHF", "AUD", "NZD", "CAD", "SEK", "NOK", "DKK",
	// Precious metals.
	"XAG", "XAU", "XPT", "XPD",

	// Additional common currencies as defined by CLDR.
	"BRL", "CNY", "INR", "RUB", "HKD", "IDR", "KRW", "MXN", "PLN", "SAR",
	"THB", "TRY", "TWD", "ZAR",
}

type builder struct {
	currencies    tag.Index
	numCurrencies int
}

func (b *builder) genCurrencies(w *cw.Writer, db *cldr.SupplementalData) {
	// 3-letter ISO currency codes
	// Start with dummy to let index start at 1.
	currencies := []string{"\x00\x00\x00\x00"}

	// currency codes
	for _, reg := range db.CurrencyData.Region {
		for _, cur := range reg.Currency {
			currencies = append(currencies, cur.Iso4217)
		}
	}
	sort.Strings(currencies)
	// Unique the elements
	k := 0
	for i := 1; i < len(currencies); i++ {
		if currencies[k] != currencies[i] {
			currencies[k+1] = currencies[i]
			k++
		}
	}
	currencies = currencies[:k+1]

	// close with dummy
	currencies = append(currencies, "\xff\xff\xff\xff")

	// Write currency values
	fmt.Fprintln(w, "const (")
	for _, c := range constants {
		idx := sort.SearchStrings(currencies, c)
		fmt.Fprintf(w, "\t%s = %d\n", strings.ToLower(c), idx)
	}
	fmt.Fprint(w, ")")

	// compute currency-related data that we merge into the table
	for _, info := range db.CurrencyData.Fractions[0].Info {
		if info.Iso4217 == "DEFAULT" {
			continue
		}
		standard := getRoundingIndex(info.Digits, info.Rounding, 0)
		cash := getRoundingIndex(info.CashDigits, info.CashRounding, standard)

		idx := sort.SearchStrings(currencies, info.Iso4217)
		currencies[idx] += mkCurrencyInfo(standard, cash)
	}

	// Set default values for entries that weren't touched
	for i, c := range currencies {
		if len(c) == 3 {
			currencies[i] += mkCurrencyInfo(0, 0)
		}
	}
	b.currencies = tag.Index(strings.Join(currencies, ""))
	w.WriteComment(`
	currency holds an alphabetically sorted list of canonical 3-letter currency
	identifiers. Each identifier is followed by a byte of type currencyInfo,
	defined in gen_common.go.`)
	w.WriteConst("currency", b.currencies)
	b.numCurrencies = (len(b.currencies) / 4) - 2
	w.WriteConst("numCurrencies", b.numCurrencies)
}

func mkCurrencyInfo(standard, cash int) string {
	return (string([]byte{byte(cash<<cashShift | standard)}))
}

func getRoundingIndex(digits, rounding string, defIdx int) int {
	round := roundings[defIdx] // default
	if digits != "" {
		round.scale = parseUint8(digits)
	}
	if rounding != "" && rounding != "0" { // 0 means 1 in CLDR
		round.increment = parseUint8(rounding)
	}

	// panic if entry doesn't exists
	for i, r := range roundings {
		if r == round {
			return i
		}
	}
	log.Fatalf("error rounding entry %#v does not exists", round)
	panic("unreachable")
}

func parseUint8(str string) uint8 {
	x, err := strconv.ParseUint(str, 10, 8)
	if err != nil {
		// Show line number of where this function was called.
		log.New(os.Stderr, "", log.Lshortfile).Output(2, err.Error())
		os.Exit(1)
	}
	return uint8(x)
}

func openCLDRCoreZip() *os.File {
	wd, err := os.Getwd()
	if err != nil {
		log.Fatalf("error opening working directory: %v", err)
	}

	zip, err := os.Open(filepath.Join(wd, "core.zip"))
	if err != nil {
		log.Fatalf("error opening core.zip file: %v", err)
	}
	return zip
}
