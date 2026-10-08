package errcode

import "errors"

// Error extends standard library errors with custom error codes.
type Error struct {
	Code    string
	Message string
}

func New(code, message string) error {
	return Error{
		Code:    code,
		Message: message,
	}
}

func (e Error) Error() string {
	return e.Message
}

// Code returns the code of the error, or an empty string if none found.
func Code(err error) string {
	if err, ok := errors.AsType[Error](err); ok {
		return err.Code
	}
	return ""
}
