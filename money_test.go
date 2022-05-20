package money

import (
	"testing"
)

func TestMoney_New(t *testing.T) {
	tcs := []struct {
		m                    *Money
		expectedAmount       int64
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
			t.Errorf("expected currency %s got %s", tc.expectedCurrency, tc.m.currency)
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
		ok, err := tc.m1.Equal(tc.m2)
		if ok != tc.want {
			t.Errorf("expected %+v and %+v equality to be %v, got %v ", tc.m1, tc.m2, tc.want, ok)
		}
		if err != tc.wantErr {
			t.Errorf("expected %+v and %+v equality to be %v, got %v ", tc.m1, tc.m2, tc.want, ok)
		}
	}
}
