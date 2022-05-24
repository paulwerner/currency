package money

import "testing"

func TestCurrency_Equals(t *testing.T) {
	tcs := []struct {
		c1   Currency
		c2   Currency
		want bool
	}{
		{EUR, EUR, true},
		{USD, EUR, false},
		{EUR, USD, false},
		{USD, USD, true},
	}

	for _, tc := range tcs {
		if eq := tc.c1.Equals(&tc.c2); eq != tc.want {
			t.Errorf("expected %v and %v to be %v, got %v", tc.c1, tc.c2, tc.want, eq)
		}
	}
}
