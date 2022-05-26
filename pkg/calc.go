package money

import (
	"math"
)

type calculator struct{}

var calc = calculator{}

func (c *calculator) add(a, b int64) (int64, bool) {
	z := a + b
	if (z > a) == (b > 0) {
		return z, true
	}
	return z, false
}

func (c *calculator) sub(a, b int64) (int64, bool) {
	z := a - b
	if (z < a) == (b > 0) {
		return z, true
	}
	return z, false
}

func (c *calculator) mul(a Value, m int64) (int64, bool) {
	if a == 0 || m == 0 {
		return 0, true
	}
	z := a * m
	if (z < 0) == ((a < 0) != (m < 0)) {
		if z/m == a {
			return z, true
		}
	}
	return z, false
}

func (c *calculator) div(a Value, d int64) (int64, bool) {
	q, _, ok := c.quotient(a, d)
	return q, ok
}

func (c *calculator) mod(a Value, d int64) (int64, bool) {
	_, r, ok := c.quotient(a, d)
	return r, ok
}

func (c *calculator) quotient(a Value, b int64) (int64, int64, bool) {
	if b == 0 {
		return 0, 0, false
	}
	z := a / b
	return z, a % b, (z < 0) == ((a < 0) != (b < 0))
}

func (c *calculator) alloc(a Value, r, s int64) (int64, bool) {
	if r > s {
		return a, true
	}
	z1, ok := c.mul(a, r)
	if !ok {
		return 0, false
	}
	z2, ok := c.div(z1, s)
	if !ok {
		return 0, false
	}
	return z2, true
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

func (c *calculator) round(a int64, e int) (int64, bool) {
	if a == 0 {
		return a, true
	}
	absa := c.abs(a)
	exp := int64(math.Pow(10, float64(e)))
	m, ok := c.mod(absa, exp)
	if !ok {
		return 0, false
	}
	if m > (exp / 2) {
		absa += exp
	}
	absa = (absa / exp) * exp
	if a < 0 {
		return -absa, true
	} else {
		return absa, true
	}
}
