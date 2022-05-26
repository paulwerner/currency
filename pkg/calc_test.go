package money

import "testing"

func TestCalculator_add(t *testing.T) {
	calc := calculator{}
	tcs := []struct {
		a    Value
		b    Value
		want Value
	}{
		{1, 2, 3},
		{-1, 2, 1},
		{-1, -2, -3},
	}

	for i, tc := range tcs {
		if sum, ok := calc.add(tc.a, tc.b); !ok || sum != tc.want {
			t.Errorf("[%v]: expected %d + %d to be %d, got %d",
				i, tc.a, tc.b, tc.want, sum)
		}
	}
}

func TestCalculator_sub(t *testing.T) {
	calc := calculator{}
	tcs := []struct {
		a    Value
		b    Value
		want Value
	}{
		{1, 2, -1},
		{-1, 2, -3},
		{10, -2, 12},
	}

	for i, tc := range tcs {
		if diff, ok := calc.sub(tc.a, tc.b); !ok || diff != tc.want {
			t.Errorf("[%v]: expected %d - %d to be %d, got %d",
				i, tc.a, tc.b, tc.want, diff)
		}
	}
}

func TestCalculator_mul(t *testing.T) {
	calc := calculator{}
	tcs := []struct {
		a    Value
		b    Value
		want Value
	}{
		{1, 2, 2},
		{-1, 2, -2},
		{2, -2, -4},
		{10, -2, -20},
	}

	for i, tc := range tcs {
		if prod, ok := calc.mul(tc.a, tc.b); !ok || prod != tc.want {
			t.Errorf("[%v]: expected %d * %d to be %d, got %d",
				i, tc.a, tc.b, tc.want, prod)
		}
	}
}

func TestCalculator_div(t *testing.T) {
	calc := calculator{}
	tcs := []struct {
		a    Value
		b    Value
		want Value
	}{
		{10, 2, 5},
		{11, 2, 5},
		{-10, 2, -5},
		{-11, 2, -5},
	}

	for i, tc := range tcs {
		if quot, ok := calc.div(tc.a, tc.b); !ok || quot != tc.want {
			t.Errorf("[%v]: expected %d / %d to be %d, got %d",
				i, tc.a, tc.b, tc.want, quot)
		}
	}
}

func TestCalculator_mod(t *testing.T) {
	calc := calculator{}
	tcs := []struct {
		a    Value
		d    Value
		want Value
	}{
		{100, 10, 0},
		{100, 11, 1},
		{-100, 2, 0},
		{-100, -2, 0},
		{-100, -3, -1},
		{100, 80, 20},
		{100, 51, 49},
	}

	for i, tc := range tcs {
		if mod, ok := calc.mod(tc.a, tc.d); !ok || mod != tc.want {
			t.Errorf("[%v]: expected %d mod %d to be %d, got %d",
				i, tc.a, tc.d, tc.want, mod)
		}
	}
}

func TestCalculator_alloc(t *testing.T) {
	calc := calculator{}
	tcs := []struct {
		a    Value
		r    int64
		s    int64
		want Value
	}{
		{100, 10, 50, 20},
		{100, 2, 50, 4},
		{100, 2, 500, 0},
		{100, 80, 500, 16},
		{100, 51, 50, 100},
		{100, 2, 1, 100},
	}

	for i, tc := range tcs {
		if alloc, ok := calc.alloc(tc.a, tc.r, tc.s); !ok || alloc != tc.want {
			t.Errorf("[%v]: expected %d allocated in %d ratios with sum %d to be %d, got %d, ok %v",
				i, tc.a, tc.r, tc.s, tc.want, alloc, ok)
		}
	}
}

func TestCalculator_abs(t *testing.T) {
	calc := calculator{}
	tcs := []struct {
		a    Value
		want Value
	}{
		{100, 100},
		{-100, 100},
	}

	for i, tc := range tcs {
		if abs := calc.abs(tc.a); abs != tc.want {
			t.Errorf("[%v]: expected absolute of %d to be %d, got %d",
				i, tc.a, tc.want, abs)
		}
	}
}

func TestCalculator_neg(t *testing.T) {
	calc := calculator{}
	tcs := []struct {
		a    Value
		want Value
	}{
		{100, -100},
		{-100, -100},
	}

	for i, tc := range tcs {
		if neg := calc.neg(tc.a); neg != tc.want {
			t.Errorf("[%v]: expected negative of %d to be %d, got %d",
				i, tc.a, tc.want, neg)
		}
	}
}

func TestCalculator_round(t *testing.T) {
	calc := calculator{}
	tcs := []struct {
		a    Value
		e    int
		want Value
	}{
		{1023, 2, 1000},
		{1023, 1, 1020},
		{1, 1, 0},
		{1, 0, 1},
	}

	for i, tc := range tcs {
		if round, ok := calc.round(tc.a, tc.e); !ok || round != tc.want {
			t.Errorf("[%v]: expected round of %d with %d decimal places to be %d, got %d",
				i, tc.a, tc.e, tc.want, round)
		}
	}
}
