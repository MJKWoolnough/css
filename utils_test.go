package css

import (
	"io"
	"testing"
)

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
		{ // 18
			Input: `"\"`,
			Err:   io.ErrUnexpectedEOF,
		},
	} {
		if out, err := Unquote(test.Input); err != test.Err {
			t.Errorf("test %d: expecting error %v, got %v", n+1, test.Err, err)
		} else if out != test.Output {
			t.Errorf("test %d: expecting output %q, got %q", n+1, test.Output, out)
		}
	}
}

func TestUnURL(t *testing.T) {
	for n, test := range [...]struct {
		Input, Output string
		Err           error
	}{
		{ // 1
			Input: ``,
			Err:   ErrBadURL,
		},
		{ // 2
			Input:  `url()`,
			Output: ``,
		},
		{ // 3
			Input:  `url(abc)`,
			Output: `abc`,
		},
		{ // 4
			Input:  `URL( abc )`,
			Output: `abc`,
		},
		{ // 5
			Input:  `URL( abc )`,
			Output: `abc`,
		},
		{ // 6
			Input:  `URL( \3D \4D )`,
			Output: `=M`,
		},
		{ // 7
			Input: `url(\)`,
			Err:   io.ErrUnexpectedEOF,
		},
		{ // 8
			Input: `URL( abc def )`,
			Err:   ErrBadURL,
		},
	} {
		if out, err := UnURL(test.Input); err != test.Err {
			t.Errorf("test %d: expecting error %v, got %v", n+1, test.Err, err)
		} else if out != test.Output {
			t.Errorf("test %d: expecting output %q, got %q", n+1, test.Output, out)
		}
	}
}
