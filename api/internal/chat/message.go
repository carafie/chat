package chat

import (
	"unicode"

	"github.com/carafie/chat/api/pkg/errcode"
	"github.com/carafie/chat/api/pkg/text"
)

var (
	ErrMessageTooShort = errcode.New("MESSAGE_TOO_SHORT", "message is too short")
	ErrMessageTooLong  = errcode.New("MESSAGE_TOO_LONG", "message is too long")
	ErrMessageInvalid  = errcode.New("MESSAGE_INVALID", "message is invalid")
)

type Message = string

func NewMessage(message string) (Message, error) {
	message, err := text.Parser{
		MinChars: 1,
		MaxChars: 500,
		CharValid: func(char rune) bool {
			return unicode.IsPrint(char)
		},
	}.Parse(message)
	switch err {
	case nil:
		return Message(message), nil
	case text.ErrTooShort:
		return "", ErrMessageTooShort
	case text.ErrTooLong:
		return "", ErrMessageTooLong
	default:
		return "", ErrMessageInvalid
	}
}
