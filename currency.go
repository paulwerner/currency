package money

// Currency is an ISO 4217 currency designator
type Currency struct {
	Code string
}

// Equals returns true, if both Codes are equal,
// false otherwise
func (c *Currency) Equals(oc Currency) bool {
	return c.Code == oc.Code
}

// String returns the string representation of the currency
func (c *Currency) String() string {
	return c.Code
}
