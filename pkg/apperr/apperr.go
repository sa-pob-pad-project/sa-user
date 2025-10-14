package apperr

import "errors"

type Code int

const (
	CodeBadRequest Code = iota + 1
	CodeUnauthorized
	CodeForbidden
	CodeNotFound
	CodeConflict
	CodeInternal
)

type Error struct {
	Code   Code
	Msg    string
	Err    error
	Fields map[string]any
}

func (e *Error) Error() string {
	if e.Err != nil {
		return e.Msg + ": " + e.Err.Error()
	}
	return e.Msg
}
func (e *Error) Unwrap() error { return e.Err }

func New(code Code, msg string, err error) *Error {
	return &Error{Code: code, Msg: msg, Err: err}
}

func IsCode(err error, code Code) bool {
	var ae *Error
	if errors.As(err, &ae) {
		return ae.Code == code
	}
	return false
}
