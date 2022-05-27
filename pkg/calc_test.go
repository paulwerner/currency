package money

import (
	"fmt"
	"math"
	"strconv"
	"testing"
)

func TestCalc_add(t *testing.T) {
	maxVal := value(math.MaxUint64)
	a := &Amount{val: maxVal - 1, neg: true}
	b := &Amount{val: maxVal, neg: false}

	r, ok := add(a, b)
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

	r, ok := add(a, b)
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

	r, ok := add(a, b)
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

	r, ok := sub(a, b)
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

	r, ok := sub(a, b)
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

	r, ok := sub(a, b)
	if ok {
		t.Errorf("should not be ok")
	}
	if r != nil {
		t.Errorf("result should be nil")
	}
}

func TestFoo(t *testing.T) {
	maxInt64 := int64(math.MaxInt64)
	fmt.Printf("maxInt64:\t %v\n", maxInt64)
	fmt.Printf("maxInt64:\t %s\n", strconv.FormatInt(maxInt64, 2))

	minInt64 := int64(math.MinInt64)
	fmt.Printf("minInt64:\t %v\n", minInt64)
	fmt.Printf("minInt64:\t %s\n", strconv.FormatInt(minInt64, 2))

	maxUint64 := uint64(maxInt64)
	fmt.Printf("maxUint64:\t %v\n", maxUint64)
	fmt.Printf("maxUint64:\t %s\n", strconv.FormatUint(maxUint64, 2))

	minUint64 := uint64(minInt64)
	fmt.Printf("minUint64:\t %v\n", minUint64)
	fmt.Printf("minUint64:\t %s\n", strconv.FormatUint(minUint64, 2))

	fmt.Println("")
	valUint64 := uint64(minInt64 >> 1 << 1)
	fmt.Printf("valUint64:\t %v\n", valUint64)
	fmt.Printf("valUint64:\t %s\n", strconv.FormatUint(valUint64, 2))
	fmt.Println("")
}
