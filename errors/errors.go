package errors

import (
	"errors" //nolint:depguard // we need this for aliasing
	"fmt"
	"strings"
)

var (
	Is = errors.Is
	As = errors.As
)

// Error is our custom error type that allows wrapping errors with additional data.
// It should be go 1.13 compatible, implementing Unwrap (formerly Cause)
type Error interface {
	Error() string
	Unwrap() error
	Msg() string
	Data() []interface{}
}

type hasData interface {
	Data() []interface{}
}

type internalError interface {
	Error
	addDetails([]interface{})
}

type errStruct struct { //nolint:errname // disabled
	msg   string
	data  []interface{}
	cause error
}

func (e *errStruct) Msg() string {
	var msg string

	if e.msg != "" {
		msg = e.msg

		if e.cause != nil {
			msg += ": "
		}
	}

	if e2, ok := e.cause.(Error); ok {
		msg += e2.Msg()
	} else if e.cause != nil {
		msg += e.cause.Error()
	}

	return msg
}

func (e *errStruct) Error() string {
	var ret string

	if e.cause != nil {
		ret = e.cause.Error()

		if e.msg != "" {
			ret = ": " + ret
		}
	}

	data := e.data
	if len(data) > 0 {
		ret += " " + formatData(data)
	}

	if e.msg == "" {
		return ret
	}

	return e.msg + ret
}

func (e *errStruct) Unwrap() error {
	return e.cause
}

func (e *errStruct) Data() []interface{} {
	if d, ok := e.cause.(hasData); ok {
		subdata := d.Data()
		combined := make([]interface{}, 0, len(e.data)+len(subdata)+2)
		combined = append(combined, subdata...)
		combined = append(combined, e.data...)
		return combined
	}

	return e.data
}

func (e *errStruct) addDetails(data []interface{}) {
	e.data = append(e.data, data...)
}

func formatData(data []interface{}) string {
	kvs := make([]string, 0, len(data)/2)
	for i := 0; i < len(data)-1; i += 2 {
		kvs = append(kvs, fmt.Sprintf("%v=%v", data[i], data[i+1]))
	}

	return strings.Join(kvs, " ")
}

func New(msg string) error {
	return &errStruct{
		msg:   msg,
		data:  nil,
		cause: nil,
	}
}

func Newf(msg string, args ...interface{}) error {
	return &errStruct{
		msg:   fmt.Sprintf(msg, args...),
		data:  nil,
		cause: nil,
	}
}

func Wrap(err error, msg string, data ...interface{}) error {
	if err == nil {
		return nil
	}

	if len(data)%2 != 0 {
		data = append(data, "")
	}

	return &errStruct{
		msg:   msg,
		data:  data,
		cause: err,
	}
}

func WithDetails(err error, data ...interface{}) error {
	if err == nil {
		return nil
	}

	if len(data)%2 != 0 {
		data = append(data, "")
	}

	if e, ok := err.(internalError); ok {
		e.addDetails(data)
		return e
	}

	return Wrap(err, "", data...)
}
