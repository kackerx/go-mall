package errcode

import (
	"encoding/json"
	"fmt"
	"net/http"
	"path"
	"runtime"
)

/*
*
* 第三方组件不确定是什么err, 用Wrap包装一下且返回500
* 直到是什么类型的错误的话用预定义错误且追加根因方式, WithCause
*
* Handle Errors Once
* 底层(Dao, 基础设施层): 抛出错误
* 中层(领域层, 应用层): 包装错误
* 上层(控制器接口层): 日志记录错误, errors.Is区分错误
*
 */

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
	if code > -1 {
		if _, ok := codes[code]; ok {
			panic(fmt.Sprintf("预定义错误码不能重复: %d", code))
		}
	}

	codes[code] = struct{}{}
	return &AppError{code: code, msg: msg, cause: nil}
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
	newErr := e.Clone()
	newErr.cause = err
	newErr.occurred = getAppErrOccurredInfo()
	return newErr
}

func (e *AppError) Clone() *AppError {
	return &AppError{e.code, e.msg, e.cause, e.occurred}
}

// Wrap 底层错误包装成应用层错误, 当调用第三方组件不确定是什么err时Wrap一下, 返回500
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

func (e *AppError) UnWrap() error {
	return e.cause
}

func (e *AppError) Is(target error) bool {
	targerErr, ok := target.(*AppError)
	if !ok {
		return false
	}

	return targerErr.code == e.code
}
