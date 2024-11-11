package log

import (
	"context"
	"fmt"
	"path"
	"runtime"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Logger struct {
	ctx     context.Context
	traceID string
	spanID  string
	pspanID string
	_logger *zap.Logger
}

func New(ctx context.Context) *Logger {
	var (
		traceID, spanID, pspanID string
	)

	if ctx.Value("traceid") != nil {
		traceID = ctx.Value("traceid").(string)
	}
	if ctx.Value("spanid") != nil {
		spanID = ctx.Value("spanid").(string)
	}
	if ctx.Value("pspanid") != nil {
		pspanID = ctx.Value("pspanid").(string)
	}

	return &Logger{
		ctx:     ctx,
		traceID: traceID,
		spanID:  spanID,
		pspanID: pspanID,
		_logger: _logger,
	}
}

func (l *Logger) log(lvl zapcore.Level, msg string, kv ...any) {
	if len(kv)%2 != 0 {
		kv = append(kv, "UNKNOW")
	}

	// kv增加trace信息
	kv = append(kv, "traceid", l.traceID, "spanid", l.spanID, "pspanid", l.pspanID)
	// kv增加日志调用者信息
	funcName, fileName, line := l.getLoggerCallerInfo()
	kv = append(kv, "func", funcName, "file", fileName, "line", line)

	fields := make([]zap.Field, 0, len(kv)/2)
	for i := 0; i < len(kv); i += 2 {
		k := fmt.Sprintf("%v", kv[i])
		fields = append(fields, zap.Any(
			k,
			kv[i+1],
		))
	}

	ce := l._logger.Check(lvl, msg)
	ce.Write(fields...)
}

func (l *Logger) Debug(msg string, kv ...any) {
	l.log(zap.DebugLevel, msg, kv...)
}

func (l *Logger) Info(msg string, kv ...any) {
	l.log(zap.InfoLevel, msg, kv...)
}

func (l *Logger) Warn(msg string, kv ...any) {
	l.log(zap.WarnLevel, msg, kv...)
}

func (l *Logger) Error(msg string, kv ...any) {
	l.log(zap.ErrorLevel, msg, kv...)
}

// getLoggerCallerInfo 获取函数调用者的信息
func (l *Logger) getLoggerCallerInfo() (funcName, file string, line int) {
	pc, f, line, ok := runtime.Caller(3)
	if !ok {
		return
	}

	file = path.Base(f)
	funcName = runtime.FuncForPC(pc).Name()
	return
}
