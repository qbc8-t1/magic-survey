package logger

import (
	"context"
	"time"

	"gorm.io/gorm/logger"
)

type GormLogger struct {
	appLogger *AppLogger
}

func NewGormLogger(appLogger *AppLogger) *GormLogger {
	return &GormLogger{appLogger: appLogger}
}

func (gl *GormLogger) LogMode(level logger.LogLevel) logger.Interface {
	return gl
}

func (gl *GormLogger) Info(ctx context.Context, msg string, args ...interface{}) {
	gl.appLogger.Infof(msg, args...)
}

func (gl *GormLogger) Warn(ctx context.Context, msg string, args ...interface{}) {
	gl.appLogger.Warnf(msg, args...)
}

func (gl *GormLogger) Error(ctx context.Context, msg string, args ...interface{}) {
	gl.appLogger.Errorf(msg, args...)
}

func (gl *GormLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	elapsed := time.Since(begin)
	sql, rows := fc()

	if err != nil {
		gl.appLogger.Errorf("SQL [%s] executed with error: %v (Rows affected: %d) in %s", sql, err, rows, elapsed)
	} else {
		gl.appLogger.Infof("SQL [%s] executed successfully (Rows affected: %d) in %s", sql, rows, elapsed)
	}
}
