package e

import (
	"fmt"
	"strconv"

	"github.com/pkg/errors"
)

var _ Codes = &E{}

// E is a Custom error message
// implement ecode.Codes
type E struct {
	code Code
	err  error // Custom errors
}

// New a E
func New(code Code, err error) Codes {
	return &E{code, errors.WithStack(err)}
}

// NewI a E
func NewI(code int, err error) Codes {
	return &E{Int(code), errors.WithStack(err)}
}

// Error implement error
func (e *E) Error() string {
	return fmt.Sprintf("Error: %s, Code: %d, Message: %s", e.err.Error(), e.Code(), e.Message())
}

// Code return error code
func (e *E) Code() int {
	return int(e.code.Code())
}

// Message return error message for developer
func (e *E) Message() string {
	if e.err != nil {
		return e.err.Error() // Prefer custom error messages
	}
	if e.code.Message() == "" {
		return strconv.Itoa(e.code.Code())
	}
	return e.code.Message()
}

// US .
func (e *E) US(msg string) Code {
	Register("en-us", Code(e.Code()), msg)
	return Code(e.Code())
}
