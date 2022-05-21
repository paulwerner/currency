package money

import (
	"testing"
)

func TestMoney_New(t *testing.T) {
	tcs := []struct {
		m                    *Money
		expectedAmount       Amount
		expectedCurrency     *Currency
		expectedCurrencyCode string
	}{
		{New(1, EUR), 1, MustGetCurrency(EUR), EUR},
		{New(-100, USD), -100, MustGetCurrency(USD), USD},
		{New(1, "ABC"), 1, MustGetCurrency(USD), USD},
	}

	for _, tc := range tcs {
		if tc.m.Amount() != tc.expectedAmount {
			t.Errorf("expected amount %d got %d", tc.expectedAmount, tc.m.amount)
		}
		if tc.m.Currency() != tc.expectedCurrency {
			t.Errorf("expected currency %v got %v", tc.expectedCurrency, tc.m.currency)
		}
		if tc.m.CurrencyCode() != tc.expectedCurrencyCode {
			t.Errorf("expected currency %s got %s", tc.expectedCurrencyCode, tc.m.currency.code)
		}
	}
}

func TestMoney_SameCurrency(t *testing.T) {
	tcs := []struct {
		m1   *Money
		m2   *Money
		want bool
	}{
		{New(1, EUR), New(1, EUR), true},
		{New(1, USD), New(1, EUR), false},
		{New(1, EUR), New(1, USD), false},
		{New(1, "ABC"), New(1, USD), true},
	}

	for _, tc := range tcs {
		if result := tc.m1.SameCurrency(tc.m2); result != tc.want {
			t.Errorf("expected same currency to be %v for %s and %s, got %v",
				tc.want, tc.m1.CurrencyCode(), tc.m2.CurrencyCode(), result)
		}
	}
}

func TestMoney_Equals(t *testing.T) {
	tcs := []struct {
		m1      *Money
		m2      *Money
		want    bool
		wantErr error
	}{
		{New(1, USD), New(1, USD), true, nil},
		{New(1, EUR), New(1, EUR), true, nil},
		{New(1, EUR), New(10, EUR), false, nil},
		{New(10, EUR), New(1, EUR), false, nil},
		{New(11, USD), New(1, EUR), false, ErrCurrencyMismatch},
		{New(11, EUR), New(1, USD), false, ErrCurrencyMismatch},
	}

	for _, tc := range tcs {
		ok, err := tc.m1.Equals(tc.m2)
		if ok != tc.want {
			t.Errorf("expected %+v and %+v equality to be %v, got %v", tc.m1, tc.m2, tc.want, ok)
		}
		if err != tc.wantErr {
			t.Errorf("expected %+v and %+v equality error to be %s, got %s", tc.m1, tc.m2, tc.wantErr, err)
		}
	}
}

func TestMoney_GreaterThan(t *testing.T) {
	tcs := []struct {
		m1      *Money
		m2      *Money
		want    bool
		wantErr error
	}{
		{New(1, EUR), New(1, EUR), false, nil},
		{New(1, EUR), New(10, EUR), false, nil},
		{New(2, EUR), New(1, EUR), true, nil},
		{New(0, EUR), New(-1, EUR), true, nil},
		{New(1, EUR), New(-1, EUR), true, nil},
		{New(1, EUR), New(0, EUR), true, nil},
		{New(0, EUR), New(0, EUR), false, nil},

		{New(1, EUR), New(1, USD), false, ErrCurrencyMismatch},
		{New(1, USD), New(1, EUR), false, ErrCurrencyMismatch},
	}

	for _, tc := range tcs {
		ok, err := tc.m1.GreaterThan(tc.m2)
		if ok != tc.want {
			t.Errorf("expected %+v greater than %+v to be %v, got %v", tc.m1, tc.m2, tc.want, ok)
		}
		if err != tc.wantErr {
			t.Errorf("expected %+v greater than %+v error to be %s, got %s", tc.m1, tc.m2, tc.wantErr, err)
		}
	}
}

