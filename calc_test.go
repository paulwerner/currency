package currency

import (
	"math"
	"testing"
)

func TestCalc_boundaries(t *testing.T) {
	if loBound != math.MinInt {
		t.Errorf("loBound %v is not as expected: %v", loBound, math.MinInt)
	}
	if hiBound != math.MaxInt {
		t.Errorf("hiBound %v is not as expected: %v", hiBound, math.MinInt)
	}
}

func TestCalc_addition(t *testing.T) {
	for i, tc := range []struct {
		x      int
		y      int
		wantOk bool
		wantZ  int
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
		z, ok := add(tc.x, tc.y)
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
		x      int
		y      int
		wantOk bool
		wantZ  int
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
		z, ok := sub(tc.x, tc.y)
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
		x      int
		m      int
		wantOk bool
		wantZ  int
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
		z, ok := mul(tc.x, tc.m)
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
		x      int
		d      int
		wantOk bool
		wantZ  int
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
		z, ok := div(tc.x, tc.d)
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
		x      int
		d      int
		wantOk bool
		wantZ  int
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
		z, ok := mod(tc.x, tc.d)
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
		x      int
		r      int
		s      int
		wantOk bool
		wantZ  int
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
		z, ok := alloc(tc.x, tc.r, tc.s)
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
		x      int
		want   int
		wantOk bool
	}{
		{0, 0, true},
		{1, -1, true},
		{-1, 1, true},
		{-2, 2, true},
		{2, -2, true},
		{math.MaxInt, math.MinInt + 1, true},
		{math.MinInt + 1, math.MaxInt, true},

		// overflow
		{math.MinInt, 0, false},
	} {
		z, ok := neg(tc.x)
		if ok != tc.wantOk {
			t.Errorf("[%v]: want ok: %v, got: %v", i, tc.wantOk, ok)
		}
		if z != tc.want {
			t.Errorf("[%v]: want z: %v, got: %v", i, tc.want, z)
		}
	}
}
func TestCalc_absolute(t *testing.T) {
	for i, tc := range []struct {
		x      int
		want   int
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
		z, ok := abs(tc.x)
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
		x      int
		e      int
		want   int
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
		z, ok := pow(tc.x, tc.e)
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
		x      int
		step   int
		want   int
		wantOk bool
	}{
		// step 1 is the identity
		{0, 1, 0, true},
		{1, 1, 1, true},
		{-1, 1, -1, true},
		{math.MaxInt, 1, math.MaxInt, true},
		{math.MinInt, 1, math.MinInt, true},

		// positive values, step 100
		{110, 100, 100, true},
		{111, 100, 100, true},
		{301, 100, 300, true},
		{401, 100, 400, true},

		// ties round away from zero
		{449, 100, 400, true},
		{450, 100, 500, true},
		{451, 100, 500, true},
		{499, 100, 500, true},
		{501, 100, 500, true},

		{111, 1000, 0, true},
		{1111, 1000, 1000, true},
		{1499, 1000, 1000, true},
		{1500, 1000, 2000, true},

		// negative values
		{-110, 100, -100, true},
		{-111, 100, -100, true},
		{-301, 100, -300, true},
		{-401, 100, -400, true},

		{-449, 100, -400, true},
		{-450, 100, -500, true},
		{-451, 100, -500, true},
		{-499, 100, -500, true},
		{-501, 100, -500, true},

		{-111, 1000, 0, true},
		{-1111, 1000, -1000, true},
		{-1499, 1000, -1000, true},
		{-1500, 1000, -2000, true},

		// cash rounding increments (e.g. 0.05 CHF, 50 øre)
		{102, 5, 100, true},
		{103, 5, 105, true},
		{-102, 5, -100, true},
		{-103, 5, -105, true},
		{1224, 50, 1200, true},
		{1225, 50, 1250, true},

		// invalid steps
		{1, 0, 0, false},
		{1, -2, 0, false},

		// overflow: math.MaxInt ends in ...807, so with step 10 the
		// remainder 7 rounds up past the representable maximum
		{math.MaxInt, 10, 0, false},
		{math.MinInt, 2, 0, false},
	} {
		z, ok := round(tc.x, tc.step)
		if ok != tc.wantOk {
			t.Errorf("[%v]: want ok: %v, got: %v", i, tc.wantOk, ok)
		}
		if z != tc.want {
			t.Errorf("[%v]: want z: %v, got: %v", i, tc.want, z)
		}
	}
}
