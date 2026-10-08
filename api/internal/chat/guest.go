package chat

import (
	"unicode"
	"uuid"

	"github.com/carafie/chat/api/pkg/errcode"
	"github.com/carafie/chat/api/pkg/text"
)

var (
	ErrGuestNameTooShort = errcode.New("GUEST_NAME_TOO_SHORT", "guest name is too short")
	ErrGuestNameTooLong  = errcode.New("GUEST_NAME_TOO_LONG", "guest name is too long")
	ErrGuestNameInvalid  = errcode.New("GUEST_NAME_INVALID", "guest name is invalid")
)

type Guest struct {
	ID    uuid.UUID
	Name  string
	Color Color
}

func NewGuest(id uuid.UUID, name, color string) (Guest, error) {
	name, err := text.Parser{
		MinChars: 2,
		MaxChars: 20,
		CharValid: func(char rune) bool {
			return unicode.IsPrint(char) && !unicode.IsSpace(char)
		},
	}.Parse(name)
	switch err {
	case nil:
		return Guest{ID: id, Name: name, Color: NewColor(color)}, nil
	case text.ErrTooShort:
		return Guest{}, ErrGuestNameTooShort
	case text.ErrTooLong:
		return Guest{}, ErrGuestNameTooLong
	default:
		return Guest{}, ErrGuestNameInvalid
	}
}

type Color = string

const (
	ColorRed    Color = "red"
	ColorOrange Color = "orange"
	ColorGreen  Color = "green"
	ColorTeal   Color = "teal"
	ColorBlue   Color = "blue"
	ColorPurple Color = "purple"
	ColorPink   Color = "pink"
)

func NewColor(color string) Color {
	switch color {
	case ColorRed, ColorOrange, ColorGreen, ColorTeal, ColorBlue, ColorPurple, ColorPink:
		return color
	default:
		return ColorRed
	}
}
