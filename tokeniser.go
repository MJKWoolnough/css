package css

import (
	"io"

	"vimagination.zapto.org/parser"
)

const (
	whitespace   = " \t\n"
	digit        = "0123456789"
	upperLetters = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	lowerLetters = "abcdefghijklmnopqrstuvwxyz"
	letters      = upperLetters + lowerLetters
	identStart   = letters + "_"
	hexDigits    = digit + "abcdefABCDEF"
	noEscape     = "\n"
)

const (
	TokenWhitespace parser.TokenType = iota
	TokenComment
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
	case '\f':
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

	if tk.Accept("/") {
		if tk.Accept("*") {
			return t.parseComment(tk)
		}
	} else if tk.Accept(whitespace) {
		tk.AcceptRun(whitespace)

		return tk.Return(TokenWhitespace, t.start)
	} else if tk.Accept(`"`) {
		return t.string(tk)
	} else if tk.Accept("'") {
		return t.string(tk)
	} else if tk.Accept("#") {
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

func (t *tokeniser) parseComment(tk *parser.Tokeniser) (parser.Token, parser.TokenFunc) {
	for {
		if tk.ExceptRun("*") == -1 {
			return tk.ReturnError(io.ErrUnexpectedEOF)
		}

		tk.Accept("*")

		if tk.Accept("/") {
			return tk.Return(TokenComment, t.start)
		}
	}
}

func (t *tokeniser) string(tk *parser.Tokeniser) (parser.Token, parser.TokenFunc) {
	tk.Reset()

	var chars string

	switch tk.Next() {
	case '"':
		chars = "\"\\\n"
	case '\'':
		chars = "'\\\n"
	}

	for {
		switch tk.ExceptRun(chars) {
		case '\n':
			tk.Next()

			fallthrough
		case -1:
			return tk.Return(TokenBadString, t.start)
		case '"', '\'':
			tk.Next()

			return tk.Return(TokenString, t.start)
		case '\\':
			tk.Next()

			if !tk.Accept(noEscape) && tk.Accept(hexDigits) {
				for range 5 {
					tk.Accept(hexDigits)
				}

				tk.Accept(whitespace)
			}
		}
	}
}
