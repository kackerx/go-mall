package logger

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"gopkg.in/natefinch/lumberjack.v2"

	"github.com/kackerx/go-mall/common/enum"
	"github.com/kackerx/go-mall/config"
)

var _logger *zap.Logger

func init() {
	encoderConfg := zap.NewProductionEncoderConfig()
	encoderConfg.EncodeTime = zapcore.ISO8601TimeEncoder
	encoder := zapcore.NewJSONEncoder(encoderConfg)

	fileWriteSyncer := getFileLogWriter()

	var cores []zapcore.Core
	switch config.Conf.App.Env {
	// 测试环境和生产环境 --> 文件, info
	case enum.ModeTest, enum.ModeProd:
		cores = append(cores, zapcore.NewCore(encoder, fileWriteSyncer, zapcore.InfoLevel))
	// 开发环境 --> 控制台 & 文件, debug
	case enum.ModeDev:
		cores = append(
			cores,
			zapcore.NewCore(encoder, fileWriteSyncer, zap.DebugLevel),
			zapcore.NewCore(encoder, zapcore.AddSync(os.Stdout), zap.DebugLevel),
		)
	}

	core := zapcore.NewTee(cores...)
	_logger = zap.New(core)
}

func getFileLogWriter() zapcore.WriteSyncer {
	lumberJackLogger := &lumberjack.Logger{
		Filename:   config.Conf.App.Log.Path,
		MaxSize:    config.Conf.App.Log.MaxSize,
		MaxAge:     config.Conf.App.Log.MaxAge,
		MaxBackups: 0,
		LocalTime:  true,
		Compress:   false,
	}
	return zapcore.AddSync(lumberJackLogger)
}

func ZapTest() {
	_logger.Info(
		"hehe",
		zap.Any("app", config.Conf.App),
		zap.Any("db", config.Conf.DB),
		zap.Any("data", "呵呵呵呵呵呵你是个 人才\n挺厉害的\n"),
	)
}
