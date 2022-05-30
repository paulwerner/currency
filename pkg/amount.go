package money

type value = uint
type value2 = int

type Amount struct {
	val value
	neg bool
}

func (a *Amount) Equals(oa *Amount) bool {
	return a.val == oa.val && a.neg == oa.neg
}
func amount(v int) *Amount {
	return &Amount{val: value(v), neg: v < 0}
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
