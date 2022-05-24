package money

// TODO: structure is to be clarified
type Region struct {
	language string
	code     string
	from     string
	to       string
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
