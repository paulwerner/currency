package money

import (
	"fmt"
	"math"
	"testing"
)

func TestAmount(t *testing.T) {
	a := amount(int64(math.MinInt64))
	if !a.neg {
		t.Errorf("expected %v to be positive", a)
	}
	if int64(a.val) < 0 {
		t.Logf("a.val: %v", a.val)
	}

	b := uint64(0)
	fmt.Printf("%b\n", b)

	c := b | 0xf
	fmt.Printf("%v\n", c)

}