func TestMoney_GreaterThanOrEqual(t *testing.T) {
	tcs := []struct {
		m1      *Money
		m2      *Money
		want    bool
		wantErr error
	}{
		{New(1, EUR), New(1, EUR), true, nil},
		{New(1, EUR), New(10, EUR), false, nil},
		{New(2, EUR), New(1, EUR), true, nil},
		{New(0, EUR), New(-1, EUR), true, nil},
		{New(1, EUR), New(-1, EUR), true, nil},
		{New(1, EUR), New(0, EUR), true, nil},
		{New(0, EUR), New(0, EUR), true, nil},

		{New(1, EUR), New(1, USD), false, ErrCurrencyMismatch},
		{New(1, USD), New(1, EUR), false, ErrCurrencyMismatch},
	}

	for _, tc := range tcs {
		ok, err := tc.m1.GreaterThanOrEqual(tc.m2)
		if ok != tc.want {
			t.Errorf("expected %+v greater than or equal %+v to be %v, got %v", tc.m1, tc.m2, tc.want, ok)
		}
		if err != tc.wantErr {
			t.Errorf("expected %+v greater than or equal %+v error to be %s, got %s", tc.m1, tc.m2, tc.wantErr, err)
		}
	}
}

func TestMoney_LessThan(t *testing.T) {
	tcs := []struct {
		m1      *Money
		m2      *Money
		want    bool
		wantErr error
	}{
		{New(1, EUR), New(1, EUR), false, nil},
		{New(1, EUR), New(10, EUR), true, nil},
		{New(2, EUR), New(1, EUR), false, nil},
		{New(0, EUR), New(-1, EUR), false, nil},
		{New(1, EUR), New(-1, EUR), false, nil},
		{New(1, EUR), New(0, EUR), false, nil},
		{New(0, EUR), New(0, EUR), false, nil},
		{New(-1, EUR), New(0, EUR), true, nil},
		{New(-1, EUR), New(1, EUR), true, nil},

		{New(1, EUR), New(1, USD), false, ErrCurrencyMismatch},
		{New(1, USD), New(1, EUR), false, ErrCurrencyMismatch},
	}

	for _, tc := range tcs {
		ok, err := tc.m1.LessThan(tc.m2)
		if ok != tc.want {
			t.Errorf("expected %+v less than %+v to be %v, got %v", tc.m1, tc.m2, tc.want, ok)
		}
		if err != tc.wantErr {
			t.Errorf("expected %+v less than %+v error to be %s, got %s", tc.m1, tc.m2, tc.wantErr, err)
		}
	}
}

func TestMoney_LessThanOrEqual(t *testing.T) {
	tcs := []struct {
		m1      *Money
		m2      *Money
		want    bool
		wantErr error
	}{
		{New(1, EUR), New(1, EUR), true, nil},
		{New(1, EUR), New(10, EUR), true, nil},
		{New(2, EUR), New(1, EUR), false, nil},
		{New(0, EUR), New(-1, EUR), false, nil},
		{New(1, EUR), New(-1, EUR), false, nil},
		{New(1, EUR), New(0, EUR), false, nil},
		{New(0, EUR), New(0, EUR), true, nil},
		{New(-1, EUR), New(0, EUR), true, nil},
		{New(-1, EUR), New(1, EUR), true, nil},

		{New(1, EUR), New(1, USD), false, ErrCurrencyMismatch},
		{New(1, USD), New(1, EUR), false, ErrCurrencyMismatch},
	}

	for _, tc := range tcs {
		ok, err := tc.m1.LessThanOrEqual(tc.m2)
		if ok != tc.want {
			t.Errorf("expected %+v less than or equal %+v to be %v, got %v", tc.m1, tc.m2, tc.want, ok)
		}
		if err != tc.wantErr {
			t.Errorf("expected %+v less than or equal %+v error to be %s, got %s", tc.m1, tc.m2, tc.wantErr, err)
		}
	}
}

func TestMoney_IsZero(t *testing.T) {
	tcs := []struct {
		m    *Money
		want bool
	}{
		{New(1, EUR), false},
		{New(0, EUR), true},
		{New(-1, EUR), false},
	}

	for _, tc := range tcs {
		if tc.m.IsZero() != tc.want {
			t.Errorf("expected money value %v to be zero", tc.m.amount)
		}
	}
}

func TestMoney_IsPositive(t *testing.T) {
	tcs := []struct {
		m    *Money
		want bool
	}{
		{New(1, EUR), true},
		{New(0, EUR), false},
		{New(-1, EUR), false},
	}

	for _, tc := range tcs {
		if tc.m.IsPositive() != tc.want {
			t.Errorf("expected money value %v to be positive", tc.m.amount)
		}
	}
}

