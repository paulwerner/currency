package money

// TODO: structure is to be clarified
type Region struct {
	region string
	code   string // BCP 47 code, e.g., en, de, fr
	from   string
	to     string
}

// String returns the Region's BCP 47 representations
// it returns "ZZ" for an unspecified region
func (r Region) String() string {
	// TODO: fetch from index
	if r.code == "" {
		return "ZZ"
	}
	return r.code
}
