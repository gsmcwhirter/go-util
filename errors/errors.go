package errors

import (
	"fmt"
	"strings"
)

type Error interface {
	error
	Cause() error
	Msg() string
	Data() []interface{}
}

type internalError interface {
	Error
	addDetails([]interface{})
}

type errStruct struct {
	msg  string
	data []interface{}
}

func (e *errStruct) Msg() string {
	return e.msg
}

func (e *errStruct) Error() string {
	if len(e.data) > 0 {
		return e.msg + " " + formatData(e.data)
	}

	return e.msg
}

func (e *errStruct) Cause() error {
	return nil
}

func (e *errStruct) Data() []interface{} {
	return e.data
}

func (e *errStruct) addDetails(data []interface{}) {
	e.data = append(e.data, data...)
}

type wrappedErr struct {
	errStruct
	cause error
}

func (e *wrappedErr) Error() string {
	var ret string

	if e.cause == nil {
		e.cause = New("(unknown error -- nil cause)")
	}

	data := e.Data()
	if len(data) > 0 {
		ret = e.cause.Error() + " " + formatData(data)
	} else {
		ret = e.cause.Error()
	}

	if e.msg == "" {
		return ret
	}

	return e.msg + ": " + ret
}

func (e *wrappedErr) Cause() error {
	return e.cause
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
		msg:  msg,
		data: nil,
	}
}

func Wrap(err error, msg string, data ...interface{}) error {
	if err == nil {
		return nil
	}

	if len(data)%2 != 0 {
		data = append(data, "")
	}

	return &wrappedErr{
		errStruct: errStruct{msg, data},
		cause:     err,
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
