package money

import (
	"math"
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

func TestCalc_boundaries(t *testing.T) {
	if loBound2 != math.MinInt {
		t.Errorf("loBound %v is not as expected: %v", loBound2, math.MinInt)
	}
	if hiBound2 != math.MaxInt {
		t.Errorf("hiBound %v is not as expected: %v", hiBound2, math.MinInt)
	}
}

func TestCalc_addition(t *testing.T) {
	for i, tc := range []struct {
		x      value2
		y      value2
		wantOk bool
		wantZ  value2
	}{
		// base cases
		{0, 0, true, 0},
		{0, 1, true, 1},
		{1, 0, true, 1},
		{1, 1, true, 2},

		{-1, 0, true, -1},
		{0, -1, true, -1},
		{-1, -1, true, -2},

		// boundaries
		// no overflow
		{math.MaxInt, 0, true, math.MaxInt},
		{math.MaxInt - 1, 1, true, math.MaxInt},

		{math.MinInt, 0, true, math.MinInt},
		{math.MinInt + 1, -1, true, math.MinInt},

		// overflow
		{math.MaxInt, 1, false, 0},
		{math.MinInt, -1, false, 0},
	} {
		z, ok := calc.add(tc.x, tc.y)
		if ok != tc.wantOk {
			t.Errorf("[%v]: want ok: %v, got: %v", i, tc.wantOk, ok)
		}
		if z != tc.wantZ {
			t.Errorf("[%v]: want z: %v, got: %v", i, tc.wantZ, z)
		}
	}
}

func TestCalc_subtraction(t *testing.T) {
	for i, tc := range []struct {
		x      value2
		y      value2
		wantOk bool
		wantZ  value2
	}{
		// base cases
		{0, 0, true, 0},
		{0, 1, true, -1},
		{1, 0, true, 1},
		{1, 1, true, 0},

		{-1, 0, true, -1},
		{0, -1, true, 1},
		{-1, -1, true, 0},

		// boundaries
		// no overflow
		{math.MaxInt, 1, true, math.MaxInt - 1},
		{math.MaxInt - 1, -1, true, math.MaxInt},

		{math.MinInt, -1, true, math.MinInt + 1},
		{math.MinInt + 1, 1, true, math.MinInt},

		// overflow
		{math.MaxInt, -1, false, 0},
		{math.MinInt, 1, false, 0},
	} {
		z, ok := calc.sub(tc.x, tc.y)
		if ok != tc.wantOk {
			t.Errorf("[%v]: want ok: %v, got: %v", i, tc.wantOk, ok)
		}
		if z != tc.wantZ {
			t.Errorf("[%v]: want z: %v, got: %v", i, tc.wantZ, z)
		}
	}
}

func TestCalc_multiplication(t *testing.T) {
	for i, tc := range []struct {
		x      value2
		m      int
		wantOk bool
		wantZ  value2
	}{
		// base cases
		{0, 0, true, 0},
		{0, 1, true, 0},
		{1, 0, true, 0},
		{1, 1, true, 1},

		{-1, 0, true, 0},
		{0, -1, true, 0},
		{-1, -1, true, 1},

		// boundaries
		// no overflow
		{math.MaxInt, 1, true, math.MaxInt},
		{math.MaxInt, -1, true, math.MinInt + 1},
		{(math.MaxInt / 2), 2, true, math.MaxInt - 1},

		{math.MinInt, 1, true, math.MinInt},
		{(math.MinInt / 2), 2, true, math.MinInt},
		{(math.MinInt / 2) + 1, 2, true, math.MinInt + 2},

		// overflow
		{math.MinInt, -1, false, 0},
		{math.MinInt, 2, false, 0},
		{(math.MinInt / 2) - 1, 2, false, 0},

		{math.MaxInt, 2, false, 0},
	} {
		z, ok := calc.mul(tc.x, tc.m)
		if ok != tc.wantOk {
			t.Errorf("[%v]: want ok: %v, got: %v", i, tc.wantOk, ok)
		}
		if z != tc.wantZ {
			t.Errorf("[%v]: want z: %v, got: %v", i, tc.wantZ, z)
		}
	}
}

func TestCalc_division(t *testing.T) {
	for i, tc := range []struct {
		x      value2
		d      int
		wantOk bool
		wantZ  value2
	}{
		// base cases
		{0, 1, true, 0},
		{1, 1, true, 1},
		{1, -1, true, -1},
		{-1, 1, true, -1},
		{-1, -1, true, 1},

		// division by zero
		{0, 0, false, 0},
		{1, 0, false, 0},
		{-1, 0, false, 0},

		// x < d
		{1, 2, true, 0},
		{-1, 2, true, 0},
		{1, -2, true, 0},
		{-1, -2, true, 0},
		{-1, math.MinInt, true, 0},
		{-1, math.MaxInt, true, 0},

		// boundaries
		// no overflow
		{math.MaxInt, 1, true, math.MaxInt},
		{math.MaxInt, -1, true, -math.MaxInt},
		{math.MaxInt, 2, true, math.MaxInt / 2},
		{math.MaxInt, -2, true, -(math.MaxInt / 2)},

		{math.MinInt, 1, true, math.MinInt},
		{math.MinInt + 1, -1, true, -(math.MinInt + 1)},
		{math.MinInt, 2, true, math.MinInt / 2},
		{math.MinInt, -2, true, math.MaxInt/2 + 1},

		// overflow
		{math.MinInt, -1, false, 0},
	} {
		z, ok := calc.div(tc.x, tc.d)
		if ok != tc.wantOk {
			t.Errorf("[%v]: want ok: %v, got: %v", i, tc.wantOk, ok)
		}
		if z != tc.wantZ {
			t.Errorf("[%v]: want z: %v, got: %v", i, tc.wantZ, z)
		}
	}
}

func TestCalc_modulo(t *testing.T) {
	for i, tc := range []struct {
		x      value2
		d      int
		wantOk bool
		wantZ  value2
	}{
		// base cases
		{0, 1, true, 0},
		{1, 1, true, 0},
		{1, -1, true, 0},
		{-1, 1, true, 0},
		{-1, -1, true, 0},

		// division by zero
		{0, 0, false, 0},
		{1, 0, false, 0},
		{-1, 0, false, 0},

		{1, 2, true, 1},
		{-1, 2, true, -1},
		{1, -2, true, 1},
		{-1, -2, true, -1},
		{-1, math.MinInt, true, -1},
		{-1, math.MaxInt, true, -1},

		// boundaries
		// no overflow
		{math.MaxInt, 1, true, 0},
		{math.MaxInt, -1, true, 0},
		{math.MaxInt, 2, true, 1},
		{math.MaxInt, -2, true, 1},

		{math.MinInt, 1, true, 0},
		{math.MinInt + 1, -1, true, 0},
		{math.MinInt, 2, true, 0},
		{math.MinInt, -2, true, 0},

		// overflow
		{math.MinInt, -1, false, 0},
	} {
		z, ok := calc.mod(tc.x, tc.d)
		if ok != tc.wantOk {
			t.Errorf("[%v]: want ok: %v, got: %v", i, tc.wantOk, ok)
		}
		if z != tc.wantZ {
			t.Errorf("[%v]: want z: %v, got: %v", i, tc.wantZ, z)
		}
	}
}

func TestCalc_allocation(t *testing.T) {
	for i, tc := range []struct {
		x      value2
		r      int
		s      int
		wantOk bool
		wantZ  value2
	}{
		// error cases
		{1, -1, 1, false, 0},
		{1, 1, 0, false, 0},
		{1, 1, -1, false, 0},
		{1, 1, math.MinInt, false, 0},
		{1, 2, 1, false, 0},

		// success
		{10, 5, 10, true, 5},
		{10, 7, 10, true, 7},
		{10, 3, 10, true, 3},
	} {
		z, ok := calc.alloc(tc.x, tc.r, tc.s)
		if ok != tc.wantOk {
			t.Errorf("[%v]: want ok: %v, got: %v", i, tc.wantOk, ok)
		}
		if z != tc.wantZ {
			t.Errorf("[%v]: want z: %v, got: %v", i, tc.wantZ, z)
		}
	}
}

func TestCalc_negation(t *testing.T) {
	for i, tc := range []struct {
		x    value2
		want value2
	}{
		{1, -1},
		{-1, -1},
		{-2, -2},
		{2, -2},
		{math.MinInt, math.MinInt},
		{math.MaxInt, math.MinInt + 1},
	} {
		z := calc.neg(tc.x)
		if z != tc.want {
			t.Errorf("[%v]: want z: %v, got: %v", i, tc.want, z)
		}
	}
}
func TestCalc_absolute(t *testing.T) {
	for i, tc := range []struct {
		x      value2
		want   value2
		wantOk bool
	}{
		{1, 1, true},
		{-1, 1, true},
		{-2, 2, true},
		{2, 2, true},
		{math.MaxInt, math.MaxInt, true},

		// overflow
		{math.MinInt, 0, false},
	} {
		z, ok := calc.abs(tc.x)
		if ok != tc.wantOk {
			t.Errorf("[%v]: want ok: %v, got: %v", i, tc.wantOk, ok)
		}
		if z != tc.want {
			t.Errorf("[%v]: want z: %v, got: %v", i, tc.want, z)
		}
	}
}

func TestCalc_power(t *testing.T) {
	for i, tc := range []struct {
		x      value2
		e      int
		want   value2
		wantOk bool
	}{
		// base cases
		{0, 0, 1, true},
		{1, 0, 1, true},
		{1, 1, 1, true},
		{0, 1, 0, true},

		// higher cases
		{2, 2, 4, true},
		{2, 4, 16, true},
		{3, 2, 9, true},
		{3, 3, 27, true},

		// negative values
		{-1, 0, 1, true},
		{-1, 1, -1, true},
		{-1, 2, 1, true},
		{-2, 2, 4, true},
		{-2, 3, -8, true},
		{-2, 4, 16, true},

		// negative exponent
		{-1, -1, 0, false},
		{0, -1, 0, false},
		{2, math.MinInt, 0, false},

		// overflow
		{math.MaxInt, 2, 0, false},
		{2, math.MaxInt, 0, false},
		{math.MinInt, 2, 0, false},
	} {
		z, ok := calc.pow(tc.x, tc.e)
		if ok != tc.wantOk {
			t.Errorf("[%v]: want ok: %v, got: %v", i, tc.wantOk, ok)
		}
		if z != tc.want {
			t.Errorf("[%v]: want z: %v, got: %v", i, tc.want, z)
		}
	}
}

func TestCalc_rounding(t *testing.T) {
	for i, tc := range []struct {
		x      value2
		s      int
		want   value2
		wantOk bool
	}{
		// positive values
		{110, 2, 100, true},
		{111, 2, 100, true},
		{111, 2, 100, true},
		{301, 2, 300, true},
		{401, 2, 400, true},

		// edge cases
		{449, 2, 400, true},
		{450, 2, 500, true},
		{451, 2, 500, true},
		{499, 2, 500, true},
		{501, 2, 500, true},

		{111, 3, 0, true},
		{1111, 3, 1000, true},
		{1499, 3, 1000, true},
		{1500, 3, 2000, true},

		// negative values
		{-110, 2, -100, true},
		{-111, 2, -100, true},
		{-111, 2, -100, true},
		{-301, 2, -300, true},
		{-401, 2, -400, true},

		// edge cases
		{-449, 2, -400, true},
		{-450, 2, -500, true},
		{-451, 2, -500, true},
		{-499, 2, -500, true},
		{-501, 2, -500, true},

		{-111, 3, 0, true},
		{-1111, 3, -1000, true},
		{-1499, 3, -1000, true},
		{-1500, 3, -2000, true},

		// negative scale
		{1, -2, 0, false},
	} {
		z, ok := calc.round(tc.x, tc.s)
		if ok != tc.wantOk {
			t.Errorf("[%v]: want ok: %v, got: %v", i, tc.wantOk, ok)
		}
		if z != tc.want {
			t.Errorf("[%v]: want z: %v, got: %v", i, tc.want, z)
		}
	}
}
