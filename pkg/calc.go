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
)

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
	absm := uint(_abs(m))
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
	absd := _abs(d)
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
func mod(a *Amount, d int) *Amount {
	panic("not implemented")
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
	return &Amount{val: a.val, neg: !a.neg}
}

func abs(a *Amount) *Amount {
	panic("not implemented")
}

func _abs(n int) int {
	y := n >> (intSize - 1)
	return (n ^ y) - y
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
