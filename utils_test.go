package css

import "testing"

func TestUnquote(t *testing.T) {
	for n, test := range [...]struct {
		Input, Output string
		Err           error
	}{
		{ // 1
			Input: ``,
			Err:   ErrBadString,
		},
		{ // 2
			Input:  `""`,
			Output: "",
		},
		{ // 3
			Input:  `''`,
			Output: "",
		},
		{ // 4
			Input: `"`,
			Err:   ErrBadString,
		},
		{ // 5
			Input: `'`,
			Err:   ErrBadString,
		},
		{ // 6
			Input: `"'`,
			Err:   ErrBadString,
		},
		{ // 7
			Input: `''"`,
			Err:   ErrBadString,
		},
		{ // 8
			Input:  `"'"`,
			Output: "'",
		},
		{ // 9
			Input:  `'"'`,
			Output: "\"",
		},
		{ // 10
			Input: `'''`,
			Err:   ErrBadString,
		},
		{ // 11
			Input: `"""`,
			Err:   ErrBadString,
		},
		{ // 12
			Input:  `"abc"`,
			Output: "abc",
		},
		{ // 13
			Input:  `"\20"`,
			Output: " ",
		},
		{ // 14
			Input:  `"\41"`,
			Output: "A",
		},
		{ // 15
			Input:  `"\3D 123"`,
			Output: "=123",
		},
		{ // 16
			Input:  `"\= 123"`,
			Output: "= 123",
		},
		{ // 17
			Input:  `"\=\4D A"`,
			Output: "=MA",
		},
	} {
		if out, err := Unquote(test.Input); err != test.Err {
			t.Errorf("test %d: expecting error %v, got %v", n+1, test.Err, err)
		} else if out != test.Output {
			t.Errorf("test %d: expecting output %q, got %q", n+1, test.Output, out)
		}
	}
}
