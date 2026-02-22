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

const (
	TokenWhitespace parser.TokenType = iota
)

type preprocessor struct {
	parser.Tokeniser
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

func CreateTokeniser(t parser.Tokeniser) *parser.Tokeniser {
	t = parser.NewRuneReaderTokeniser(&preprocessor{t})

	t.TokeniserState(new(tokeniser).start)

	return &t
}

type tokeniser struct {
	depth []rune
}

func (t *tokeniser) start(tk *parser.Tokeniser) (parser.Token, parser.TokenFunc) {
	if tk.Peek() == -1 {
		if len(t.depth) == 0 {
			return tk.Done()
		}

		return tk.ReturnError(io.ErrUnexpectedEOF)
	}

	if tk.Accept(whitespace) {
		tk.AcceptRun(whitespace)

		return tk.Return(TokenWhitespace, t.start)
	}

	return tk.ReturnError(nil)
}
