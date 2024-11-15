package dao

import (
	"context"
	"time"

	"gorm.io/gorm/logger"

	"github.com/kackerx/go-mall/common/log"
)

type GormLogger struct {
	SlowThreshold time.Duration
}

func (g *GormLogger) LogMode(level logger.LogLevel) logger.Interface {
	return &GormLogger{}
}

func (g *GormLogger) Info(ctx context.Context, s string, data ...interface{}) {
	log.New(ctx).Info(s, "data", data)
}

func (g *GormLogger) Warn(ctx context.Context, s string, data ...interface{}) {
	log.New(ctx).Warn(s, "data", data)
}

func (g *GormLogger) Error(ctx context.Context, s string, data ...interface{}) {
	log.New(ctx).Error(s, "data", data)
}

func (g *GormLogger) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
	dura := time.Since(begin).Milliseconds()

	// 错误日志
	sql, rows := fc()
	if err != nil {
		log.New(ctx).Error(
			"SQL ERROR",
			"sql", sql,
			"rows", rows,
			"dur(ms)", dura,
		)
	}

	// 慢日志
	if dura > g.SlowThreshold.Milliseconds() {
		log.New(ctx).Warn("SQL SLOW", "sql", sql, "rows", rows, "dur(ms)", dura)
	} else {
		log.New(ctx).Debug("SQL DEBUG", "sql", sql, "rows", rows, "dur(ms)", dura)
	}
}

func NewGormLogger(slowThreshold time.Duration) *GormLogger {
	return &GormLogger{SlowThreshold: slowThreshold}
}
