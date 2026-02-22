package css

import (
	"io"

	"vimagination.zapto.org/parser"
)

const (
	whitespace   = " \t\r\n\f"
	digits       = "0123456789"
	upperLetters = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	lowerLetters = "abcdefghijklmnopqrstuvwxyz"
	letters      = upperLetters + lowerLetters
)

type preprocessor struct {
	*parser.Tokeniser
}

func (p *preprocessor) ReadRune() (rune, int, error) {
	r := p.Next()
	if r == -1 {
		return 0, 0, io.EOF
	}

	switch r {
	case '\r':
		p.Accept("\n")

		r = '\n'
	case '\x0c':
		r = '\n'
	}

	return r, 0, nil
}

func CreateTokeniser(t *parser.Tokeniser) *parser.Tokeniser {
	tk := parser.NewRuneReaderTokeniser(&preprocessor{t})

	return &tk
}
