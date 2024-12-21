package main

import "testing"

func TestCleanInput(t *testing.T) {
	// build some test cases
	cases := []struct {
		input    string
		expected []string
	}{
		{
			input:    "  hello  world  ",
			expected: []string{"hello", "world"},
		},
		{
			input:    "   ",
			expected: []string{},
		},
		{
			input:    "Giga, Mega KILO",
			expected: []string{"giga", "mega", "kilo"},
		},
	}

	// run the test case
	for _, c := range cases {
		actual := cleanInput(c.input)
		if len(actual) != len(c.expected) {
			t.Errorf("cleanInput(%q) expected %d elements, got %d", c.input, len(c.expected), len(actual))
		}
		for i, word := range actual {
			if word != c.expected[i] {
				t.Errorf("cleanInput(%v) expected %v, got %v", c.input, c.expected[i], word)
			}
		}
	}
	return
}
