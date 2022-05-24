package money

import "math"

type calculator struct{}

var calc = calculator{}

func (c *calculator) add(a, b int64) int64 {
	return a + b
}

func (c *calculator) sub(a, b int64) int64 {
	return a - b
}

func (c *calculator) mul(a int64, m int64) int64 {
	return a * m
}

func (c *calculator) div(a int64, d int64) int64 {
	return a / d
}

func (c *calculator) mod(a int64, d int64) int64 {
	return a % d
}

func (c *calculator) alloc(a int64, r, s int) int64 {
	if r > s {
		return a
	}
	return a * int64(r) / int64(s)
}

func (c *calculator) abs(a int64) int64 {
	if a < 0 {
		return -a
	}
	return a
}

func (c *calculator) neg(a int64) int64 {
	if a > 0 {
		return -a
	}
	return a
}

func (c *calculator) round(a int64, e int) int64 {
	if a == 0 {
		return a
	}
	absa := c.abs(a)
	exp := int64(math.Pow(10, float64(e)))
	m := c.mod(absa, exp)

	if m > (exp / 2) {
		absa += exp
	}
	absa = (absa / exp) * exp
	if a < 0 {
		return -absa
	} else {
		return absa
	}
}
