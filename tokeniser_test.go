package css

import (
	"testing"

	"vimagination.zapto.org/parser"
)

func TestTokeniser(t *testing.T) {
	for n, test := range [...]struct {
		Input  string
		Output []parser.Token
	}{
		{ // 1
			"/* A Comment *//* Another Comment */",
			[]parser.Token{
				{Type: TokenComment, Data: "/* A Comment */"},
				{Type: TokenComment, Data: "/* Another Comment */"},
				{Type: parser.TokenDone},
			},
		},
		{ // 2
			" /* A Comment */\n \t",
			[]parser.Token{
				{Type: TokenWhitespace, Data: " "},
				{Type: TokenComment, Data: "/* A Comment */"},
				{Type: TokenWhitespace, Data: "\n \t"},
				{Type: parser.TokenDone},
			},
		},
	} {
		p := CreateTokeniser(parser.NewStringTokeniser(test.Input))

		for m, tkn := range test.Output {
			if tk, _ := p.GetToken(); tk.Type != tkn.Type {
				if tk.Type == parser.TokenError {
					t.Errorf("test %d.%d: unexpected error: %s", n+1, m+1, tk.Data)
				} else {
					t.Errorf("test %d.%d: Incorrect type, expecting %d, got %d", n+1, m+1, tkn.Type, tk.Type)
				}

				break
			} else if tk.Data != tkn.Data {
				t.Errorf("test %d.%d: Incorrect data, expecting %q, got %q", n+1, m+1, tkn.Data, tk.Data)

				break
			}
		}
	}
}
