package gen

import (
	"log"
	"os"

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

func openZip(zipPath string) *os.File {
	zip, err := os.Open(zipPath)
	if err != nil {
		log.Fatalf("error opening file %s: %v", zipPath, err)
	}
	return zip
}
