package status

import (
	"errors"
	"fmt"
)

var _ error = &Error{}

type Error struct {
	Status  int
	Message string
	cause   error
}

func (e *Error) String() string {
	if e == nil {
		return ""
	}

	if e.cause == nil {
		return e.Message
	}

	if e.Message == "" {
		return fmt.Sprintf("%v", e.cause)
	}

	return fmt.Sprintf("%s, err: %v", e.Message, e.cause)
}

func (e *Error) Error() string {
	return e.String()
}

func (e *Error) Unwrap() error {
	return e.cause
}

func NewError(status int, format string, a ...interface{}) error {
	return &Error{
		Status:  status,
		Message: fmt.Sprintf(format, a...),
	}
}

func NewErrorWithCause(cause error, status int, format string, a ...interface{}) error {
	return &Error{
		Status:  status,
		Message: fmt.Sprintf(format, a...),
		cause:   cause,
	}
}

func IsStatusError(err error) (*Error, bool) {
	var statusErr *Error
	if errors.As(err, &statusErr) {
		return statusErr, true
	}

	return nil, false
}

func UnwrapStatusNotZeroError(err error) (*Error, bool) {
	statusErr, ok := IsStatusError(err)
	if !ok {
		return nil, false
	}

	if statusErr.Status != 0 {
		return statusErr, true
	}

	return UnwrapStatusNotZeroError(statusErr.cause)
}
