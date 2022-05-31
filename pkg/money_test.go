package money

import "testing"

func TestMoney_Base(t *testing.T) {
	// Addition
	m, _ := New(10, EUR)
	m2, _ := New(2, EUR)

	sum, err := m.Add(m2)
	if err != nil {
		t.Errorf("err not nil, got %v", err)
	}
	if !(sum.amount == amount(12)) {
		t.Errorf("sum.amount != %v", 12)
	}
	if !sum.currency.Equals(&EUR) {
		t.Errorf("sum.currency != %v", sum.currency)
	}

	// Subtraction
	diff, err := m.Sub(m2)
	if err != nil {
		t.Errorf("err not nil, got %v", err)
	}
	if !(diff.amount == amount(8)) {
		t.Errorf("sum.amount != %v, got %v", 8, sum.amount)
	}
	if !diff.currency.Equals(&EUR) {
		t.Errorf("sum.currency != %v", diff.currency)
	}

	// Multiplication
	mul := 2
	prod, err := m.Mul(mul)
	if err != nil {
		t.Errorf("err not nil, got %v", err)
	}
	if !(prod.amount == amount(20)) {
		t.Errorf("prod.amount != %v, got %v", 8, prod.amount)
	}
	if !prod.currency.Equals(&EUR) {
		t.Errorf("prod.currency != %v", prod.currency)
	}

	// Split
	n := 2
	ps, r, err := m.Split(n)
	if err != nil {
		t.Errorf("err not nil, got %v", err)
	}
	if len(ps) != 2 {
		t.Errorf("ps.len != 2, got %v", 2)
	}
	if r.amount != 0 {
		t.Errorf("ps.reminder != 0, got %v", r.amount)
	}
	for i, p := range ps {
		if p.amount != 5 {
			t.Errorf("ps[%v].amount != 5, got %v", i, p.amount)
		}
		if !p.currency.Equals(&EUR) {
			t.Errorf("ps[%v].currency != EUR, got %v", i, p.currency)
		}
	}
}
