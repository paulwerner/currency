package currency

import (
	"encoding/json"
	"errors"
	"math"
	"testing"
)

func mustAmount(t *testing.T, v int, cur Currency) *Amount {
	t.Helper()
	a, err := NewAmount(v, cur)
	if err != nil {
		t.Fatalf("NewAmount(%v, %v): unexpected error %v", v, &cur, err)
	}
	return a
}

func TestAmount_New(t *testing.T) {
	a := mustAmount(t, 1234, EUR)
	if a.Amount() != 1234 {
		t.Errorf("Amount() = %v, want 1234", a.Amount())
	}
	if !a.Currency().Equals(&EUR) {
		t.Errorf("Currency() = %v, want EUR", a.Currency())
	}
}

func TestAmount_NewFromISO(t *testing.T) {
	a, err := NewFromISO(500, "USD")
	if err != nil {
		t.Fatalf("unexpected error %v", err)
	}
	if a.Amount() != 500 || !a.Currency().Equals(&USD) {
		t.Errorf("got %v, want USD 5.00", a)
	}

	// XXX is valid and yields the canonical undefined currency
	a, err = NewFromISO(1, "XXX")
	if err != nil {
		t.Fatalf("unexpected error %v", err)
	}
	if !a.Currency().Equals(&XXX) {
		t.Errorf("currency = %v, want XXX", a.Currency())
	}

	if _, err = NewFromISO(1, "invalid"); !errors.Is(err, ErrISOCodeMalformed) {
		t.Errorf("err = %v, want ErrISOCodeMalformed", err)
	}
	if _, err = NewFromISO(1, "BTC"); !errors.Is(err, ErrISOCodeNotRecognized) {
		t.Errorf("err = %v, want ErrISOCodeNotRecognized", err)
	}
}

func TestAmount_ZeroValue(t *testing.T) {
	var a Amount
	if !a.IsZero() {
		t.Error("zero Amount should be zero")
	}
	if got := a.Currency().Code(); got != "XXX" {
		t.Errorf("zero Amount currency = %q, want XXX", got)
	}
	if got := a.String(); got != "XXX 0.00" {
		t.Errorf("zero Amount String() = %q, want \"XXX 0.00\"", got)
	}
}

func TestAmount_Arithmetic(t *testing.T) {
	m := mustAmount(t, 10, EUR)

	// Addition
	m2 := mustAmount(t, 2, EUR)

	sum, err := m.Add(m2)
	if err != nil {
		t.Errorf("err not nil, got %v", err)
	}
	if sum.Amount() != 12 {
		t.Errorf("sum.value != %v", 12)
	}
	if !sum.Currency().Equals(&EUR) {
		t.Errorf("sum.currency != %v", sum.Currency())
	}

	// Subtraction
	diff, err := m.Sub(m2)
	if err != nil {
		t.Errorf("err not nil, got %v", err)
	}
	if diff.Amount() != 8 {
		t.Errorf("diff.value != %v, got %v", 8, diff.Amount())
	}
	if !diff.Currency().Equals(&EUR) {
		t.Errorf("diff.currency != %v", diff.Currency())
	}

	// Multiplication
	prod, err := m.Mul(2)
	if err != nil {
		t.Errorf("err not nil, got %v", err)
	}
	if prod.Amount() != 20 {
		t.Errorf("prod.value != %v, got %v", 20, prod.Amount())
	}
	if !prod.Currency().Equals(&EUR) {
		t.Errorf("prod.currency != %v", prod.Currency())
	}

	// Split
	ps, r, err := m.Split(2)
	if err != nil {
		t.Errorf("err not nil, got %v", err)
	}
	if len(ps) != 2 {
		t.Errorf("ps.len != 2, got %v", len(ps))
	}
	if r.Amount() != 0 {
		t.Errorf("remainder != 0, got %v", r.Amount())
	}
	for i, p := range ps {
		if p.Amount() != 5 {
			t.Errorf("ps[%v].value != 5, got %v", i, p.Amount())
		}
		if !p.Currency().Equals(&EUR) {
			t.Errorf("ps[%v].currency != EUR, got %v", i, p.Currency())
		}
	}

	// Allocation
	ms, rem, err := m.Alloc(11, 11, 11)
	if err != nil {
		t.Errorf("err not nil, got %v", err)
	}
	if rem.Amount() != 1 {
		t.Errorf("expected remainder value to be 1, got %v", rem.Amount())
	}
	if !rem.Currency().Equals(&EUR) {
		t.Errorf("expected currency to be EUR, got %v", rem.Currency())
	}
	if len(ms) != 3 {
		t.Errorf("expected parties to be 3, got %v", len(ms))
	}
	for i := 0; i < 3; i++ {
		if ms[i].Amount() != 3 {
			t.Errorf("expected %v. party allocation value to be 3, got %v", i, ms[i].Amount())
		}
		if !ms[i].Currency().Equals(&EUR) {
			t.Errorf("expected %v. party allocation currency to be EUR, got %v", i, ms[i].Currency())
		}
	}
}

