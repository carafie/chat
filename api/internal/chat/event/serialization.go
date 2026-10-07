package event

import (
	"encoding/json/v2"
	"errors"
)

var ErrInvalid = errors.New("event is invalid")

func Serialize(v any) ([]byte, error) {
	raw, err := json.Marshal(v)
	if err != nil {
		return nil, ErrInvalid
	}
	return raw, nil
}

func Deserialize[T any](raw []byte) (T, error) {
	var v T
	if err := json.Unmarshal(raw, &v); err != nil {
		return v, ErrInvalid
	}
	return v, nil
}
