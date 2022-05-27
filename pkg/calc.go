package money

type calculator struct{}

var calc = calculator{}

const (
	min = -1 << 63
	max = 1<<63 - 1
)

func (c *calculator) add(x, y *Amount) (*Amount, bool) {
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

func (c *calculator) sub(x, y *Amount) (*Amount, bool) {
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

func (c *calculator) mul(x *Amount, m int) (*Amount, bool) {
	panic("not implemented")
}

func (c *calculator) div(a *Amount, d int) (*Amount, bool) {
	panic("not implemented")
}

func (c *calculator) mod(a *Amount, d int) *Amount {
	panic("not implemented")
}

func (c *calculator) alloc(a *Amount, r, s int) (*Amount, bool) {
	panic("not implemented")
}

func (c *calculator) neg(a *Amount) *Amount {
	panic("not implemented")
}

func (c *calculator) abs(a *Amount) *Amount {
	panic("not implemented")
}

func (c *calculator) round(a *Amount, s, i int) *Amount {
	panic("not implemented")
}
