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
	identCont    = letters + "_" + digit + "-"
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
		if acceptWordChar(tk) {
			return t.hash(tk)
		}
	} else if tk.Accept("(") {
		t.pushState(')')

		return tk.Return(TokenOpenParen, t.start)
	} else if tk.Accept(")") {
		if t.isState(')') {
			t.popState()

			return tk.Return(TokenCloseParen, t.start)
		}
	} else if tk.Accept(",") {
		return tk.Return(TokenComma, t.start)
	} else if tk.Accept(".") {
		if tk.Accept(digit) {
			return t.number(tk)
		}
	} else if tk.Accept(":") {
		return tk.Return(TokenColon, t.start)
	} else if tk.Accept(";") {
		return tk.Return(TokenSemiColon, t.start)
	} else if tk.Accept("<") {
		s := tk.State()

		if tk.AcceptString("!--", false) == 3 {
			return tk.Return(TokenCDO, t.start)
		}

		s.Reset()
	} else if tk.Accept("@") {
		return t.ident(tk)
	} else if tk.Accept("[") {
		t.pushState(']')

		return tk.Return(TokenOpenBracket, t.start)
	} else if tk.Accept("]") {
		if t.isState(']') {
			t.popState()

			return tk.Return(TokenCloseBracket, t.start)
		}
	} else if tk.Accept("\\") {
		return t.ident(tk)
	} else if tk.Accept("{") {
		t.pushState('}')

		return tk.Return(TokenOpenBrace, t.start)
	} else if tk.Accept("}") {
		if t.isState('}') {
			t.popState()

			return tk.Return(TokenCloseBrace, t.start)
		}
	} else if tk.Accept(digit) {
		return t.number(tk)
	} else if tk.Accept("+") {
		state := tk.State()

		if tk.Accept(digit) || tk.Accept(".") && tk.Accept(digit) {
			return t.number(tk)
		}

		state.Reset()
	} else if tk.Accept("-") {
		state := tk.State()

		if tk.Accept("-") {
			if tk.Accept(">") {
				return tk.Return(TokenCDC, t.start)
			} else {
				return t.ident(tk)
			}
		} else if tk.Accept(identStart) || tk.Accept("\\") {
			return t.ident(tk)
		} else if tk.Accept(digit) || tk.Accept(".") && tk.Accept(digit) {
			return t.number(tk)
		}

		state.Reset()
	} else if tk.Accept(identStart) {
		return t.ident(tk)
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

			if !tk.Accept(noEscape) && !acceptEscape(tk) {
				return tk.Return(TokenBadString, t.start)
			}
		}
	}
}

func acceptEscape(tk *parser.Tokeniser) bool {
	if tk.Accept(noEscape) {
		return false
	} else if tk.Except(hexDigits) {
		return true
	}

	for range 6 {
		tk.Accept(hexDigits)
	}

	tk.Accept(whitespace)

	return true
}

func (t *tokeniser) number(tk *parser.Tokeniser) (parser.Token, parser.TokenFunc) {
	tk.Reset()
	tk.Accept("+-")
	tk.AcceptRun(digit)

	state := tk.State()

	if tk.Accept(".") {
		if tk.Accept(digit) {
			tk.AcceptRun(digit)
		} else {
			state.Reset()
		}
	}

	state = tk.State()

	if tk.Accept("eE") {
		tk.Accept("+-")

		if tk.Accept(digit) {
			tk.AcceptRun(digit)
		} else {
			state.Reset()
		}
	}

	state = tk.State()

	if tk.Accept("%") {
		return tk.Return(TokenPercentage, t.start)
	} else if acceptIdent(tk) {
		return tk.Return(TokenDimension, t.start)
	}

	state.Reset()

	return tk.Return(TokenNumber, t.start)
}

func (t *tokeniser) ident(tk *parser.Tokeniser) (parser.Token, parser.TokenFunc) {
	tk.Reset()

	id := TokenIdent
	state := tk.State()

	if tk.Accept("@") {
		id = TokenAtKeyword
	}

	if !acceptIdent(tk) {
		state.Reset()
		tk.Next()

		return tk.Return(TokenDelim, t.start)
	}

	if id == TokenIdent && tk.Accept("(") {
		id = TokenFunction
		t.pushState(')')
	}

	return tk.Return(id, t.start)
}

func acceptIdent(tk *parser.Tokeniser) bool {
	if !tk.Accept("-") || !tk.Accept("-") {
		if tk.Accept("\\") {
			if !acceptEscape(tk) {
				return false
			}
		} else if !tk.Accept(identStart) {
			if !acceptNonAscii(tk) {
				return false
			}
		}
	}

	for acceptWordChar(tk) {
	}

	return true
}

func acceptNonAscii(tk *parser.Tokeniser) bool {
	if c := tk.Peek(); c < 0x80 {
		return false
	}

	tk.Next()

	return true
}

func acceptWordChar(tk *parser.Tokeniser) bool {
	if tk.Accept(identCont) || acceptNonAscii(tk) {
		return true
	}

	state := tk.State()

	if tk.Accept("\\") && acceptEscape(tk) {
		return true
	}

	state.Reset()

	return false
}

func (t *tokeniser) hash(tk *parser.Tokeniser) (parser.Token, parser.TokenFunc) {
	for acceptWordChar(tk) {
	}

	return tk.Return(TokenHash, t.start)
}
