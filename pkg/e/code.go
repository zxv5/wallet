package e

import (
	"fmt"
	"strconv"
	"sync/atomic"

	"github.com/pkg/errors"
)

var (
	_messages    atomic.Value // NOTE: stored map[string]map[int]string
	_codes       = map[int]struct{}{}
	_defaultI18n = "en-us"
)

// Register register ecode message map.
func Register(i18n string, e Code, msg string) {
	all, ok := _messages.Load().(map[string]map[int]string)
	if !ok {
		all = make(map[string]map[int]string)
	}

	if _, ok := all[i18n]; !ok {
		msgs := make(map[int]string)
		all[i18n] = msgs
	}
	all[i18n][e.Code()] = msg
	_messages.Store(all)
}

// SetDefault .
func SetDefault(i18n string) {
	_defaultI18n = i18n
}

// NewCode .
func NewCode(e int) Code {
	if e <= 0 {
		panic("e code must greater than zero")
	}
	return add(e)
}

func add(e int) Code {
	if _, ok := _codes[e]; ok {
		panic(fmt.Sprintf("ecode: %d already exist", e))
	}
	_codes[e] = struct{}{}
	return Int(e)
}

// Codes ecode error interface which has a code & message.
type Codes interface {
	// sometimes Error return Code in string form
	Error() string
	// Code get error code.
	Code() int
	// Message get code message.
	Message() string

	US(string) Code
}

// A Code is an int error code spec.
type Code int

func (e Code) Error() string {
	return strconv.FormatInt(int64(e), 10)
}

// Code return error code
func (e Code) Code() int { return int(e) }

// Message return error message
func (e Code) Message() string {
	if all, ok := _messages.Load().(map[string]map[int]string); ok {
		if msgs, ok := all[_defaultI18n]; ok {
			if msg, ok := msgs[e.Code()]; ok {
				return msg
			}
		}
	}

	return e.Error()
}

// Messagef .
func (e Code) Messagef(i18n string) string {
	if all, ok := _messages.Load().(map[string]map[int]string); ok {
		msgs, ok := all[i18n]
		if !ok {
			msgs, ok = all[_defaultI18n]
		}

		if ok {
			if msg, ok := msgs[e.Code()]; ok {
				return msg
			}
		}
	}

	return e.Error()
}

// Details return details.
func (e Code) Details() []interface{} { return nil }

// US .
func (e Code) US(msg string) Code {
	Register("en-us", e, msg)
	return e
}

// Int parse code int to error.
func Int(i int) Code { return Code(i) }

// String parse code string to error.
func String(e string) Code {
	if e == "" {
		return OK
	}
	// try error string
	i, err := strconv.Atoi(e)
	if err != nil {
		return ServerErr
	}
	return Code(i)
}

// Cause cause from error to ecode.
func Cause(e error) Codes {
	if e == nil {
		return OK
	}
	ec, ok := errors.Cause(e).(Codes)
	if ok {
		return ec
	}
	return String(e.Error())
}
