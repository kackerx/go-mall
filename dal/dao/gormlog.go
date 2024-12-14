package dao

import (
	"context"
	"time"

	"gorm.io/gorm/logger"
)

type GormLogger struct {
	SlowThreshold time.Duration
}

func (g *GormLogger) LogMode(level logger.LogLevel) logger.Interface {
	return &GormLogger{}
}

func (g *GormLogger) Info(ctx context.Context, s string, data ...interface{}) {
	logger.New(ctx).Info(s, "data", data)
}

func (g *GormLogger) Warn(ctx context.Context, s string, data ...interface{}) {
	logger.New(ctx).Warn(s, "data", data)
}

func (g *GormLogger) Error(ctx context.Context, s string, data ...interface{}) {
	logger.New(ctx).Error(s, "data", data)
}

func (g *GormLogger) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
	dura := time.Since(begin).Milliseconds()

	// 错误日志
	sql, rows := fc()
	if err != nil {
		logger.New(ctx).Error(
			"SQL ERROR",
			"sql", sql,
			"rows", rows,
			"dur(ms)", dura,
		)
	}

	// 慢日志
	if dura > g.SlowThreshold.Milliseconds() {
		logger.New(ctx).Warn("SQL SLOW", "sql", sql, "rows", rows, "dur(ms)", dura)
	} else {
		logger.New(ctx).Debug("SQL DEBUG", "sql", sql, "rows", rows, "dur(ms)", dura)
	}
}

func NewGormLogger(slowThreshold time.Duration) *GormLogger {
	return &GormLogger{SlowThreshold: slowThreshold}
}