func TestMoney_IsNegative(t *testing.T) {
	tcs := []struct {
		m    *Money
		want bool
	}{
		{New(1, EUR), false},
		{New(0, EUR), false},
		{New(-1, EUR), true},
	}

	for _, tc := range tcs {
		if tc.m.IsNegative() != tc.want {
			t.Errorf("expected money value %v to be negative", tc.m.amount)
		}
	}
}

func TestMoney_Abs(t *testing.T) {
	tcs := []struct {
		m       *Money
		want    *Money
		wantErr error
	}{
		{New(1, EUR), New(1, EUR), nil},
		{New(0, EUR), New(0, EUR), nil},
		{New(-1, EUR), New(1, EUR), nil},
		{New(-1, EUR), New(1, USD), ErrCurrencyMismatch},
	}

	for _, tc := range tcs {
		abs := tc.m.Abs()
		eq, err := tc.want.Equals(abs)
		if err != tc.wantErr {
			t.Errorf("expected error to be %s, got %s", tc.wantErr, err)
		}
		if err == nil && !eq {
			t.Errorf("expected absolute value of %v to be %v, got %v", tc.m, tc.wantErr, abs)
		}
	}
}

func TestMoney_Neg(t *testing.T) {
	tcs := []struct {
		m       *Money
		want    *Money
		wantErr error
	}{
		{New(1, EUR), New(-1, EUR), nil},
		{New(0, EUR), New(0, EUR), nil},
		{New(-1, EUR), New(-1, USD), ErrCurrencyMismatch},
	}

	for _, tc := range tcs {
		abs := tc.m.Neg()
		eq, err := tc.want.Equals(abs)
		if err != tc.wantErr {
			t.Errorf("expected error to be %s, got %s", tc.wantErr, err)
		}
		if err == nil && !eq {
			t.Errorf("expected negative value of %v to be %v, got %v", tc.m, tc.wantErr, abs)
		}
	}
}

func TestMoney_Add(t *testing.T) {
	tcs := []struct {
		m    *Money
		om   *Money
		want *Money
	}{
		{New(1, EUR), New(-1, EUR), New(0, EUR)},
		{New(-1, EUR), New(1, EUR), New(0, EUR)},
		{New(0, EUR), New(0, EUR), New(0, EUR)},
		{New(-10, EUR), New(-10, EUR), New(-20, EUR)},
		{New(10, EUR), New(10, EUR), New(20, EUR)},
	}

	for _, tc := range tcs {
		sum, err := tc.m.Add(tc.om)
		if err != nil {
			t.Errorf("expected error to be nil, got %s", err)
		}
		eq, err := sum.Equals(tc.want)
		if err != nil {
			t.Errorf("expected error to be nil, got %s", err)
		}
		if !eq {
			t.Errorf("expected %v + %v to be %v, got %v", tc.m, tc.om, tc.want, sum)
		}
	}
}

func TestMoney_Sub(t *testing.T) {
	tcs := []struct {
		m    *Money
		om   *Money
		want *Money
	}{
		{New(1, EUR), New(-1, EUR), New(2, EUR)},
		{New(-1, EUR), New(1, EUR), New(-2, EUR)},
		{New(0, EUR), New(0, EUR), New(0, EUR)},
		{New(-10, EUR), New(-10, EUR), New(0, EUR)},
		{New(20, EUR), New(10, EUR), New(10, EUR)},
		{New(10, EUR), New(20, EUR), New(-10, EUR)},
	}

	for _, tc := range tcs {
		sum, err := tc.m.Sub(tc.om)
		if err != nil {
			t.Errorf("expected error to be nil, got %s", err)
		}
		eq, err := sum.Equals(tc.want)
		if err != nil {
			t.Errorf("expected error to be nil, got %s", err)
		}
		if !eq {
			t.Errorf("expected %v - %v to be %v, got %v", tc.m, tc.om, tc.want, sum)
		}
	}
}

