package errors

import (
	"errors"
	"fmt"
	"github.com/SimonMorphy/grok-go/constants"
)

type Error struct {
	code int
	msg  string
	err  error
}

func (e *Error) Error() string {
	var msg string
	if e.msg != "" {
		msg = e.msg
	}
	msg = constants.ErrMsg[e.code]
	return msg + " -> " + e.err.Error()
}

func New(code int) error {
	return &Error{
		code: code,
		err:  errors.New(constants.ErrMsg[code]),
	}
}

func NewWithError(code int, err error) error {
	if err == nil {
		return New(code)
	}
	return &Error{
		code: code,
		err:  err,
	}
}

func NewWithMsgF(code int, format string, args ...any) error {
	return &Error{
		code: code,
		msg:  fmt.Sprintf(format, args...),
	}
}

func Errno(err error) int {
	if err == nil {
		return constants.ErrnoSuccess
	}
	targetError := &Error{}
	if errors.As(err, &targetError) {
		return targetError.code
	}
	return -1
}

func Output(err error) (int, string) {
	if err == nil {
		return constants.ErrnoSuccess, constants.ErrMsg[constants.ErrnoSuccess]
	}
	errno := Errno(err)
	if errno == -1 {
		return constants.ErrnoUnknownError, err.Error()
	}
	return errno, err.Error()
}
