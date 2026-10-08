package chat

import (
	"unicode"

	"github.com/carafie/chat/api/pkg/errcode"
	"github.com/carafie/chat/api/pkg/text"
)

var (
	ErrRoomNameTooShort = errcode.New("ROOM_NAME_TOO_SHORT", "room name is too short")
	ErrRoomNameTooLong  = errcode.New("ROOM_NAME_TOO_LONG", "room name is too long")
	ErrRoomNameInvalid  = errcode.New("ROOM_NAME_INVALID", "room name is invalid")
)

type Room struct {
	Name string
}

func NewRoom(name string) (Room, error) {
	name, err := text.Parser{
		MinChars: 2,
		MaxChars: 20,
		CharValid: func(char rune) bool {
			return unicode.IsPrint(char) && !unicode.IsSpace(char)
		},
	}.Parse(name)
	switch err {
	case nil:
		return Room{Name: name}, nil
	case text.ErrTooShort:
		return Room{}, ErrRoomNameTooShort
	case text.ErrTooLong:
		return Room{}, ErrRoomNameTooLong
	default:
		return Room{}, ErrRoomNameInvalid
	}
}
