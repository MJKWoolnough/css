package css

import (
	"fmt"
	"slices"
	"strings"

	"vimagination.zapto.org/parser"
)

// Token represents a single parsed token with source positioning.
type Token struct {
	parser.Token
	Pos, Line, LinePos uint64
}

// Tokens is a collection of Token values.
type Tokens []Token

// Comments is a collection of Comment Tokens.
type Comments []*Token

type cssParser Tokens

// Tokeniser is an interface representing a tokeniser.
type Tokeniser interface {
	TokeniserState(parser.TokenFunc)
	Iter(func(parser.Token) bool)
	GetError() error
}

func newCSSParser(t Tokeniser) (cssParser, error) {
	var (
		tokens             cssParser
		pos, line, linePos uint64
		err                error
	)

	for tk := range t.Iter {
		typ := tk.Type

		tokens = append(tokens, Token{
			Token: parser.Token{
				Type: typ,
				Data: tk.Data,
			},
			Pos:     pos,
			Line:    line,
			LinePos: linePos,
		})

		switch typ {
		case parser.TokenError:
			err = Error{
				Err:     t.GetError(),
				Parsing: "Tokens",
				Token:   tokens[len(tokens)-1],
			}
		case TokenWhitespace:
			var (
				lastLT   int
				lastChar rune
			)

			for n, c := range tk.Data {
				if strings.ContainsRune(whitespace, c) {
					lastLT = n + 1
					linePos = 0

					if lastChar != '\r' || c != '\n' {
						line++
					}
				}

				lastChar = c
			}

			linePos += uint64(len(tk.Data) - lastLT)
		default:
			linePos += uint64(len(tk.Data))
		}

		pos += uint64(len(tk.Data))
	}

	return tokens[0:0:len(tokens)], err
}

func (c cssParser) NewGoal() cssParser {
	return c[len(c):]
}

func (c *cssParser) Score(k cssParser) {
	*c = (*c)[:len(*c)+len(k)]
}

func (c *cssParser) next() *Token {
	l := len(*c)
	if l == cap(*c) {
		return &(*c)[l-1]
	}

	*c = (*c)[:l+1]
	tk := (*c)[l]

	return &tk
}

func (c *cssParser) backup() {
	*c = (*c)[:len(*c)-1]
}

func (c *cssParser) Peek() parser.Token {
	tk := c.next().Token

	c.backup()

	return tk
}

func (c *cssParser) Accept(ts ...parser.TokenType) bool {
	if slices.Contains(ts, c.next().Type) {
		return true
	}

	c.backup()

	return false
}

func (c *cssParser) AcceptRun(ts ...parser.TokenType) parser.TokenType {
Loop:
	for {
		tt := c.next().Type

		for _, pt := range ts {
			if pt == tt {
				continue Loop
			}
		}

		c.backup()

		return tt
	}
}

func (c *cssParser) Skip() {
	c.next()
}

func (c *cssParser) Next() *Token {
	return c.next()
}

var depths = [...][2]parser.Token{
	{{Type: TokenOpenBracket, Data: "["}, {Type: TokenCloseBracket, Data: "]"}},
	{{Type: TokenOpenParen, Data: "("}, {Type: TokenCloseParen, Data: ")"}},
	{{Type: TokenOpenBrace, Data: "{"}, {Type: TokenCloseBrace, Data: "}"}},
}

func (c *cssParser) SkipDepth() bool {
	var (
		on    = -1
		depth = 1
	)

	for n, d := range depths {
		if c.AcceptToken(d[0]) {
			on = n

			break
		}
	}

	if on == -1 {
		return false
	}

	for depth > 0 {
		if c.AcceptToken(depths[on][0]) {
			depth++
		} else if c.AcceptToken(depths[on][1]) {
			depth--
		} else {
			c.Skip()
		}
	}

	return true
}

func (c *cssParser) AcceptToken(tk parser.Token) bool {
	if c.next().Token == tk {
		return true
	}

	c.backup()

	return false
}

func (c *cssParser) ToTokens() Tokens {
	return Tokens((*c)[:len(*c):len(*c)])
}

func (c *cssParser) AcceptRunWhitespace() parser.TokenType {
	return c.AcceptRun(TokenWhitespace, TokenComment)
}

func (c *cssParser) AcceptRunWhitespaceNoNewLine() {
	for {
		d := c.NewGoal()

		if d.Accept(TokenWhitespace) {
			if l := d.GetLastToken().Data; l != "\n" && l != "\r\n" {
				break
			}
		}

		c.Skip()
	}
}

func (c *cssParser) AcceptRunWhitespaceNoComment() parser.TokenType {
	return c.AcceptRun(TokenWhitespace)
}

func (c *cssParser) AcceptRunWhitespaceComments() Comments {
	var cs Comments

	d := c.NewGoal()

Loop:
	for {
		switch d.AcceptRunWhitespaceNoComment() {
		case TokenComment:
		default:
			break Loop
		}

		cs = append(cs, d.Next())

		c.Score(d)

		d = c.NewGoal()
	}

	return cs
}

func (c *cssParser) AcceptRunWhitespaceNoNewLineNoComment() parser.TokenType {
	return c.AcceptRun(TokenWhitespace)
}

func (c *cssParser) AcceptRunWhitespaceNoNewlineComments() Comments {
	var cs Comments

	d := c.NewGoal()

Loop:
	for {
		switch d.AcceptRunWhitespaceNoNewLineNoComment() {
		case TokenComment:
		default:
			break Loop
		}

		cs = append(cs, d.Next())

		c.Score(d)

		d = c.NewGoal()

		d.AcceptRunWhitespaceNoNewLineNoComment()

		if d.Accept(TokenWhitespace) {
			if l := d.GetLastToken().Data; l != "\n" && l != "\r\n" {
				break
			}
		}
	}

	return cs
}

func (c *cssParser) GetLastToken() *Token {
	return &(*c)[len(*c)-1]
}

// Error is a parsing error with trace details.
type Error struct {
	Err     error
	Parsing string
	Token   Token
}

// Error returns the error string.
func (e Error) Error() string {
	return fmt.Sprintf("%s: error at position %d (%d:%d):\n%s", e.Parsing, e.Token.Pos+1, e.Token.Line+1, e.Token.LinePos+1, e.Err)
}

// Unwrap returns the wrapped error.
func (e Error) Unwrap() error {
	return e.Err
}

func (c *cssParser) Error(parsingFunc string, err error) error {
	tk := c.next()

	c.backup()

	return Error{
		Err:     err,
		Parsing: parsingFunc,
		Token:   *tk,
	}
}
