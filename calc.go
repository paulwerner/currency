package currency

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

func add(x, y int) (int, bool) {
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

func sub(x, y int) (int, bool) {
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

func mul(x int, m int) (int, bool) {
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

func div(x int, d int) (int, bool) {
	if d == 0 ||
		(x == loBound && d == -1) {
		return 0, false
	}
	return x / d, true
}

func mod(x int, d int) (int, bool) {
	if d == 0 || (x == loBound && d == -1) {
		return 0, false
	}
	return x % d, true
}

func alloc(x int, r, s int) (int, bool) {
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

func neg(x int) (int, bool) {
	if x == loBound {
		return 0, false
	}
	return -x, true
}

func abs(x int) (int, bool) {
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
func pow(x, e int) (int, bool) {
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
		if e == 0 {
			// p is complete; squaring x once more could report a
			// false overflow for a representable result
			break
		}
		r, ok := mul(x, x)
		if !ok {
			return 0, false
		}
		x = r
	}
	return p, true
}

// round rounds x to the nearest multiple of step, with ties rounded away
// from zero (i.e. "half up" in terms of magnitude). step must be positive.
func round(x int, step int) (int, bool) {
	if step <= 0 {
		return 0, false
	}
	if x == 0 || step == 1 {
		return x, true
	}
	xabs, ok := abs(x)
	if !ok {
		return 0, false
	}

	m := xabs % step
	// m >= step-m is equivalent to 2*m >= step without risking overflow
	if m >= step-m {
		xabs, ok = add(xabs, step-m)
		if !ok {
			return 0, false
		}
	} else {
		xabs -= m
	}
	if x < 0 {
		// xabs > 0 here, so -xabs cannot underflow
		return -xabs, true
	}
	return xabs, true
}
