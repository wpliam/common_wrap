package errs

import "fmt"

const (
	RetOk           = 0
	RetInvalidParam = 101
	RetUnknown      = 999
)

type Error struct {
	Code int32
	Msg  string
}

func New(code int32, msg string) error {
	return &Error{
		Code: code,
		Msg:  msg,
	}
}

func Newf(code int32, format string, args ...interface{}) error {
	return &Error{
		Code: code,
		Msg:  fmt.Sprintf(format, args...),
	}
}

func Code(e error) int32 {
	if e == nil {
		return RetOk
	}
	err, ok := e.(*Error)
	if !ok {
		return RetUnknown
	}
	if err == (*Error)(nil) {
		return RetOk
	}
	return err.Code
}

func Msg(e error) string {
	if e == nil {
		return ""
	}
	err, ok := e.(*Error)
	if !ok {
		return e.Error()
	}
	if err == (*Error)(nil) {
		return ""
	}
	return err.Msg
}

func (e *Error) Error() string {
	if e == nil {
		return ""
	}
	return fmt.Sprintf("code:%d msg:%s", e.Code, e.Msg)
}
