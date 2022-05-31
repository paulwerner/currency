package money

const (
	// 32 or 64
	intSize = 32 << (^uint(0) >> 63)

	// 32Bit: 2147483648
	// 64bit: 9223372036854775808
	loBound uint = -(-1 << (intSize - 1))

	// 32 bit: 2147483647
	// 64bit: 9223372036854775807
	hiBound uint = 1<<(intSize-1) - 1

	// 32 bit: -2147483648
	// 64bit: -9223372036854775808
	// loBound2 int  = -1 << (intSize - 1)
	loBound2 int = -1 << (intSize - 1)
	// 32 bit: 2147483647
	// 64bit: 9223372036854775807
	hiBound2 int = 1<<(intSize-1) - 1
)

type calculator struct{}

var calc = calculator{}

func (c *calculator) add(x, y value2) (value2, bool) {
	if y > 0 {
		if x > hiBound2-y {
			return 0, false
		}
	} else {
		if x < loBound2-y {
			return 0, false
		}
	}
	return x + y, true
}

func (c *calculator) sub(x, y value2) (value2, bool) {
	if y > 0 {
		if x < loBound2+y {
			return 0, false
		}
	} else {
		if x > hiBound2+y {
			return 0, false
		}
	}
	return x - y, true
}

func (c *calculator) mul(x value2, m int) (value2, bool) {
	if x == 0 || m == 0 {
		return 0, true
	}
	if (m > 0 && x > hiBound2/m) || (m < 0 && x < hiBound2/m) {
		return 0, false
	}
	if (m > 0 && x < loBound2/m) || (m < -1 && x > loBound2/m) {
		return 0, false
	}
	return x * m, true
}

func (c *calculator) div(x value2, d int) (value2, bool) {
	if d == 0 ||
		(x == loBound2 && d == -1) {
		return 0, false
	}
	return x / d, true
}

func (c *calculator) mod(x value2, d int) (value2, bool) {
	if d == 0 || (x == loBound2 && d == -1) {
		return 0, false
	}
	return x % d, true
}

func (c *calculator) alloc(x value2, r, s int) (value2, bool) {
	if r < 0 || s <= 0 {
		return 0, false
	}
	if r > s {
		return 0, false
	}
	if r == 0 {
		return 0, true
	}

	p, ok := c.mul(x, r)
	if !ok {
		return 0, false
	}
	z, ok := c.div(p, s)
	if !ok {
		return 0, false
	}
	return z, true
}

func (c *calculator) neg(x value2) value2 {
	if x > 0 {
		return -x
	}
	return x
}

func (c *calculator) abs(x value2) (value2, bool) {
	if x < 0 {
		if x == loBound2 {
			return 0, false
		}
		return -x, true
	}
	return x, true
}

// pow computes x**e using binary powering algorithm
// for a positive exponent
// see Donald Knuth: The Art of Computer Programming
func (c *calculator) pow(x, e int) (value2, bool) {
	if e < 0 {
		return 0, false
	}
	p := 1
	for e > 0 {
		if e&1 != 0 {
			r, ok := c.mul(p, x)
			if !ok {
				return 0, false
			}
			p = r
		}
		e >>= 1
		r, ok := c.mul(x, x)
		if !ok {
			return 0, false
		}
		x = r
	}
	return p, true
}

func (c *calculator) round(x value2, s int) (value2, bool) {
	if x == 0 {
		return 0, true
	}
	if s < 0 {
		return 0, false
	}
	xabs, ok := c.abs(x)
	if !ok {
		return 0, false
	}

	exp, ok := c.pow(10, s)
	if !ok {
		return 0, false
	}

	m, ok := c.mod(xabs, exp)
	if !ok {
		return 0, false
	}

	if m >= (exp / 2) {
		xabs, ok = c.add(xabs, exp)
		if !ok {
			return 0, false
		}
	}
	q, ok := c.div(xabs, exp)
	if !ok {
		return 0, false
	}
	xabs, ok = c.mul(q, exp)
	if !ok {
		return 0, false
	}
	if x < 0 {
		return -xabs, true
	} else {
		return xabs, true
	}
}

//
//
// To Be Deleted
//
//

func add(x, y *Amount) (*Amount, bool) {
	var r Amount
	neg := x.neg
	if x.neg == y.neg {
		// x + y == x + y
		// (-x) + (-y) == -(x + y)
		val, ok := _add(x.val, y.val)
		if !ok {
			return nil, false
		}
		r.val = val
	} else {
		// x + (-y) == x - y == -(y - x)
		// (-x) + y == y - x == -(x - y)
		if x.cmpByValue(y) >= 0 {
			val, ok := _sub(x.val, y.val)
			if !ok {
				return nil, false
			}
			r.val = val
		} else {
			neg = !neg
			val, ok := _sub(y.val, x.val)
			if !ok {
				return nil, false
			}
			r.val = val
		}
	}
	r.neg = neg
	return &r, true
}

func sub(x, y *Amount) (*Amount, bool) {
	var r Amount
	neg := x.neg
	if x.neg != y.neg {
		// x - (-y) == x + y
		// (-x) - y == -(x + y)
		val, ok := _add(x.val, y.val)
		if !ok {
			return nil, false
		}
		r.val = val
	} else {
		// x - y == x - y == -(y - x)
		// (-x) - (-y) == y - x == -(x - y)
		if x.cmpByValue(y) >= 0 {
			val, ok := _sub(x.val, y.val)
			if !ok {
				return nil, false
			}
			r.val = val
		} else {
			neg = !neg
			val, ok := _sub(y.val, x.val)
			if !ok {
				return nil, false
			}
			r.val = val
		}
	}
	r.neg = neg
	return &r, true
}

