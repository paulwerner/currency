package money

import (
	"testing"
)

func TestAmount_NewAmount(t *testing.T) {
	tcs := []struct {
		aVal       Value
		aCur       string
		wantAmount *Amount
	}{
		{1, "EUR", &Amount{val: 1, currency: EUR}},
		{-100, "EUR", &Amount{val: -100, currency: EUR}},
		{-100, "USD", &Amount{val: -100, currency: USD}},
	}

	for _, tc := range tcs {
		a, err := NewAmount(tc.aVal, tc.aCur)
		if err != nil {
			t.Errorf("unexpected error creating amount with value %d and currency %s", tc.aVal, tc.aCur)
		}
		if a.Value() != tc.wantAmount.Value() {
			t.Errorf("expected amount value %d, got %d", tc.wantAmount.Value(), a.Value())
		}
		if a.Currency() != tc.wantAmount.Currency() {
			t.Errorf("expected currency %v, got %v", tc.wantAmount.Currency(), a.Currency())
		}
	}
}

func TestAmount_NewAmountFailsWithUnrecognizedCurrency(t *testing.T) {
	tcs := []struct {
		currency  string
		wantError bool
	}{
		{"ABC", true},
		{"XYZ", true},
		{"EUR", false},
	}

	for _, tc := range tcs {
		_, err := NewAmount(1, tc.currency)
		if err == nil && tc.wantError {
			t.Error("expected error not occurred")
		}
	}
}

func TestAmount_AssertSameCurrency(t *testing.T) {
	tcs := []struct {
		a1   Amount
		a2   Amount
		want error
	}{
		{Amount{1, EUR}, Amount{1, EUR}, nil},
		{Amount{1, EUR}, Amount{1, USD}, errCurrencyMismatch},
		{Amount{1, USD}, Amount{1, EUR}, errCurrencyMismatch},
		{Amount{1, USD}, Amount{1, USD}, nil},
	}

	for _, tc := range tcs {
		if err := tc.a1.assertSameCurrency(&tc.a2); err != tc.want {
			t.Errorf("expected assertion error for %v and %v to be %v, got %v", tc.a1, tc.a2, tc.want, err)
		}
	}
}
