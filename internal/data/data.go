package data

import "sort"

// Index converts tags to a compact numeric value
//
// Table represents data stored as string of bytes.
//
// All elements are of size 4. Data may be up to 4 bytes long. Excess bytes can
// be used to store additional information about the data.
type Table string

// Elem returns the element data at the given index
func (s Table) Elem(x int) string {
	return string(s[x*4 : x*4+4])
}

// Index reports the index of the given key or -1 if it could not be found.
// Only the first len(key) bytes from the start of the 4-byte entries will be
// considered for the search and the first match in Index will be returned
func (s Table) Index(key []byte) int {
	n := len(key)
	// search the index of the first entry with an equal or higher value than
	// key in s
	idx := sort.Search(len(s)/4, func(i int) bool {
		return cmp(s[i*4:i*4+n], key) != -1
	})
	i := idx * 4
	if cmp(s[i:i+len(key)], key) != 0 {
		return -1
	}
	return idx
}

// cmp returns an integer comparing a and b lexicographically.
func cmp(a Table, b []byte) int {
	n := len(a)
	if len(b) < n {
		n = len(b)
	}
	for i, c := range b[:n] {
		switch {
		case a[i] > c:
			return 1
		case a[i] < c:
			return -1
		}
	}
	switch {
	case len(a) < len(b):
		return -1
	case len(a) > len(b):
		return 1
	}
	return 0
}

// FixCase reformats b to the same pattern of cases as form.
// It returns false if string b is malformed.
func FixCase(form string, b []byte) bool {
	if len(form) != len(b) {
		return false
	}
	for i, c := range b {
		if form[i] <= 'Z' {
			if c >= 'a' {
				c -= 'z' - 'Z'
			}
			if c < 'A' || 'Z' < c {
				return false
			}
		} else {
			if c <= 'Z' {
				c += 'z' - 'Z'
			}
			if c < 'a' || 'z' < c {
				return false
			}
		}
		b[i] = c
	}
	return true
}
