package css

import (
	"io"

	"vimagination.zapto.org/parser"
)

const (
	whitespace   = " \t\r\n\f"
	digit        = "0123456789"
	upperLetters = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	lowerLetters = "abcdefghijklmnopqrstuvwxyz"
	letters      = upperLetters + lowerLetters
	identStart   = letters + "_"
)

const (
	TokenWhitespace parser.TokenType = iota
	TokenIdent
	TokenString
	TokenHash
	TokenNumber
	TokenComma
	TokenCDC
	TokenCDO
	TokenColon
	TokenSemiColon
	TokenAtKeyword
	TokenOpenParen
	TokenCloseParen
	TokenOpenBracket
	TokenCloseBracket
	TokenOpenBrace
	TokenCloseBrace
	TokenFunction
	TokenBadString
	TokenURL
	TokenBadURL
	TokenPercentage
	TokenDimension
	TokenDelim
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

func (t *tokeniser) isState(r rune) bool {
	if len(t.depth) == 0 {
		return false
	}

	return t.depth[len(t.depth)-1] == r
}

func (t *tokeniser) pushState(r rune) {
	t.depth = append(t.depth, r)
}

func (t *tokeniser) popState() {
	t.depth = t.depth[:len(t.depth)-1]
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
	} else if tk.Accept("\"") {
	} else if tk.Accept("#") {
	} else if tk.Accept("'") {
	} else if tk.Accept("(") {
		t.pushState(')')

		return tk.Return(TokenOpenParen, t.start)
	} else if tk.Accept(")") {
		if t.isState(')') {
			t.popState()

			return tk.Return(TokenCloseParen, t.start)
		}
	} else if tk.Accept("+") {
	} else if tk.Accept(",") {
		return tk.Return(TokenComma, t.start)
	} else if tk.Accept("-") {
	} else if tk.Accept(".") {
	} else if tk.Accept(":") {
		return tk.Return(TokenColon, t.start)
	} else if tk.Accept(";") {
		return tk.Return(TokenSemiColon, t.start)
	} else if tk.Accept("<") {
		s := tk.State()

		if tk.Accept("-") && tk.Accept("-") {
			return tk.Return(TokenCDO, t.start)
		}

		s.Reset()
	} else if tk.Accept("@") {
	} else if tk.Accept("[") {
		t.pushState(']')

		return tk.Return(TokenOpenBracket, t.start)
	} else if tk.Accept("]") {
		if t.isState(']') {
			t.popState()

			return tk.Return(TokenCloseBracket, t.start)
		}
	} else if tk.Accept("\\") {
	} else if tk.Accept("{") {
		t.pushState('}')

		return tk.Return(TokenOpenBrace, t.start)
	} else if tk.Accept("}") {
		if t.isState('}') {
			t.popState()

			return tk.Return(TokenCloseBrace, t.start)
		}
	} else if tk.Accept(digit) {
	} else if tk.Accept(identStart) {
	}

	return tk.Return(TokenDelim, t.start)
}
