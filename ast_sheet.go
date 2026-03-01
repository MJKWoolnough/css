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

type Sheet struct{}

func (s *Sheet) parse(c *cssParser) error {
	return nil
}
