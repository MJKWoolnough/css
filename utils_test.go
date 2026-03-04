package css

import "testing"

func TestUnquote(t *testing.T) {
	for n, test := range [...]struct {
		Input, Output string
		Err           error
	}{
		{
			Input: ``,
			Err:   ErrBadString,
		},
		{
			Input:  `""`,
			Output: "",
		},
		{
			Input:  `''`,
			Output: "",
		},
		{
			Input: `"`,
			Err:   ErrBadString,
		},
		{
			Input: `'`,
			Err:   ErrBadString,
		},
		{
			Input: `"'`,
			Err:   ErrBadString,
		},
		{
			Input: `''"`,
			Err:   ErrBadString,
		},
		{
			Input:  `"'"`,
			Output: "'",
		},
		{
			Input:  `'"'`,
			Output: "\"",
		},
		{
			Input: `'''`,
			Err:   ErrBadString,
		},
		{
			Input: `"""`,
			Err:   ErrBadString,
		},
		{
			Input:  `"abc"`,
			Output: "abc",
		},
		{
			Input:  `"\20"`,
			Output: " ",
		},
		{
			Input:  `"\41"`,
			Output: "A",
		},
		{
			Input:  `"\3D 123"`,
			Output: "=123",
		},
		{
			Input:  `"\= 123"`,
			Output: "= 123",
		},
		{
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