func TestAmount_CurrencyMismatch(t *testing.T) {
	eur := mustAmount(t, 10, EUR)
	usd := mustAmount(t, 10, USD)

	if _, err := eur.Add(usd); !errors.Is(err, ErrCurrencyMismatch) {
		t.Errorf("Add: err = %v, want ErrCurrencyMismatch", err)
	}
	if _, err := eur.Sub(usd); !errors.Is(err, ErrCurrencyMismatch) {
		t.Errorf("Sub: err = %v, want ErrCurrencyMismatch", err)
	}
	if _, err := eur.Cmp(usd); !errors.Is(err, ErrCurrencyMismatch) {
		t.Errorf("Cmp: err = %v, want ErrCurrencyMismatch", err)
	}
	if _, err := eur.GreaterThan(usd); !errors.Is(err, ErrCurrencyMismatch) {
		t.Errorf("GreaterThan: err = %v, want ErrCurrencyMismatch", err)
	}
	if _, err := eur.LessThanOrEqual(usd); !errors.Is(err, ErrCurrencyMismatch) {
		t.Errorf("LessThanOrEqual: err = %v, want ErrCurrencyMismatch", err)
	}
	if eur.SameCurrency(usd) {
		t.Error("SameCurrency(EUR, USD) = true, want false")
	}
	if eur.Equals(usd) {
		t.Error("Equals(EUR 10, USD 10) = true, want false")
	}
}

func TestAmount_Overflow(t *testing.T) {
	max := mustAmount(t, math.MaxInt, EUR)
	min := mustAmount(t, math.MinInt, EUR)
	one := mustAmount(t, 1, EUR)

	if _, err := max.Add(one); !errors.Is(err, ErrInvalidOperation) {
		t.Errorf("Add overflow: err = %v, want ErrInvalidOperation", err)
	}
	if _, err := min.Sub(one); !errors.Is(err, ErrInvalidOperation) {
		t.Errorf("Sub overflow: err = %v, want ErrInvalidOperation", err)
	}
	if _, err := max.Mul(2); !errors.Is(err, ErrInvalidOperation) {
		t.Errorf("Mul overflow: err = %v, want ErrInvalidOperation", err)
	}
	if _, err := min.Abs(); !errors.Is(err, ErrInvalidOperation) {
		t.Errorf("Abs overflow: err = %v, want ErrInvalidOperation", err)
	}
	if _, err := min.Neg(); !errors.Is(err, ErrInvalidOperation) {
		t.Errorf("Neg overflow: err = %v, want ErrInvalidOperation", err)
	}
}

func TestAmount_SplitErrors(t *testing.T) {
	a := mustAmount(t, 10, EUR)
	for _, n := range []int{0, -1} {
		if _, _, err := a.Split(n); !errors.Is(err, ErrInvalidSplitNumber) {
			t.Errorf("Split(%v): err = %v, want ErrInvalidSplitNumber", n, err)
		}
	}
}

func TestAmount_SplitNegativeAmount(t *testing.T) {
	a := mustAmount(t, -10, EUR)
	ps, r, err := a.Split(3)
	if err != nil {
		t.Fatalf("unexpected error %v", err)
	}
	sum := r.Amount()
	for _, p := range ps {
		sum += p.Amount()
	}
	if sum != -10 {
		t.Errorf("parts and remainder sum to %v, want -10", sum)
	}
}

func TestAmount_AllocErrors(t *testing.T) {
	a := mustAmount(t, 10, EUR)

	if _, _, err := a.Alloc(); !errors.Is(err, ErrNoRatiosSpecified) {
		t.Errorf("Alloc(): err = %v, want ErrNoRatiosSpecified", err)
	}
	if _, _, err := a.Alloc(1, -1); !errors.Is(err, ErrInvalidOperation) {
		t.Errorf("Alloc(1, -1): err = %v, want ErrInvalidOperation", err)
	}
	if _, _, err := a.Alloc(0, 0); !errors.Is(err, ErrInvalidOperation) {
		t.Errorf("Alloc(0, 0): err = %v, want ErrInvalidOperation", err)
	}
}

