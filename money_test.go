package money

import (
	"testing"
)

func TestNew(t *testing.T) {
	m := New(1, EUR)
	if m.Amount() != 1 {
		t.Errorf("Expected %d got %d", 1, m.amount)
	}
	if m.CurrencyCode() != EUR {
		t.Errorf("Expected currency %s got %s", EUR, m.currency)
	}

	m = New(-100, EUR)
	if m.Amount() != -100 {
		t.Errorf("Expected %d got %d", -100, m.amount)
	}
}
