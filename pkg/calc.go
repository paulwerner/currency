package money

import "math"

type calculator struct{}

var calc = calculator{}

func (c *calculator) add(a, b Value) Value {
	return a + b
}

func (c *calculator) sub(a, b Value) Value {
	return a - b
}

func (c *calculator) mul(a Value, m Value) Value {
	return a * m
}

func (c *calculator) div(a Value, d Value) Value {
	return a / d
}

func (c *calculator) mod(a Value, d Value) Value {
	return a % d
}

func (c *calculator) alloc(a Value, r, s int) Value {
	if r > s {
		return a
	}
	return a * Value(r) / Value(s)
}

func (c *calculator) abs(a Value) Value {
	if a < 0 {
		return -a
	}
	return a
}

func (c *calculator) neg(a Value) Value {
	if a > 0 {
		return -a
	}
	return a
}

func (c *calculator) round(a Value, e int) Value {
	if a == 0 {
		return a
	}
	absa := c.abs(a)
	exp := Value(math.Pow(10, float64(e)))
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
