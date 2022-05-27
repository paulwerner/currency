package money

type Amount struct {
	val uint64
	neg bool
}

func amount(v int64) *Amount {
	return &Amount{val: uint64(v), neg: v < 0}
}
