package money

import "testing"

func TestNewGenericInt(t *testing.T) {
	m := New[Int](1, EUR)
	if m.amount != 1 {
		t.Errorf("Expected %d got %d", 1, m.amount)
	}
	if m.currency.Code != EUR {
		t.Errorf("Expected currency %s got %s", EUR, m.currency)
	}

	m = New[Int](-100, EUR)
	if m.amount != Int(-100) {
		t.Errorf("Expected %d got %d", -100, m.amount)
	}
}
