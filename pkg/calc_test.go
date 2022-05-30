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
	for i, tc := range map[int]struct {
		x      value2
		m      int
		wantOk bool
		wantZ  value2
	}{
		// base cases
		0: {0, 0, true, 0},
		1: {0, 1, true, 0},
		2: {1, 0, true, 0},
		3: {1, 1, true, 1},

		4: {-1, 0, true, 0},
		5: {0, -1, true, 0},
		6: {-1, -1, true, 1},

		// boundaries
		// no overflow
		7: {math.MaxInt, 1, true, math.MaxInt},
		8: {math.MaxInt, -1, true, math.MinInt + 1},
		9: {(math.MaxInt / 2), 2, true, math.MaxInt - 1},

		10: {math.MinInt, 1, true, math.MinInt},
		11: {(math.MinInt / 2), 2, true, math.MinInt},
		12: {(math.MinInt / 2) + 1, 2, true, math.MinInt + 2},

		// overflow
		13: {math.MinInt, -1, false, 0},
		14: {math.MinInt, 2, false, 0},
		15: {(math.MinInt / 2) - 1, 2, false, 0},

		16: {math.MaxInt, 2, false, 0},
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
	for i, tc := range map[int]struct {
		x      value2
		d      int
		wantOk bool
		wantZ  value2
	}{
		// base cases
		0: {0, 1, true, 0},
		1: {1, 1, true, 1},
		2: {1, -1, true, -1},
		3: {-1, 1, true, -1},
		4: {-1, -1, true, 1},

		// division by zero
		5: {0, 0, false, 0},
		6: {1, 0, false, 0},
		7: {-1, 0, false, 0},

		// x < d
		8:  {1, 2, true, 0},
		9:  {-1, 2, true, 0},
		10: {1, -2, true, 0},
		11: {-1, -2, true, 0},
		12: {-1, math.MinInt, true, 0},
		13: {-1, math.MaxInt, true, 0},

		// boundaries
		// no overflow
		14: {math.MaxInt, 1, true, math.MaxInt},
		15: {math.MaxInt, -1, true, -math.MaxInt},
		16: {math.MaxInt, 2, true, math.MaxInt / 2},
		17: {math.MaxInt, -2, true, -(math.MaxInt / 2)},

		18: {math.MinInt, 1, true, math.MinInt},
		19: {math.MinInt + 1, -1, true, -(math.MinInt + 1)},
		20: {math.MinInt, 2, true, math.MinInt / 2},
		21: {math.MinInt, -2, true, math.MaxInt/2 + 1},

		// overflow
		22: {math.MinInt, -1, false, 0},
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
	for i, tc := range map[int]struct {
		x      value2
		d      int
		wantOk bool
		wantZ  value2
	}{
		// base cases
		0: {0, 1, true, 0},
		1: {1, 1, true, 0},
		2: {1, -1, true, 0},
		3: {-1, 1, true, 0},
		4: {-1, -1, true, 0},

		// division by zero
		5: {0, 0, false, 0},
		6: {1, 0, false, 0},
		7: {-1, 0, false, 0},

		8:  {1, 2, true, 1},
		9:  {-1, 2, true, -1},
		10: {1, -2, true, 1},
		11: {-1, -2, true, -1},
		12: {-1, math.MinInt, true, -1},
		13: {-1, math.MaxInt, true, -1},

		// boundaries
		// no overflow
		14: {math.MaxInt, 1, true, 0},
		15: {math.MaxInt, -1, true, 0},
		16: {math.MaxInt, 2, true, 1},
		17: {math.MaxInt, -2, true, 1},

		18: {math.MinInt, 1, true, 0},
		19: {math.MinInt + 1, -1, true, 0},
		20: {math.MinInt, 2, true, 0},
		21: {math.MinInt, -2, true, 0},

		// overflow
		22: {math.MinInt, -1, false, 0},
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
	for i, tc := range map[int]struct {
		x      value2
		r      int
		s      int
		wantOk bool
		wantZ  value2
	}{
		// error cases
		0: {1, -1, 1, false, 0},
		1: {1, 1, 0, false, 0},
		2: {1, 1, -1, false, 0},
		3: {1, 1, math.MinInt, false, 0},
		4: {1, 2, 1, false, 0},

		// success
		5: {10, 5, 10, true, 5},
		6: {10, 7, 10, true, 7},
		7: {10, 3, 10, true, 3},
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
	for i, tc := range map[int]struct {
		x    value2
		want value2
	}{
		0: {1, -1},
		1: {-1, -1},
		2: {-2, -2},
		3: {2, -2},
		4: {math.MinInt, math.MinInt},
		5: {math.MaxInt, math.MinInt + 1},
	} {
		z := calc.neg(tc.x)
		if z != tc.want {
			t.Errorf("[%v]: want z: %v, got: %v", i, tc.want, z)
		}
	}
}
func TestCalc_absolute(t *testing.T) {
	for i, tc := range map[int]struct {
		x      value2
		want   value2
		wantOk bool
	}{
		0: {1, 1, true},
		1: {-1, 1, true},
		2: {-2, 2, true},
		3: {2, 2, true},
		5: {math.MaxInt, math.MaxInt, true},

		// overflow
		4: {math.MinInt, 0, false},
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
