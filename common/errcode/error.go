package errcode

import (
	"encoding/json"
	"fmt"
	"net/http"
	"path"
	"runtime"
)

type AppError struct {
	code     int    `json:"code"`
	msg      string `json:"msg"`
	cause    error  `json:"cause"`    // 保存根因, 如数据库错误这种底层错误, 生成项目的AppError或者预定义的
	occurred string `json:"occurred"` // 错误发生的位置
}

func (e *AppError) String() string {
	return e.Error()
}

func (e *AppError) Error() string {
	if e == nil {
		return ""
	}

	formattedErr := struct {
		Code     int    `json:"code"`
		Msg      string `json:"msg"`
		Cause    string `json:"cause"`
		Occurred string `json:"occurred"`
	}{
		Code:     e.code,
		Msg:      e.msg,
		Occurred: e.occurred,
	}

	if e.cause != nil {
		formattedErr.Cause = e.cause.Error()
	}

	errByte, _ := json.Marshal(formattedErr)
	return string(errByte)
}

func newError(code int, msg string) *AppError {
	return &AppError{
		code:  code,
		msg:   msg,
		cause: nil,
	}
}

func (e *AppError) Code() int {
	return e.code
}

func (e *AppError) Msg() string {
	return e.msg
}

func (e *AppError) HttpStatusCode() int {
	switch e.code {
	case Success.code:
		return http.StatusOK
	case ErrServer.code, ErrPanic.code:
		return http.StatusInternalServerError
	case ErrParams.code:
		return http.StatusBadRequest
	case ErrNotFound.code:
		return http.StatusNotFound
	case ErrTooManyRequests.code:
		return http.StatusTooManyRequests
	case ErrToken.code:
		return http.StatusUnauthorized
	case ErrForbidden.code:
		return http.StatusForbidden
	default:
		return http.StatusInternalServerError
	}
}

// WithCause AppError添加根因
func (e *AppError) WithCause(err error) *AppError {
	e.cause = err
	e.occurred = getAppErrOccurredInfo()
	return e
}

// Wrap 底层错误包装成应用层错误
func Wrap(msg string, err error) *AppError {
	if err == nil {
		return nil
	}

	return &AppError{-1, msg, err, getAppErrOccurredInfo()}
}

func getAppErrOccurredInfo() string {
	pc, file, line, ok := runtime.Caller(2)
	if !ok {
		return ""
	}

	file = path.Base(file)

	funcName := runtime.FuncForPC(pc).Name()
	return fmt.Sprintf("func: %s, file: %s, line: %d", funcName, file, line)
}