func TestAmount_AllocSumsToOriginal(t *testing.T) {
	a := mustAmount(t, 1001, EUR)
	ms, lo, err := a.Alloc(50, 30, 20)
	if err != nil {
		t.Fatalf("unexpected error %v", err)
	}
	sum := lo.Amount()
	for _, m := range ms {
		sum += m.Amount()
	}
	if sum != 1001 {
		t.Errorf("parts and leftover sum to %v, want 1001", sum)
	}
	if ms[0].Amount() != 500 || ms[1].Amount() != 300 || ms[2].Amount() != 200 {
		t.Errorf("parts = %v %v %v, want 500 300 200",
			ms[0].Amount(), ms[1].Amount(), ms[2].Amount())
	}
	if lo.Amount() != 1 {
		t.Errorf("leftover = %v, want 1", lo.Amount())
	}
}

func TestAmount_Round(t *testing.T) {
	for i, tc := range []struct {
		value int
		cur   Currency
		kind  Kind
		want  int
	}{
		// standard rounding is a no-op for amounts in minor units
		{1234, EUR, Standard, 1234},
		{1234, JPY, Standard, 1234},
		{-1234, EUR, Accounting, -1234},

		// CHF cash rounds to 0.05
		{102, CHF, Cash, 100},
		{103, CHF, Cash, 105},
		{-102, CHF, Cash, -100},
		{-103, CHF, Cash, -105},

		// DKK cash rounds to 0.50
		{1224, DKK, Cash, 1200},
		{1225, DKK, Cash, 1250},

		// SEK cash rounds to whole kronor (cash scale 0)
		{149, SEK, Cash, 100},
		{150, SEK, Cash, 200},

		// USD and JPY have no special cash rounding
		{1234, USD, Cash, 1234},
		{1234, JPY, Cash, 1234},
	} {
		a := mustAmount(t, tc.value, tc.cur)
		r, err := a.Round(tc.kind)
		if err != nil {
			t.Errorf("[%v] %v %v: unexpected error %v", i, tc.value, &tc.cur, err)
			continue
		}
		if r.Amount() != tc.want {
			t.Errorf("[%v] %v %v: got %v, want %v", i, tc.value, &tc.cur, r.Amount(), tc.want)
		}
		if !r.Currency().Equals(&tc.cur) {
			t.Errorf("[%v]: currency changed to %v", i, r.Currency())
		}
	}
}

func TestAmount_Comparisons(t *testing.T) {
	small := mustAmount(t, 5, EUR)
	big := mustAmount(t, 10, EUR)
	alsoBig := mustAmount(t, 10, EUR)

	for i, tc := range []struct {
		a, b    *Amount
		wantCmp int
	}{
		{small, big, -1},
		{big, small, 1},
		{big, alsoBig, 0},
	} {
		c, err := tc.a.Cmp(tc.b)
		if err != nil {
			t.Errorf("[%v] Cmp: unexpected error %v", i, err)
		}
		if c != tc.wantCmp {
			t.Errorf("[%v] Cmp = %v, want %v", i, c, tc.wantCmp)
		}

		gt, _ := tc.a.GreaterThan(tc.b)
		if want := tc.wantCmp > 0; gt != want {
			t.Errorf("[%v] GreaterThan = %v, want %v", i, gt, want)
		}
		gte, _ := tc.a.GreaterThanOrEqual(tc.b)
		if want := tc.wantCmp >= 0; gte != want {
			t.Errorf("[%v] GreaterThanOrEqual = %v, want %v", i, gte, want)
		}
		lt, _ := tc.a.LessThan(tc.b)
		if want := tc.wantCmp < 0; lt != want {
			t.Errorf("[%v] LessThan = %v, want %v", i, lt, want)
		}
		lte, _ := tc.a.LessThanOrEqual(tc.b)
		if want := tc.wantCmp <= 0; lte != want {
			t.Errorf("[%v] LessThanOrEqual = %v, want %v", i, lte, want)
		}
	}

	if !big.Equals(alsoBig) {
		t.Error("Equals(EUR 10, EUR 10) = false, want true")
	}
}

func TestAmount_Predicates(t *testing.T) {
	pos := mustAmount(t, 1, EUR)
	zero := mustAmount(t, 0, EUR)
	negV := mustAmount(t, -1, EUR)

	if !pos.IsPositive() || pos.IsZero() || pos.IsNegative() {
		t.Error("1 should be positive only")
	}
	if zero.IsPositive() || !zero.IsZero() || zero.IsNegative() {
		t.Error("0 should be zero only, neither positive nor negative")
	}
	if negV.IsPositive() || negV.IsZero() || !negV.IsNegative() {
		t.Error("-1 should be negative only")
	}
}

