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

type Rule struct{}

func (r *Rule) parse(c *cssParser) error {
	return nil
}
