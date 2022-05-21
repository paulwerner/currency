package money

import (
	"testing"
)

func TestSerializer_UnmarshalJSON(t *testing.T) {
	tcs := []struct {
		given string
		want  *Money
	}{
		{`{"amount": 0, "currency": "EUR"}`, New(0, EUR)},
		{`{"amount": 1, "currency": "EUR"}`, New(1, EUR)},
		{`{"amount": 1, "currency": "USD"}`, New(1, USD)},
		{`{"amount": -1, "currency": "USD"}`, New(-1, USD)},
		{`{"amount": -100, "currency": "USD"}`, New(-100, USD)},
	}

	for _, tc := range tcs {
		var m Money
		err := defaultUnmarshalJSON(&m, []byte(tc.given))
		if err != nil {
			t.Errorf("expected unmarshal error to be nil, got %s", err)
		}
		if ok, _ := m.Equals(tc.want); !ok {
			t.Errorf("expected unmarshal value %v to be equal to %v", tc.want, m)
		}
	}
}

func TestSerializer_MarshalJSON(t *testing.T) {
	tcs := []struct {
		given *Money
		want  string
	}{
		{New(0, EUR), `{"amount": 0, "currency": "EUR"}`},
		{New(1, EUR), `{"amount": 1, "currency": "EUR"}`},
		{New(1, USD), `{"amount": 1, "currency": "USD"}`},
		{New(-1, USD), `{"amount": -1, "currency": "USD"}`},
		{New(-100, USD), `{"amount": -100, "currency": "USD"}`},
	}

	for _, tc := range tcs {
		b, err := defaultMarshalJSON(*tc.given)
		if err != nil {
			t.Errorf("expected marshal error to be nil, got %s", err)
		}
		if string(b) != tc.want {
			t.Errorf("expected marshal value to be %s, got %s", tc.want, b)

		}
	}
}
