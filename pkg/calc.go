package money

type calculator struct{}

var calc = calculator{}

func (c *calculator) add(a, b int64) (int64, bool) {
	panic("not implemented")
}

func (c *calculator) sub(a, b int64) (int64, bool) {
	panic("not implemented")
}

func (c *calculator) mul(a int64, m int64) (int64, bool) {
	panic("not implemented")
}

func (c *calculator) div(a int64, d int64) (int64, bool) {
	panic("not implemented")
}

func (c *calculator) mod(a int64, d int64) int64 {
	panic("not implemented")
}

func (c *calculator) alloc(a int64, r, s int) (int64, bool) {
	panic("not implemented")
}

func (c *calculator) neg(a int64) int64 {
	panic("not implemented")
}

func (c *calculator) abs(a int64) int64 {
	panic("not implemented")
}

func (c *calculator) round(a int64, s, i int) int64 {
	panic("not implemented")
}
