package money

const (
	// 32 or 64
	intSize = 32 << (^uint(0) >> 63)

	// 32 bit: -2147483648
	// 64bit: -9223372036854775808
	loBound int = -1 << (intSize - 1)

	// 32 bit: 2147483647
	// 64bit: 9223372036854775807
	hiBound int = 1<<(intSize-1) - 1
)

func add(x, y amount) (amount, bool) {
	if y > 0 {
		if x > hiBound-y {
			return 0, false
		}
	} else {
		if x < loBound-y {
			return 0, false
		}
	}
	return x + y, true
}

func sub(x, y amount) (amount, bool) {
	if y > 0 {
		if x < loBound+y {
			return 0, false
		}
	} else {
		if x > hiBound+y {
			return 0, false
		}
	}
	return x - y, true
}

func mul(x amount, m int) (amount, bool) {
	if x == 0 || m == 0 {
		return 0, true
	}
	if (m > 0 && x > hiBound/m) || (m < 0 && x < hiBound/m) {
		return 0, false
	}
	if (m > 0 && x < loBound/m) || (m < -1 && x > loBound/m) {
		return 0, false
	}
	return x * m, true
}

func div(x amount, d int) (amount, bool) {
	if d == 0 ||
		(x == loBound && d == -1) {
		return 0, false
	}
	return x / d, true
}

func mod(x amount, d int) (amount, bool) {
	if d == 0 || (x == loBound && d == -1) {
		return 0, false
	}
	return x % d, true
}

func alloc(x amount, r, s int) (amount, bool) {
	if r < 0 || s <= 0 {
		return 0, false
	}
	if r > s {
		return 0, false
	}
	if r == 0 {
		return 0, true
	}

	p, ok := mul(x, r)
	if !ok {
		return 0, false
	}
	z, ok := div(p, s)
	if !ok {
		return 0, false
	}
	return z, true
}

func neg(x amount) amount {
	if x > 0 {
		return -x
	}
	return x
}

func abs(x amount) (amount, bool) {
	if x < 0 {
		if x == loBound {
			return 0, false
		}
		return -x, true
	}
	return x, true
}

// pow computes x**e using binary powering algorithm
// for a positive exponent
// see Donald Knuth: The Art of Computer Programming
func pow(x, e int) (amount, bool) {
	if e < 0 {
		return 0, false
	}
	p := 1
	for e > 0 {
		if e&1 != 0 {
			r, ok := mul(p, x)
			if !ok {
				return 0, false
			}
			p = r
		}
		e >>= 1
		r, ok := mul(x, x)
		if !ok {
			return 0, false
		}
		x = r
	}
	return p, true
}

func round(x amount, s, i int) (amount, bool) {
	if x == 0 {
		return 0, true
	}
	if s < 0 {
		return 0, false
	}
	xabs, ok := abs(x)
	if !ok {
		return 0, false
	}

	exp, ok := pow(10, s)
	if !ok {
		return 0, false
	}

	m, ok := mod(xabs, exp)
	if !ok {
		return 0, false
	}

	if m >= (exp / 2) {
		xabs, ok = add(xabs, exp)
		if !ok {
			return 0, false
		}
	}
	q, ok := div(xabs, exp)
	if !ok {
		return 0, false
	}
	xabs, ok = mul(q, exp)
	if !ok {
		return 0, false
	}
	if x < 0 {
		return -xabs, true
	} else {
		return xabs, true
	}
}
