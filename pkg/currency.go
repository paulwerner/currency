package money

type Currency struct {
	index uint16
}

func (c *Currency) Code() string {
	panic("not implemented")
}

func (c *Currency) CodeNumeric() uint8 {
	panic("not implemented")
}

func (c *Currency) String() string {
	panic("not implemented")
}

func (c *Currency) template(r Region) string {
	panic("not implemented")
}
