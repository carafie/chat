package event

import (
	"encoding/json/v2"

	"github.com/carafie/chat/api/pkg/errcode"
)

var ErrInvalid = errcode.New("EVENT_INVALID", "event is invalid")

func Serialize(v any) ([]byte, error) {
	raw, err := json.Marshal(v)
	if err != nil {
		return nil, ErrInvalid
	}
	return raw, nil
}

func Deserialize[T any](raw []byte) (T, error) {
	var v T
	if err := json.Unmarshal(raw, &v, json.RejectUnknownMembers(true)); err != nil {
		return v, ErrInvalid
	}
	return v, nil
}
