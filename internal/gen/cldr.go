package gen

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"golang.org/x/text/unicode/cldr"
)

func GetCLDRDataFromZip(zipPath string) *cldr.CLDR {
	dc := &cldr.Decoder{}
	dc.SetDirFilter("supplemental", "main")
	dc.SetSectionFilter("numbers")

	zip := openZip(zipPath)
	defer zip.Close()

	data, err := dc.DecodeZip(zip)
	if err != nil {
		log.Fatalf("error decoding cldr core.zip: %v", err)
	}
	data.SetDraftLevel(cldr.Contributed, false)
	return data
}

type country struct {
	code string
	name string

	currencies []string
}

type currency struct {
	code         string
	digits       uint
	rounding     uint
	cashDigits   uint
	cashRounding uint

	countries []string
}

func GetCountries(db *cldr.CLDR) *map[string]*country {
	ldml := db.RawLDML("en")
	supd := db.Supplemental()

	territories := ldml.LocaleDisplayNames.Territories.Territory
	countries := make(map[string]*country, len(territories))
	for _, ter := range territories {
		// Skip alt="short", alt="variant"
		if len(ter.Alt) > 0 {
			continue
		}
		// Skip continents
		if len(ter.Type) != 2 {
			continue
		}
		// Remove "QO" (in CLDR, but invalid ISO-3166)
		if ter.Type[0] == 'Q' && ter.Type[1] >= 'M' {
			continue
		}
		if ter.Type == "ZZ" {
			continue
		}
		// Remove "X?" except "XK" (special code for Kosovo)
		if ter.Type[0] == 'X' && ter.Type[1] != 'K' {
			continue
		}
		countries[ter.Type] = &country{code: ter.Type, name: ter.Data()}
	}

	for _, r := range supd.CurrencyData.Region {
		c := countries[r.Iso3166]
		if c == nil {
			continue
		}

		for _, cu := range r.Currency {
			if cu.Tender == "false" {
				continue
			}

			// keep only current currencies
			if len(cu.To) > 0 {
				continue
			}
			c.currencies = append(c.currencies, cu.Iso4217)
		}
		switch len(c.currencies) {
		case 1:
		case 0:
			log.Printf("%s: no currencies", c.code)
		default:
			log.Printf("%s.Currencies: %v", c.code, c.currencies)
		}
	}

	fmt.Println(len(countries), "countries.")
	return &countries
}

func GetCurrencies(db *cldr.CLDR) *map[string]*currency {
	currencies := make(map[string]*currency)
	supld := db.Supplemental()
	for _, frac := range supld.CurrencyData.Fractions[0].Info {
		if len(frac.Iso4217) != 3 {
			continue
		}
		var c currency
		c.code = frac.Iso4217
		c.digits = atou(frac.Digits, 2)
		c.rounding = atou(frac.Rounding, 0)
		c.cashDigits = atou(frac.CashDigits, c.digits)
		c.cashRounding = atou(frac.CashRounding, c.rounding)
		currencies[c.code] = &c
	}

	for _, r := range supld.CurrencyData.Region {
		for _, cu := range r.Currency {
			// take only tender currencies
			if cu.Tender == "false" {
				continue
			}
			// keep only current currencies
			if len(cu.To) > 0 {
				continue
			}
			c := currencies[cu.Iso4217]
			if c == nil {
				c = &currency{
					code:         cu.Iso4217,
					digits:       2,
					rounding:     0,
					cashDigits:   2,
					cashRounding: 0,
				}
				currencies[c.code] = c
			}
			c.countries = append(c.countries, r.Iso3166)
		}
	}
	return &currencies
}

func openZip(zipPath string) *os.File {
	zip, err := os.Open(zipPath)
	if err != nil {
		log.Fatalf("error opening file %s: %v", zipPath, err)
	}
	return zip
}

func atou(s string, defaul uint) uint {
	if len(s) == 0 {
		return defaul
	}
	u, err := strconv.ParseUint(s, 10, 8)
	if err != nil {
		panic(err)
	}
	return uint(u)
}
