package logger

import (
	"backend-blog/internal/pkg"
	"context"
	"errors"
	"fmt"
	"log/slog"
	"runtime"
	"time"

	"gorm.io/gorm"
	glogger "gorm.io/gorm/logger"
)

type GormLogger struct {
	LogLevel                  glogger.LogLevel
	SlowThreshold             time.Duration
	IgnoreRecordNotFoundError bool
}

func NewGormLogger() *GormLogger {
	return &GormLogger{
		LogLevel:                  glogger.Info,
		SlowThreshold:             200 * time.Millisecond,
		IgnoreRecordNotFoundError: true,
	}
}

// LogMode 实现 glogger.Interface
func (l *GormLogger) LogMode(level glogger.LogLevel) glogger.Interface {
	newLogger := *l
	newLogger.LogLevel = level
	return &newLogger
}

// Info/Warn/Error 实现（对接 slog）
func (l *GormLogger) Info(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= glogger.Info {
		slog.InfoContext(ctx, fmt.Sprintf(msg, data...))
	}
}

func (l *GormLogger) Warn(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= glogger.Warn {
		slog.WarnContext(ctx, fmt.Sprintf(msg, data...))
	}
}

func (l *GormLogger) Error(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= glogger.Error {
		slog.ErrorContext(ctx, fmt.Sprintf(msg, data...))
	}
}

func (l *GormLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	if l.LogLevel <= glogger.Silent {
		return
	}

	elapsed := time.Since(begin)
	sql, rows := fc()

	// 1. 确定日志级别
	var level slog.Level
	msg := "GORM SQL"

	switch {
	case err != nil && l.LogLevel >= glogger.Error && !(errors.Is(err, gorm.ErrRecordNotFound) && l.IgnoreRecordNotFoundError):
		level = slog.LevelError
		msg = "GORM ERROR"
	case elapsed > l.SlowThreshold && l.SlowThreshold != 0 && l.LogLevel >= glogger.Warn:
		level = slog.LevelWarn
		msg = "GORM SLOW SQL"
	case l.LogLevel >= glogger.Info:
		level = slog.LevelInfo
	default:
		return
	}

	// 2. 获取调用栈 PC (跳过当前函数、GORM 内部函数)
	// 这里的 4 或者 5 通常能跳到你的 Service/DAO 层
	var pc uintptr
	var pcs [1]uintptr
	runtime.Callers(4, pcs[:])
	pc = pcs[0]
	// 3. 手动构建 Record 并处理
	r := slog.NewRecord(time.Now(), level, msg, pc)

	// 提取 TraceID：利用类型断言的 ok 模式，既安全又简洁
	if tid, ok := ctx.Value(pkg.TraceKey).(string); ok && tid != "" {
		r.AddAttrs(slog.String("trace_id", tid))
	}

	// 批量添加属性
	r.AddAttrs(
		slog.Duration("elapsed", elapsed),
		slog.Int64("rows", rows),
		slog.String("sql", sql),
	)
	if err != nil && level == slog.LevelError {
		r.AddAttrs(slog.Any("error", err))
	}

	// 使用默认的 logger 或你注入的实例
	_ = slog.Default().Handler().Handle(ctx, r)
}