func TestAmount_AbsNeg(t *testing.T) {
	a := mustAmount(t, -5, EUR)

	abs, err := a.Abs()
	if err != nil {
		t.Fatalf("Abs: unexpected error %v", err)
	}
	if abs.Amount() != 5 {
		t.Errorf("Abs = %v, want 5", abs.Amount())
	}

	neg, err := a.Neg()
	if err != nil {
		t.Fatalf("Neg: unexpected error %v", err)
	}
	if neg.Amount() != 5 {
		t.Errorf("Neg(-5) = %v, want 5", neg.Amount())
	}
	neg2, err := neg.Neg()
	if err != nil {
		t.Fatalf("Neg: unexpected error %v", err)
	}
	if neg2.Amount() != -5 {
		t.Errorf("Neg(5) = %v, want -5", neg2.Amount())
	}

	// the original is unchanged
	if a.Amount() != -5 {
		t.Errorf("original mutated to %v", a.Amount())
	}
}

func TestAmount_String(t *testing.T) {
	for i, tc := range []struct {
		value int
		cur   Currency
		want  string
	}{
		{1234, EUR, "EUR 12.34"},
		{-1234, EUR, "EUR -12.34"},
		{5, EUR, "EUR 0.05"},
		{-5, EUR, "EUR -0.05"},
		{0, EUR, "EUR 0.00"},
		{1234, JPY, "JPY 1234"},
		{-1234, JPY, "JPY -1234"},
		{100, USD, "USD 1.00"},
	} {
		a := mustAmount(t, tc.value, tc.cur)
		if got := a.String(); got != tc.want {
			t.Errorf("[%v]: String() = %q, want %q", i, got, tc.want)
		}
	}
}

func TestAmount_JSON(t *testing.T) {
	a := mustAmount(t, 1234, EUR)

	b, err := json.Marshal(a)
	if err != nil {
		t.Fatalf("Marshal: unexpected error %v", err)
	}
	if want := `{"amount":1234,"currency":"EUR"}`; string(b) != want {
		t.Errorf("Marshal = %s, want %s", b, want)
	}

	var back Amount
	if err := json.Unmarshal(b, &back); err != nil {
		t.Fatalf("Unmarshal: unexpected error %v", err)
	}
	if !back.Equals(a) {
		t.Errorf("round trip: got %v, want %v", &back, a)
	}
}

func TestAmount_JSONErrors(t *testing.T) {
	for i, tc := range []struct {
		in      string
		wantErr error
	}{
		{`{"amount":1}`, ErrInvalidJSON},                           // missing currency
		{`{"currency":"EUR"}`, ErrInvalidJSON},                     // missing amount
		{`{"amount":1,"currency":"BTC"}`, ErrISOCodeNotRecognized}, // unknown code
		{`{"amount":1,"currency":"NOPE"}`, ErrISOCodeMalformed},    // not 3 letters
		{`{"amount":1,"currency":"e"}`, ErrISOCodeMalformed},       //
		{`{"amount":1,"currency":"EUR","x":1}`, ErrInvalidJSON},    // unknown field
		{`{"amount":"1","currency":"EUR"}`, ErrInvalidJSON},        // wrong type
	} {
		var a Amount
		if err := json.Unmarshal([]byte(tc.in), &a); !errors.Is(err, tc.wantErr) {
			t.Errorf("[%v] Unmarshal(%s): err = %v, want %v", i, tc.in, err, tc.wantErr)
		}
	}

	// syntactically invalid JSON is rejected by encoding/json before
	// UnmarshalJSON is invoked; it still must surface as an error
	var a Amount
	if err := json.Unmarshal([]byte(`not json`), &a); err == nil {
		t.Error("Unmarshal(not json): expected an error")
	}
}

func TestAmount_JSONRoundTripXXX(t *testing.T) {
	var a Amount // zero value: 0 XXX
	b, err := json.Marshal(&a)
	if err != nil {
		t.Fatalf("Marshal: unexpected error %v", err)
	}
	if want := `{"amount":0,"currency":"XXX"}`; string(b) != want {
		t.Errorf("Marshal = %s, want %s", b, want)
	}
	var back Amount
	if err := json.Unmarshal(b, &back); err != nil {
		t.Fatalf("Unmarshal: unexpected error %v", err)
	}
	if !back.Equals(&a) {
		t.Errorf("round trip: got %v, want %v", &back, &a)
	}
}
