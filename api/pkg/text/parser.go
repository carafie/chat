package text

import (
	"errors"
	"strings"
	"unicode/utf8"
)

var (
	ErrTooShort = errors.New("text is too short")
	ErrTooLong  = errors.New("text is too long")
	ErrInvalid  = errors.New("text is invalid")
)

var newlineNormalizer = strings.NewReplacer("\r\n", "\n", "\r", "\n")

type Parser struct {
	MinChars  int
	MaxChars  int
	CharValid func(char rune) bool
}

func (p Parser) Parse(text string) (string, error) {
	if len(text) > p.MaxChars*utf8.UTFMax {
		return "", ErrTooLong
	}
	if p.CharValid == nil {
		p.CharValid = func(char rune) bool { return true }
	}

	text = newlineNormalizer.Replace(strings.TrimSpace(text))

	chars := 0
	for _, char := range text {
		if ok := p.CharValid(char); !ok {
			return "", ErrInvalid
		}
		chars++
	}
	if chars < p.MinChars {
		return "", ErrTooShort
	}
	if chars > p.MaxChars {
		return "", ErrTooLong
	}

	return text, nil
}
