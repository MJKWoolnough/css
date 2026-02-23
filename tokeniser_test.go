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
		{ // 3
			"\"a string\"",
			[]parser.Token{
				{Type: TokenString, Data: "\"a string\""},
				{Type: parser.TokenDone},
			},
		},
		{ // 4
			" \"a string with an escape \\20\" ",
			[]parser.Token{
				{Type: TokenWhitespace, Data: " "},
				{Type: TokenString, Data: "\"a string with an escape \\20\""},
				{Type: TokenWhitespace, Data: " "},
				{Type: parser.TokenDone},
			},
		},
		{ // 5
			"'escaped newline \\\n'",
			[]parser.Token{
				{Type: TokenString, Data: "'escaped newline \\\n'"},
				{Type: parser.TokenDone},
			},
		},
		{ // 6
			"'escape followed by newline \\A\n'",
			[]parser.Token{
				{Type: TokenString, Data: "'escape followed by newline \\A\n'"},
				{Type: parser.TokenDone},
			},
		},
		{ // 7
			"'escape followed by newline \\AaFf01\n'",
			[]parser.Token{
				{Type: TokenString, Data: "'escape followed by newline \\AaFf01\n'"},
				{Type: parser.TokenDone},
			},
		},
		{ // 8
			"'escape followed by newline \\AaFf012\n'",
			[]parser.Token{
				{Type: TokenBadString, Data: "'escape followed by newline \\AaFf012\n"},
				{Type: TokenBadString, Data: "'"},
				{Type: parser.TokenDone},
			},
		},
		{ // 9
			"'bad string\n ",
			[]parser.Token{
				{Type: TokenBadString, Data: "'bad string\n"},
				{Type: TokenWhitespace, Data: " "},
				{Type: parser.TokenDone},
			},
		},
		{ // 10
			"'\"'\"'\"",
			[]parser.Token{
				{Type: TokenString, Data: "'\"'"},
				{Type: TokenString, Data: "\"'\""},
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
