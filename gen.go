package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	cw "github.com/paulwerner/gocodewriter"
	"github.com/paulwerner/gomoney/internal/tag"
	"golang.org/x/text/unicode/cldr"
)

var (
	outputFile = flag.String("out", "./pkg/tables.go",
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

	fmt.Fprintln(w, `import "github.com/paulwerner/money/internal/tag"`)

	b := builder{}
	b.genCurrencies(w, db.Supplemental())
}

var supportedCurrencies = []string{
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
	panic("not implemented")
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