// x * m 	== 	2 * 3 	== 	6
// -x * m 	== -2 * 3 	== -6 <
// x * -m 	== 	2 * -3 	== -6 <
// -x * -m 	== -2 * -3 	== 	6
func mul(x *Amount, m int) (*Amount, bool) {
	if x.val == 0 || m == 0 {
		return &Amount{val: 0, neg: false}, true
	}

	neg := x.neg || m < 0
	absm := uint(_absInt(m))
	if !x.neg { // x is positive
		if m > 0 { // x and m is positive
			if x.val > hiBound/absm {
				return nil, false
			}
			neg = false
		} else { // x positive m negative
			if absm < loBound/x.val {
				return nil, false
			}
			neg = true
		}
	} else { // x is negative
		if m > 0 { // x is negative m is positive
			if x.val < loBound/absm {
				return nil, false
			}
			neg = true
		} else { // x and m negative
			if x.val != 0 && absm < hiBound/x.val {
				return nil, false
			}
			neg = false
		}
	}

	val, ok := _mul(x.val, absm, _bound(neg))
	if !ok {
		return nil, false
	}

	return &Amount{val: value(val), neg: neg}, true
}

/*
	quotient:
	 x / y 	== 	4 /	 2 	== 	2
 	-x / y 	== -4 /	 2 	== -2
 	x / -y 	== 	4 /	-2 	== -2
 	-x / -y == -4 /	-2 	== 	2
 	sign determined by -x || -y

*/
func div(a *Amount, d int) (*Amount, bool) {
	// check division by zero and overflow
	if d == 0 || (a.neg && a.val == loBound && d == -1) {
		return nil, false
	}
	neg := a.neg || d < 0
	absd := _absInt(d)
	if _is32() {
		uval := uint32(a.val)
		uabsd := uint32(absd)
		q := uval / uabsd
		return &Amount{val: value(q), neg: neg},
			true
	}
	uval := uint64(a.val)
	uabsd := uint64(absd)
	q := uval / uabsd
	return &Amount{val: value(q), neg: neg}, true
}

/*
remainder:
 	x % y == 5 % 2 == 1
 	x % y == -5 % 2 == -1
 	x % y == 5 % -2 == 1
 	x % y == -5 % -2 == -1
 	sign determined by sign of x
*/
func mod(a *Amount, d int) (amount *Amount, ok bool) {
	if d <= 0 {
		return

	} else {
		if _is32() {
			amount = &Amount{val: value(uint32(a.val) % uint32(_absInt(d))), neg: a.neg}
			ok = true
		} else {
			amount = &Amount{val: value(uint64(a.val) % uint64(_absInt(d))), neg: a.neg}
			ok = true
		}
	}
	return
}

func alloc(a *Amount, r, s int) (*Amount, bool) {
	if r <= 0 {
		return nil, false
	}
	prod, ok := mul(a, r)
	if !ok {
		return nil, false
	}
	z, ok := div(prod, s)
	if !ok {
		return nil, false
	}
	return z, true
}

func neg(a *Amount) *Amount {
	return &Amount{val: a.val, neg: true}
}

func abs(a *Amount) (*Amount, bool) {
	if a.neg && a.val == loBound {
		return nil, false
	}
	return &Amount{val: a.val, neg: false}, true
}

func round(a *Amount, s, i int) *Amount {
	panic("not implemented")
}

/*
	PRIVATE
*/

func _is32() bool {
	return intSize == 32
}

func _bound(neg bool) (bound uint) {
	if neg {
		bound = loBound
	} else {
		bound = hiBound
	}
	return
}

func _add(x, y value) (sum value, ok bool) {
	sum = x + y
	ok = ((x&y)|(x|y)&^sum)>>(intSize-1) == 0
	return
}

func _sub(a, b value) (diff value, ok bool) {
	diff = a - b
	ok = ((^a&b)|(^(a^b)&diff))>>(intSize-1) == 0
	return
}

func _mul(x, y, bound uint) (p value, ok bool) {
	if _is32() {
		p, ok = _mul32(uint32(x), uint32(y), uint32(bound))
	} else {
		p, ok = _mul64(uint64(x), uint64(y), uint64(bound))
	}
	return
}

func _mul32(x, y, bound uint32) (p value, ok bool) {
	tmp := uint64(x) * uint64(y)
	hi, lo := uint32(tmp>>32), uint32(tmp)
	p = value(lo)
	ok = lo <= bound && hi == 0
	return
}

func _mul64(x, y, bound uint64) (p value, ok bool) {
	const mask32 = 1<<32 - 1
	x0 := x & mask32
	x1 := x >> 32
	y0 := y & mask32
	y1 := y >> 32
	w0 := x0 * y0
	t := x1*y0 + w0>>32
	w1 := t & mask32
	w2 := t >> 32
	w1 += x0 * y1
	hi := x1*y1 + w2 + w1>>32
	lo := x * y

	p = value(lo)
	ok = lo <= bound && hi == 0
	return
}

func _absInt(n int) int {
	y := n >> (intSize - 1)
	return (n ^ y) - y
}
