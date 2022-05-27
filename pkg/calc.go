package money

type calculator struct{}

var calc = calculator{}

const (
	min = -1 << 63
	max = 1<<63 - 1
)

func (c *calculator) add(a, b *Amount) (*Amount, bool) {
	panic("not implemented")
}

func (c *calculator) sub(a, b *Amount) (*Amount, bool) {
	panic("not implemented")
}

func (c *calculator) mul(a *Amount, m int) (*Amount, bool) {
	panic("not implemented")
}

func (c *calculator) div(a *Amount, d int) (*Amount, bool) {
	panic("not implemented")
}

func (c *calculator) mod(a *Amount, d int) *Amount {
	panic("not implemented")
}

func (c *calculator) alloc(a *Amount, r, s int) (*Amount, bool) {
	panic("not implemented")
}

func (c *calculator) neg(a *Amount) *Amount {
	panic("not implemented")
}

func (c *calculator) abs(a *Amount) *Amount {
	panic("not implemented")
}

func (c *calculator) round(a *Amount, s, i int) *Amount {
	panic("not implemented")
}
