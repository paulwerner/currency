package data

import "testing"

// table mirrors the layout of the generated currency table: sorted 4-byte
// entries, a leading dummy so indices start at 1, and a trailing sentinel
const table Table = "\x00\x00\x00\x00" + "AAA\x00" + "BBB\x01" + "CCC\x02" + "\xff\xff\xff\xff"

func TestTable_Elem(t *testing.T) {
	for i, want := range []string{"\x00\x00\x00\x00", "AAA\x00", "BBB\x01", "CCC\x02", "\xff\xff\xff\xff"} {
		if got := table.Elem(i); got != want {
			t.Errorf("Elem(%v) = %q, want %q", i, got, want)
		}
	}
}

func TestTable_Index(t *testing.T) {
	for _, tc := range []struct {
		key  string
		want int
	}{
		{"AAA", 1},
		{"BBB", 2},
		{"CCC", 3},
		{"ABC", -1},
		{"ZZZ", -1},
	} {
		if got := table.Index([]byte(tc.key)); got != tc.want {
			t.Errorf("Index(%q) = %v, want %v", tc.key, got, tc.want)
		}
	}
}

func TestFixCase(t *testing.T) {
	for _, tc := range []struct {
		form   string
		in     string
		want   string
		wantOk bool
	}{
		{"XXX", "usd", "USD", true},
		{"XXX", "USD", "USD", true},
		{"XXX", "uSd", "USD", true},
		{"xxx", "USD", "usd", true},
		{"XXX", "US1", "", false},
		{"XXX", "US~", "", false},
		{"XXX", "US", "", false},   // too short
		{"XXX", "USDX", "", false}, // too long
		{"XXX", "", "", false},
	} {
		b := []byte(tc.in)
		ok := FixCase(tc.form, b)
		if ok != tc.wantOk {
			t.Errorf("FixCase(%q, %q) = %v, want %v", tc.form, tc.in, ok, tc.wantOk)
		}
		if ok && string(b) != tc.want {
			t.Errorf("FixCase(%q, %q) rewrote to %q, want %q", tc.form, tc.in, b, tc.want)
		}
	}
}
