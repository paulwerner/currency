package money

import (
	"math"
	"testing"
)

func TestCalc_add(t *testing.T) {
	maxVal := value(math.MaxUint64)
	a := &Amount{val: maxVal - 1, neg: true}
	b := &Amount{val: maxVal, neg: false}

	r, ok := calc.add(a, b)
	if !ok {
		t.Errorf("should be ok")
	}
	if r.neg {
		t.Errorf("should positive")
	}
	if r.val != 1 {
		t.Errorf("result should be 1, got: %v", r)
	}
}

func TestCalc_addOverflow(t *testing.T) {
	maxVal := value(math.MaxUint64)
	a := &Amount{val: maxVal, neg: true}
	b := &Amount{val: 1, neg: true}

	r, ok := calc.add(a, b)
	if ok {
		t.Errorf("should not be ok")
	}
	if r != nil {
		t.Errorf("result should be nil")
	}
}

func TestCalc_addUnderflow(t *testing.T) {
	minVal := value(math.MaxUint64)
	a := &Amount{val: minVal, neg: true}
	b := &Amount{val: 1, neg: true}

	r, ok := calc.add(a, b)
	if ok {
		t.Errorf("should not be ok")
	}
	if r != nil {
		t.Errorf("result should be nil")
	}
}
