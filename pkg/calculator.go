package money

import "math"

type calculator struct{}

var calc = calculator{}

func (c *calculator) add(a, b Amount) Amount {
	return a + b
}

func (c *calculator) sub(a, b Amount) Amount {
	return a - b
}

func (c *calculator) mul(a Amount, m int64) Amount {
	return a * m
}

func (c *calculator) div(a Amount, d int64) Amount {
	return a / d
}

func (c *calculator) mod(a Amount, d int64) Amount {
	return a % d
}

func (c *calculator) alloc(a Amount, r, s int) Amount {
	if r > s {
		return a
	}
	return a * int64(r) / int64(s)
}

func (c *calculator) abs(a Amount) Amount {
	if a < 0 {
		return -a
	}
	return a
}

func (c *calculator) neg(a Amount) Amount {
	if a > 0 {
		return -a
	}
	return a
}

func (c *calculator) round(a Amount, e int) Amount {
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
