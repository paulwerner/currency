package money

type value uint64

func (v value) cmp(ov value) int {
	if v < ov {
		return -1
	}
	if v > 0 {
		return 1
	}
	return 0
}

type Amount struct {
	val value
	neg bool
}

func (a *Amount) valueInt64() int64 {
	return int64(a.val)
}

func amount(v int64) *Amount {
	return &Amount{val: value(v), neg: v < 0}
}
