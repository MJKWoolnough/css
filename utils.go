package css

import (
	"errors"
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
			tk.Next()
			tk.Get()

			if tk.Accept(hexDigits) {
				for range 5 {
					tk.Accept(hexDigits)
				}

				unicode, err := strconv.ParseUint(tk.Get(), 16, 32)
				if err != nil {
					return "", err
				}

				buf.WriteRune(rune(unicode))

				tk.Accept(whitespace)
				tk.Get()
			} else {
				tk.Get()
				tk.Next()
			}
		default:
			tk.Next()
		}
	}
}

var ErrBadString = errors.New("bad string")
