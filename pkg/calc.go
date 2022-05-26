package money

type calculator struct{}

var calc = calculator{}

const (
	min = -1 << 63
	max = 1<<63 - 1
)

func (c *calculator) add(a, b int64) (int64, bool) {
	if b > 0 {
		if a > max-b {
			return 0, false
		}
	} else {
		if b < 0 {
			if a < min+b {
				return 0, false
			}
		}
	}
	return a + b, true
}

func (c *calculator) sub(a, b int64) (int64, bool) {
	if b > 0 {
		if a < min+b {
			return 0, false
		}
	} else {
		if a > max+b {
			return 0, false
		}
	}
	return a - b, true
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