func TestMoney_Multi(t *testing.T) {
	tcs := []struct {
		m    *Money
		mul  int64
		want *Money
	}{
		{New(1, EUR), 1, New(1, EUR)},
		{New(-1, EUR), 2, New(-2, EUR)},
		{New(0, EUR), 1, New(0, EUR)},
		{New(0, EUR), 10, New(0, EUR)},
		{New(-10, EUR), 2, New(-20, EUR)},
		{New(20, EUR), 5, New(100, EUR)},
		{New(10, EUR), -1, New(-10, EUR)},
	}

	for _, tc := range tcs {
		sum := tc.m.Multi(tc.mul)
		eq, err := sum.Equals(tc.want)
		if err != nil {
			t.Errorf("expected error to be nil, got %s", err)
		}
		if !eq {
			t.Errorf("expected %v * %v to be %v, got %v", tc.m, tc.mul, tc.want, sum)
		}
	}
}

func TestMoney_Round(t *testing.T) {
	tcs := []struct {
		m    *Money
		want *Money
	}{
		{New(1024, EUR), New(1000, EUR)},
		{New(10, EUR), New(0, EUR)},
		{New(99, EUR), New(100, EUR)},
		{New(51, EUR), New(100, EUR)},
		{New(49, EUR), New(0, EUR)},
		{New(100, EUR), New(100, EUR)},
		{New(2345, EUR), New(2300, EUR)},
	}

	for _, tc := range tcs {
		round := tc.m.Round()
		eq, err := round.Equals(tc.want)
		if err != nil {
			t.Errorf("expected error to be nil, got %s", err)
		}
		if !eq {
			t.Errorf("expected rounded value of %v to be %v, got %v", tc.m, tc.want, round)
		}
	}
}

func TestMoney_Split(t *testing.T) {
	tcs := []struct {
		m    *Money
		want []*Money
	}{
		{New(1024, EUR), []*Money{New(1024, EUR)}},
		{New(1024, EUR), []*Money{New(512, EUR), New(512, EUR)}},
		{New(1025, EUR), []*Money{New(513, EUR), New(512, EUR)}},
		{New(1024, EUR), []*Money{New(342, EUR), New(341, EUR), New(341, EUR)}},
		{New(1025, EUR), []*Money{New(342, EUR), New(342, EUR), New(341, EUR)}},
	}

	for _, tc := range tcs {
		l := len(tc.want)
		ps, err := tc.m.Split(l)
		if err != nil {
			t.Errorf("expected error to be nil, got %s", err)
		}
		psl := len(ps)
		if psl != l {
			t.Errorf("expected parties count to be %d, got %v", l, len(ps))
		}
		for i := 0; i < l; i++ {
			eq, err := ps[i].Equals(tc.want[i])
			if err != nil {
				t.Errorf("expected error to be nil, got %s", err)
			}
			if !eq {
				t.Errorf("expected %d. party's value to be %v, got %v", i, tc.want[i], ps[i])
			}
		}
	}
}

func TestMoney_SplitWitReminder(t *testing.T) {
	tcs := []struct {
		m     *Money
		want  []*Money
		wantR *Money
	}{
		{New(1024, EUR), []*Money{New(1024, EUR)}, New(0, EUR)},
		{New(1024, EUR), []*Money{New(512, EUR), New(512, EUR)}, New(0, EUR)},
		{New(1025, EUR), []*Money{New(512, EUR), New(512, EUR)}, New(1, EUR)},
		{New(1024, EUR), []*Money{New(341, EUR), New(341, EUR), New(341, EUR)}, New(1, EUR)},
		{New(1025, EUR), []*Money{New(341, EUR), New(341, EUR), New(341, EUR)}, New(2, EUR)},
	}

	for _, tc := range tcs {
		l := len(tc.want)
		ps, r, err := tc.m.SplitWithReminder(l)
		if err != nil {
			t.Errorf("expected error to be nil, got %s", err)
		}
		psl := len(ps)
		if psl != l {
			t.Errorf("expected parties count to be %d, got %v", l, len(ps))
		}
		for i := 0; i < l; i++ {
			eq, err := ps[i].Equals(tc.want[i])
			if err != nil {
				t.Errorf("expected error to be nil, got %s", err)
			}
			if !eq {
				t.Errorf("expected %d. party's value to be %v, got %v", i, tc.want[i], ps[i])
			}
		}
		eq, err := r.Equals(tc.wantR)
		if err != nil {
			t.Errorf("expected error to be nil, got %s", err)
		}
		if !eq {
			t.Errorf("expected reminder value to be %v, got %v", tc.wantR, r)
		}
	}
}

// TODO
// func TestMoney_Alloc(t *testing.T) {
// func TestMoney_AllocWitReminder(t *testing.T) {
