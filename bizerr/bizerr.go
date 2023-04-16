package bizerr

import (
	"errors"
	"fmt"
)

var _ error = &bizError{}

type bizError struct {
	Code    int32
	Message string
	cause   error
}

func (e *bizError) String() string {
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

func (e *bizError) Error() string {
	return e.String()
}

func (e *bizError) Unwrap() error {
	return e.cause
}

func New(code int32, format string, a ...interface{}) error {
	return &bizError{
		Code:    code,
		Message: fmt.Sprintf(format, a...),
	}
}

func NewWithCause(cause error, code int32, format string, a ...interface{}) error {
	return &bizError{
		Code:    code,
		Message: fmt.Sprintf(format, a...),
		cause:   cause,
	}
}

func IsBizError(err error) (*bizError, bool) {
	var bizErr *bizError
	if errors.As(err, &bizErr) {
		return bizErr, true
	}

	return nil, false
}

func UnwrapCodeNotZeroBizError(err error) (*bizError, bool) {
	bizErr, ok := IsBizError(err)
	if !ok {
		return nil, false
	}

	if bizErr.Code != 0 {
		return bizErr, true
	}

	return UnwrapCodeNotZeroBizError(bizErr.cause)
}

func IsCodeBizError(err error, code int32) bool {
	bizErr, ok := IsBizError(err)
	if !ok {
		return false
	}

	if bizErr.Code == code {
		return true
	}

	return IsCodeBizError(bizErr.cause, code)
}
