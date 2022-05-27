/*
maxInt64:        9223372036854775807
maxInt64:        111111111111111111111111111111111111111111111111111111111111111
minInt64:        -9223372036854775808
minInt64:        -1000000000000000000000000000000000000000000000000000000000000000
maxUint64:       9223372036854775807
maxUint64:       111111111111111111111111111111111111111111111111111111111111111
minUint64:       9223372036854775808
minUint64:       1000000000000000000000000000000000000000000000000000000000000000

valUint64:       9223372036854775808
valUint64:       1000000000000000000000000000000000000000000000000000000000000000
*/
package money

const (
	min int64 = -1 << 63
	max int64 = 1<<63 - 1
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
		if x.val.cmp(y.val) >= 0 {
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

func _add(x, y value) (sum value, ok bool) {
	sum = x + y
	ok = ((x&y)|(x|y)&^sum)>>63 == 0
	return
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
		if x.val.cmp(y.val) >= 0 {
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

func _sub(a, b value) (diff value, ok bool) {
	diff = a - b
	ok = ((^a&b)|(^(a^b)&diff))>>63 == 0
	return
}

func mul(x *Amount, m int64) (*Amount, bool) {
	if x.val == 0 || m == 0 {
		return &Amount{val: 0, neg: false}, true
	}

	neg := x.neg
	if !x.neg { // x is positive
		if m > 0 { // x and m is positive
			if x.Int64() > max/m {
				return nil, false
			}
			neg = false
		} else { // x positive m negative
			if m < (min / x.Int64()) {
				return nil, false
			}
			neg = true
		}
	} else { // x is negative
		if m > 0 { // x is negative m is positive
			if x.Int64() < (min / m) {
				return nil, false
			}
			neg = true
		} else { // x and m negative
			if x.val != 0 && m < (max/x.Int64()) {
				return nil, false
			}
			neg = false
		}
	}
	val := _abs(x.Int64() * m)
	return &Amount{val: value(val), neg: neg}, true
}

func div(a *Amount, d int64) (*Amount, bool) {
	panic("not implemented")
}

func mod(a *Amount, d int64) *Amount {
	panic("not implemented")
}

func alloc(a *Amount, r, s int) (*Amount, bool) {
	panic("not implemented")
}

func neg(a *Amount) *Amount {
	panic("not implemented")
}

func abs(a *Amount) *Amount {
	panic("not implemented")
}

func _abs(x int64) int64 {
	if x < 0 {
		return -x
	}
	return x
}

func round(a *Amount, s, i int) *Amount {
	panic("not implemented")
}
