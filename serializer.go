package money

import (
	"bytes"
	"encoding/json"
	"fmt"
)

// Injection points for backward compatibility.
// If you need to use a different marshal/unmarshal way,
// overwrite this like following:
// 	money.UnmarshalJSON = func(m *Money, b []byte) error { ... }
// 	money.MarshalJSON = func(m Money) ([]byte, error) { return nil, nil }
var (
	UnmarshalJSON = defaultUnmarshalJSON
	MarshalJSON   = defaultMarshalJSON
)

func defaultUnmarshalJSON(m *Money, b []byte) error {
	data := make(map[string]any)
	err := json.Unmarshal(b, &data)
	if err != nil {
		return err
	}

	var amount float64
	if amountRaw, ok := data["amount"]; ok {
		amount, ok = amountRaw.(float64)
		if !ok {
			return ErrInvalidJSONUnmarshal
		}
	}
	
	var currency string
	if currencyRaw, ok := data["currency"]; ok {
		currency, ok = currencyRaw.(string)
		if !ok {
			return ErrInvalidJSONUnmarshal
		}
	}
	
	var ref *Money
	if amount == 0 && currency == "" {
		ref = &Money{}
		} else {
			ref = New(int64(amount), currency)
		}
		*m = *ref
		return nil
	}

func defaultMarshalJSON(m Money) ([]byte, error) {
	if m == (Money{}) {
		m = *New(0, "")
	}
	buff := bytes.NewBufferString(fmt.Sprintf(`{"amount": %d, "currency": "%s"}`, m.Amount(), m.CurrencyCode()))
	return buff.Bytes(), nil
}
