package css

import (
	"errors"
	"io"
	"strconv"
	"strings"

	"vimagination.zapto.org/parser"
)

// Unquote unquotes a CSS quoted string.
func Unquote(str string) (string, error) {
	if len(str) < 2 || str[0] != '\'' && str[0] != '"' || str[0] != str[len(str)-1] {
		return "", ErrBadString
	}

	tk := parser.NewStringTokeniser(str[1 : len(str)-1])

	var buf strings.Builder

	for {
		next := tk.ExceptRun("\"'\\")

		buf.WriteString(tk.Get())

		switch next {
		case -1:
			return buf.String(), nil
		case rune(str[0]):
			return "", ErrBadString
		case '\\':
			if err := unescape(&tk, &buf); err != nil {
				return "", err
			}
		default:
			tk.Next()
		}
	}
}

func unescape(tk *parser.Tokeniser, buf *strings.Builder) error {
	tk.Next()
	tk.Get()

	if tk.Accept(hexDigits) {
		for range 5 {
			tk.Accept(hexDigits)
		}

		unicode, _ := strconv.ParseUint(tk.Get(), 16, 32)

		buf.WriteRune(rune(unicode))

		tk.Accept(whitespace)
		tk.Get()
	} else {
		tk.Get()

		if tk.Next() == -1 {
			return io.ErrUnexpectedEOF
		}
	}

	return nil
}

// UnURL retrieves the escaped URL from a 'url(...)' string value."
func UnURL(str string) (string, error) {
	if len(str) < 5 || strings.ToLower(str[:4]) != "url(" || str[len(str)-1] != ')' {
		return "", ErrBadURL
	}

	tk := parser.NewStringTokeniser(str[4 : len(str)-1])

	tk.AcceptRun(whitespace)
	tk.Get()

	var buf strings.Builder

	for {
		next := tk.ExceptRun("\"'\\()" + whitespace)

		buf.WriteString(tk.Get())

		switch next {
		case ' ', '\t', '\r', '\n', '\f':
			if tk.AcceptRun(whitespace) != -1 {
				return "", ErrBadURL
			}

			fallthrough
		case -1:
			return buf.String(), nil
		case '\\':
			if err := unescape(&tk, &buf); err != nil {
				return "", err
			}
		default:
			return "", ErrBadURL
		}
	}
}

// Errors
var (
	ErrBadString = errors.New("bad string")
	ErrBadURL    = errors.New("bad url")
)
