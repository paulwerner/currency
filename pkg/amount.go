package money

type value = uint

type Amount struct {
	val value
	neg bool
}

func (a *Amount) Int64() int64 {
	return int64(a.val)
}

func amount(v int) *Amount {
	return &Amount{val: value(_abs(v)), neg: v < 0}
}

func (a *Amount) cmpByValue(oa *Amount) int {
	if a.val < oa.val {
		return -1
	}
	if a.val > oa.val {
		return 1
	}
	return 0
}
