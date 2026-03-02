package css

import "vimagination.zapto.org/parser"

func ParseSheet(t parser.Tokeniser) (*Sheet, error) {
	c, err := newCSSParser(createTokeniser(t))
	if err != nil {
		return nil, err
	}

	s := new(Sheet)
	if err := s.parse(&c); err != nil {
		return nil, err
	}

	return s, nil
}

type Sheet struct {
	Rules  []Rule
	Tokens Tokens
}

func (s *Sheet) parse(c *cssParser) error {
	for c.AcceptRunWhitespace() != parser.TokenDone {
		d := c.NewGoal()
		var r Rule

		if err := r.parse(&d); err != nil {
			return c.Error("Sheet", err)
		}

		s.Rules = append(s.Rules, r)

		c.Score(d)
	}

	s.Tokens = c.ToTokens()

	return nil
}

type Rule struct {
	CommentDelimiter *Token
	AtRule           *AtRule
	QualifiedRule    *QualifiedRule
	Tokens           Tokens
}

func (r *Rule) parse(c *cssParser) error {
	if c.Accept(TokenCDO, TokenCDC) {
		r.CommentDelimiter = c.GetLastToken()
	} else if tk := c.Peek(); tk.Type == TokenAtKeyword {
		d := c.NewGoal()
		r.AtRule = new(AtRule)

		if err := r.AtRule.parse(&d); err != nil {
			return c.Error("Rule", err)
		}

		c.Score(d)
	} else {
		d := c.NewGoal()
		r.QualifiedRule = new(QualifiedRule)

		if err := r.QualifiedRule.parse(&d); err != nil {
			return c.Error("Rule", err)
		}

		c.Score(d)
	}

	r.Tokens = c.ToTokens()

	return nil
}

type AtRule struct{}

func (a *AtRule) parse(c *cssParser) error {
	return nil
}

type QualifiedRule struct{}

func (q *QualifiedRule) parse(c *cssParser) error {
	return nil
}
