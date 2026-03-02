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
			" \t\n\r\r\n\f",
			[]parser.Token{
				{Type: TokenWhitespace, Data: " \t\n\n\n\n"},
				{Type: parser.TokenDone},
			},
		},
		{ // 2
			"/* A Comment *//* Another Comment */",
			[]parser.Token{
				{Type: TokenComment, Data: "/* A Comment */"},
				{Type: TokenComment, Data: "/* Another Comment */"},
				{Type: parser.TokenDone},
			},
		},
		{ // 3
			"/* A Comment",
			[]parser.Token{
				{Type: parser.TokenError, Data: "unexpected EOF"},
			},
		},
		{ // 4
			" /* A Comment */\n \t",
			[]parser.Token{
				{Type: TokenWhitespace, Data: " "},
				{Type: TokenComment, Data: "/* A Comment */"},
				{Type: TokenWhitespace, Data: "\n \t"},
				{Type: parser.TokenDone},
			},
		},
		{ // 5
			"\"a string\"",
			[]parser.Token{
				{Type: TokenString, Data: "\"a string\""},
				{Type: parser.TokenDone},
			},
		},
		{ // 6
			" \"a string with an escape \\20\" ",
			[]parser.Token{
				{Type: TokenWhitespace, Data: " "},
				{Type: TokenString, Data: "\"a string with an escape \\20\""},
				{Type: TokenWhitespace, Data: " "},
				{Type: parser.TokenDone},
			},
		},
		{ // 7
			"'escaped newline \\\n'",
			[]parser.Token{
				{Type: TokenString, Data: "'escaped newline \\\n'"},
				{Type: parser.TokenDone},
			},
		},
		{ // 8
			"'escape followed by newline \\A\n'",
			[]parser.Token{
				{Type: TokenString, Data: "'escape followed by newline \\A\n'"},
				{Type: parser.TokenDone},
			},
		},
		{ // 9
			"'escape followed by newline \\AaFf01\n'",
			[]parser.Token{
				{Type: TokenString, Data: "'escape followed by newline \\AaFf01\n'"},
				{Type: parser.TokenDone},
			},
		},
		{ // 10
			"'escaped newline \\\n'",
			[]parser.Token{
				{Type: TokenString, Data: "'escaped newline \\\n'"},
				{Type: parser.TokenDone},
			},
		},
		{ // 11
			"'escaped EOF \\",
			[]parser.Token{
				{Type: TokenBadString, Data: "'escaped EOF \\"},
				{Type: parser.TokenDone},
			},
		},
		{ // 12
			"'escape followed by newline \\AaFf012\n'",
			[]parser.Token{
				{Type: TokenBadString, Data: "'escape followed by newline \\AaFf012\n"},
				{Type: TokenBadString, Data: "'"},
				{Type: parser.TokenDone},
			},
		},
		{ // 13
			"'bad string\n ",
			[]parser.Token{
				{Type: TokenBadString, Data: "'bad string\n"},
				{Type: TokenWhitespace, Data: " "},
				{Type: parser.TokenDone},
			},
		},
		{ // 14
			"'\"'\"'\"",
			[]parser.Token{
				{Type: TokenString, Data: "'\"'"},
				{Type: TokenString, Data: "\"'\""},
				{Type: parser.TokenDone},
			},
		},
		{ // 15
			"{}[",
			[]parser.Token{
				{Type: TokenOpenBrace, Data: "{"},
				{Type: TokenCloseBrace, Data: "}"},
				{Type: TokenOpenBracket, Data: "["},
				{Type: parser.TokenError, Data: "unexpected EOF"},
			},
		},
		{ // 16
			")[(]",
			[]parser.Token{
				{Type: TokenDelim, Data: ")"},
				{Type: TokenOpenBracket, Data: "["},
				{Type: TokenOpenParen, Data: "("},
				{Type: TokenDelim, Data: "]"},
				{Type: parser.TokenError, Data: "unexpected EOF"},
			},
		},
		{ // 17
			"[abc]",
			[]parser.Token{
				{Type: TokenOpenBracket, Data: "["},
				{Type: TokenIdent, Data: "abc"},
				{Type: TokenCloseBracket, Data: "]"},
				{Type: parser.TokenDone},
			},
		},
		{ // 18
			",:;",
			[]parser.Token{
				{Type: TokenComma, Data: ","},
				{Type: TokenColon, Data: ":"},
				{Type: TokenSemiColon, Data: ";"},
				{Type: parser.TokenDone},
			},
		},
		{ // 19
			"1 2 12 +3.14 -1 10e+2 1.2E-9 .5 .+- 123.a 123e-a 456E+a 789ea",
			[]parser.Token{
				{Type: TokenNumber, Data: "1"},
				{Type: TokenWhitespace, Data: " "},
				{Type: TokenNumber, Data: "2"},
				{Type: TokenWhitespace, Data: " "},
				{Type: TokenNumber, Data: "12"},
				{Type: TokenWhitespace, Data: " "},
				{Type: TokenNumber, Data: "+3.14"},
				{Type: TokenWhitespace, Data: " "},
				{Type: TokenNumber, Data: "-1"},
				{Type: TokenWhitespace, Data: " "},
				{Type: TokenNumber, Data: "10e+2"},
				{Type: TokenWhitespace, Data: " "},
				{Type: TokenNumber, Data: "1.2E-9"},
				{Type: TokenWhitespace, Data: " "},
				{Type: TokenNumber, Data: ".5"},
				{Type: TokenWhitespace, Data: " "},
				{Type: TokenDelim, Data: "."},
				{Type: TokenDelim, Data: "+"},
				{Type: TokenDelim, Data: "-"},
				{Type: TokenWhitespace, Data: " "},
				{Type: TokenNumber, Data: "123"},
				{Type: TokenDelim, Data: "."},
				{Type: TokenIdent, Data: "a"},
				{Type: TokenWhitespace, Data: " "},
				{Type: TokenDimension, Data: "123e-a"},
				{Type: TokenWhitespace, Data: " "},
				{Type: TokenDimension, Data: "456E"},
				{Type: TokenDelim, Data: "+"},
				{Type: TokenIdent, Data: "a"},
				{Type: TokenWhitespace, Data: " "},
				{Type: TokenDimension, Data: "789ea"},
				{Type: parser.TokenDone},
			},
		},
		{ // 20
			"1% 2% 12% +3.14% 10e+2% .5%",
			[]parser.Token{
				{Type: TokenPercentage, Data: "1%"},
				{Type: TokenWhitespace, Data: " "},
				{Type: TokenPercentage, Data: "2%"},
				{Type: TokenWhitespace, Data: " "},
				{Type: TokenPercentage, Data: "12%"},
				{Type: TokenWhitespace, Data: " "},
				{Type: TokenPercentage, Data: "+3.14%"},
				{Type: TokenWhitespace, Data: " "},
				{Type: TokenPercentage, Data: "10e+2%"},
				{Type: TokenWhitespace, Data: " "},
				{Type: TokenPercentage, Data: ".5%"},
				{Type: parser.TokenDone},
			},
		},
		{ // 21
			"1a 2abc123 12-A_b\\n +3.14--123 10e+2-\\n\\n .5\\n 10px 10 px",
			[]parser.Token{
				{Type: TokenDimension, Data: "1a"},
				{Type: TokenWhitespace, Data: " "},
				{Type: TokenDimension, Data: "2abc123"},
				{Type: TokenWhitespace, Data: " "},
				{Type: TokenDimension, Data: "12-A_b\\n"},
				{Type: TokenWhitespace, Data: " "},
				{Type: TokenDimension, Data: "+3.14--123"},
				{Type: TokenWhitespace, Data: " "},
				{Type: TokenDimension, Data: "10e+2-\\n\\n"},
				{Type: TokenWhitespace, Data: " "},
				{Type: TokenDimension, Data: ".5\\n"},
				{Type: TokenWhitespace, Data: " "},
				{Type: TokenDimension, Data: "10px"},
				{Type: TokenWhitespace, Data: " "},
				{Type: TokenNumber, Data: "10"},
				{Type: TokenWhitespace, Data: " "},
				{Type: TokenIdent, Data: "px"},
				{Type: parser.TokenDone},
			},
		},
		{ // 22
			"a abc123 -A_b\\n --123 -\\n\\n \\n abc\\\n abc£def",
			[]parser.Token{
				{Type: TokenIdent, Data: "a"},
				{Type: TokenWhitespace, Data: " "},
				{Type: TokenIdent, Data: "abc123"},
				{Type: TokenWhitespace, Data: " "},
				{Type: TokenIdent, Data: "-A_b\\n"},
				{Type: TokenWhitespace, Data: " "},
				{Type: TokenIdent, Data: "--123"},
				{Type: TokenWhitespace, Data: " "},
				{Type: TokenIdent, Data: "-\\n\\n"},
				{Type: TokenWhitespace, Data: " "},
				{Type: TokenIdent, Data: "\\n"},
				{Type: TokenWhitespace, Data: " "},
				{Type: TokenIdent, Data: "abc"},
				{Type: TokenDelim, Data: "\\"},
				{Type: TokenWhitespace, Data: "\n "},
				{Type: TokenIdent, Data: "abc£def"},
				{Type: parser.TokenDone},
			},
		},
		{ // 23
			"@a @abc123 @-A_b\\n @--123 @-\\n\\n @\\n",
			[]parser.Token{
				{Type: TokenAtKeyword, Data: "@a"},
				{Type: TokenWhitespace, Data: " "},
				{Type: TokenAtKeyword, Data: "@abc123"},
				{Type: TokenWhitespace, Data: " "},
				{Type: TokenAtKeyword, Data: "@-A_b\\n"},
				{Type: TokenWhitespace, Data: " "},
				{Type: TokenAtKeyword, Data: "@--123"},
				{Type: TokenWhitespace, Data: " "},
				{Type: TokenAtKeyword, Data: "@-\\n\\n"},
				{Type: TokenWhitespace, Data: " "},
				{Type: TokenAtKeyword, Data: "@\\n"},
				{Type: parser.TokenDone},
			},
		},
		{ // 24
			"<!-- --><!----><abc>",
			[]parser.Token{
				{Type: TokenCDO, Data: "<!--"},
				{Type: TokenWhitespace, Data: " "},
				{Type: TokenCDC, Data: "-->"},
				{Type: TokenCDO, Data: "<!--"},
				{Type: TokenCDC, Data: "-->"},
				{Type: TokenDelim, Data: "<"},
				{Type: TokenIdent, Data: "abc"},
				{Type: TokenDelim, Data: ">"},
				{Type: parser.TokenDone},
			},
		},
		{ // 25
			"#a #abc123 #-A_b\\n #--123 #-\\n\\n #\\n",
			[]parser.Token{
				{Type: TokenHash, Data: "#a"},
				{Type: TokenWhitespace, Data: " "},
				{Type: TokenHash, Data: "#abc123"},
				{Type: TokenWhitespace, Data: " "},
				{Type: TokenHash, Data: "#-A_b\\n"},
				{Type: TokenWhitespace, Data: " "},
				{Type: TokenHash, Data: "#--123"},
				{Type: TokenWhitespace, Data: " "},
				{Type: TokenHash, Data: "#-\\n\\n"},
				{Type: TokenWhitespace, Data: " "},
				{Type: TokenHash, Data: "#\\n"},
				{Type: parser.TokenDone},
			},
		},
		{ // 26
			"a()abc123() -A_b\\n() @--123()",
			[]parser.Token{
				{Type: TokenFunction, Data: "a("},
				{Type: TokenCloseParen, Data: ")"},
				{Type: TokenFunction, Data: "abc123("},
				{Type: TokenCloseParen, Data: ")"},
				{Type: TokenWhitespace, Data: " "},
				{Type: TokenFunction, Data: "-A_b\\n("},
				{Type: TokenCloseParen, Data: ")"},
				{Type: TokenWhitespace, Data: " "},
				{Type: TokenAtKeyword, Data: "@--123"},
				{Type: TokenOpenParen, Data: "("},
				{Type: TokenCloseParen, Data: ")"},
				{Type: parser.TokenDone},
			},
		},
		{ // 27
			"url(abc) uRl( abc ) URL() UrL(!#$%&) url(abc\") url('abc') url(\"\") url(a b) url(a\\\nb) url(abc\\)",
			[]parser.Token{
				{Type: TokenURL, Data: "url(abc)"},
				{Type: TokenWhitespace, Data: " "},
				{Type: TokenURL, Data: "uRl( abc )"},
				{Type: TokenWhitespace, Data: " "},
				{Type: TokenURL, Data: "URL()"},
				{Type: TokenWhitespace, Data: " "},
				{Type: TokenURL, Data: "UrL(!#$%&)"},
				{Type: TokenWhitespace, Data: " "},
				{Type: TokenBadURL, Data: "url(abc\")"},
				{Type: TokenWhitespace, Data: " "},
				{Type: TokenFunction, Data: "url("},
				{Type: TokenString, Data: "'abc'"},
				{Type: TokenCloseParen, Data: ")"},
				{Type: TokenWhitespace, Data: " "},
				{Type: TokenFunction, Data: "url("},
				{Type: TokenString, Data: "\"\""},
				{Type: TokenCloseParen, Data: ")"},
				{Type: TokenWhitespace, Data: " "},
				{Type: TokenBadURL, Data: "url(a b)"},
				{Type: TokenWhitespace, Data: " "},
				{Type: TokenBadURL, Data: "url(a\\\nb)"},
				{Type: TokenWhitespace, Data: " "},
				{Type: TokenBadURL, Data: "url(abc\\)"},
				{Type: parser.TokenDone},
			},
		},
		{ // 28
			"@ # . @#.|!$&",
			[]parser.Token{
				{Type: TokenDelim, Data: "@"},
				{Type: TokenWhitespace, Data: " "},
				{Type: TokenDelim, Data: "#"},
				{Type: TokenWhitespace, Data: " "},
				{Type: TokenDelim, Data: "."},
				{Type: TokenWhitespace, Data: " "},
				{Type: TokenDelim, Data: "@"},
				{Type: TokenDelim, Data: "#"},
				{Type: TokenDelim, Data: "."},
				{Type: TokenDelim, Data: "|"},
				{Type: TokenDelim, Data: "!"},
				{Type: TokenDelim, Data: "$"},
				{Type: TokenDelim, Data: "&"},
				{Type: parser.TokenDone},
			},
		},
	} {
		p := createTokeniser(parser.NewStringTokeniser(test.Input), true)

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
