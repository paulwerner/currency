package money

import (
	"testing"
)

func TestNew(t *testing.T) {
	tcs := []struct {
		m                    *Money
		expectedAmount       int64
		expectedCurrency     *Currency
		expectedCurrencyCode string
	}{
		{New(1, EUR), 1, MustGetCurrency(EUR), EUR},
		{New(-100, USD), -100, MustGetCurrency(USD), USD},
		{New(1, "XXX"), 1, MustGetCurrency(USD), USD},
	}

	for _, tc := range tcs {
		if tc.m.Amount() != tc.expectedAmount {
			t.Errorf("Expected %d got %d", 1, tc.m.amount)
		}
		if tc.m.Currency() != tc.expectedCurrency {
			t.Errorf("Expected currency %s got %s", EUR, tc.m.currency)
		}
		if tc.m.CurrencyCode() != tc.expectedCurrencyCode {
			t.Errorf("Expected currency %s got %s", EUR, tc.m.currency)
		}
	}
}
