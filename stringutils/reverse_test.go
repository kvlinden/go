package stringutils

import "testing"

func TestReverse(t *testing.T) {
	for _, c := range []struct {
		s, r string
	}{
		{"Hello, world", "dlrow ,olleH"},
		{ "你好", "好你"},
		{"ab", "ba"},
		{"a", "a"},
		{"", ""},
	} {
		if Reverse(c.s) != c.r {
			t.Errorf("Failed to reverse: ", c)
		}
	}
}