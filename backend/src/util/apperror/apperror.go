package apperror

import (
	"errors"
	"fmt"

	"github.com/DaisukeMatsumoto0925/backend/src/util/errorcode"
	"golang.org/x/xerrors"
)

type AppError interface {
	Error() string
	Code() errorcode.ErrorCode
	SetCode(code errorcode.ErrorCode) AppError
	Info(infoMessage string) AppError
	InfoMessage() string
}

// カスタムエラー型 Error() string を持つ
type appError struct {
	err         error
	message     string
	frame       xerrors.Frame
	errCode     errorcode.ErrorCode
	infoMessage string
}

func create(msg string) *appError {
	var e appError
	e.message = msg
	e.frame = xerrors.Caller(2)
	return &e
}

func New(msg string) AppError {
	return create(msg)
}

func Errorf(format string, a ...interface{}) *appError {
	return create(fmt.Sprintf(format, a...))
}

func Wrap(err error, msg ...string) *appError {
	var m string
	if len(msg) != 0 {
		m = msg[0]
	}
	e := create(m)
	e.err = err
	return e
}

func Wrapf(err error, format string, args ...interface{}) *appError {
	e := create(fmt.Sprintf(format, args...))
	e.err = err
	return e
}

func Unwrap(e error) error {
	var appErr *appError
	if errors.As(e, &appErr) {
		return appErr.err
	}

	return e
}

func AsAppError(err error) *appError {
	if err == nil {
		return nil
	}

	var e *appError
	if errors.As(err, &e) {
		return e
	}
	return nil
}

func (err *appError) SetCode(code errorcode.ErrorCode) AppError {
	err.errCode = code
	return err
}

func (err *appError) Info(infoMessage string) AppError {
	err.infoMessage = infoMessage
	return err
}

func (err *appError) Infof(format string, a ...interface{}) AppError {
	err.infoMessage = fmt.Sprintf(format, a...)
	return err
}

func (err *appError) Error() string {
	if err.err == nil {
		return err.message
	}
	if err.message != "" {
		return err.message + ": " + err.err.Error()
	}
	return err.err.Error()
}

func (err *appError) Unwrap() error {
	return err.err
}

func (err *appError) Code() errorcode.ErrorCode {
	var next *appError = err
	for next.errCode == "" {
		if err := AsAppError(next.err); err != nil {
			next = err
		} else {
			return errorcode.Unknown
		}
	}
	return next.errCode
}

func (err *appError) InfoMessage() string {
	var next *appError = err
	for next.infoMessage == "" {
		if err := AsAppError(next.err); err != nil {
			next = err
		} else {
			return ""
		}
	}
	return next.infoMessage
}

func (err *appError) Format(f fmt.State, c rune) {
	xerrors.FormatError(err, f, c)
}

func (err *appError) FormatError(p xerrors.Printer) error {
	p.Print(err.Error())
	if p.Detail() {
		err.frame.Format(p)
	}
	return err.err
}
