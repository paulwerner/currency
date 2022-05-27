package money

import (
	"fmt"
	"math"
	"math/bits"
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
		t.Errorf("should be positive")
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

func TestCalc_sub(t *testing.T) {
	maxVal := value(math.MaxUint64)
	a := &Amount{val: 0, neg: false}
	b := &Amount{val: maxVal, neg: false}

	r, ok := calc.sub(a, b)
	if !ok {
		t.Errorf("should be ok")
	}
	if !r.neg {
		t.Errorf("should be negative")
	}
	if r.val != value(math.MaxUint64) {
		t.Errorf("result should be 1, got: %v", r)
	}
}

func TestCalc_subOverflow(t *testing.T) {
	maxVal := value(math.MaxUint64)
	a := &Amount{val: maxVal, neg: false}
	b := &Amount{val: 1, neg: true}

	r, ok := calc.sub(a, b)
	if ok {
		t.Errorf("should not be ok")
	}
	if r != nil {
		t.Errorf("result should be nil")
	}
}

func TestCalc_subUnderflow(t *testing.T) {
	maxVal := value(math.MaxUint64)
	a := &Amount{val: maxVal, neg: true}
	b := &Amount{val: 1, neg: false}

	r, ok := calc.sub(a, b)
	if ok {
		t.Errorf("should not be ok")
	}
	if r != nil {
		t.Errorf("result should be nil")
	}
}

func TestFoo(t *testing.T) {
	fmt.Println(max)
	fmt.Println(max / 2)
	hi, lo := bits.Mul64(max/2+1, 5)
	fmt.Println(hi)
	fmt.Println(lo)
	fmt.Println(hi == 0 && lo < max)
	fmt.Println((int64(lo)))

}
