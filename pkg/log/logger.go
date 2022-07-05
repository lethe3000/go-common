package log

import (
	"fmt"
	"time"

	"github.com/lethe3000/go-common/pkg/cecontext"

	"context"

	"go.uber.org/zap/zapcore"

	"go.uber.org/fx/fxevent"

	"go.uber.org/zap"

	gormlogger "gorm.io/gorm/logger"
)

type LoggerConfig interface {
	Debug() bool
	ServerName() string
}

type BasicLogger interface {
	Infof(msg string, args ...interface{})
	Debugf(msg string, args ...interface{})
	Warnf(msg string, args ...interface{})
	Errorf(msg string, args ...interface{})
}

// Logger structure
type Logger struct {
	loggers []BasicLogger
}

func (l *Logger) Infof(msg string, args ...interface{}) {
	for _, logger := range l.loggers {
		if logger != nil {
			logger.Infof(msg, args...)
		}
	}
}

func (l *Logger) Debugf(msg string, args ...interface{}) {
	for _, logger := range l.loggers {
		if logger != nil {
			logger.Debugf(msg, args...)
		}
	}
}

func (l *Logger) Warnf(msg string, args ...interface{}) {
	for _, logger := range l.loggers {
		if logger != nil {
			logger.Warnf(msg, args...)
		}
	}
}

func (l *Logger) Errorf(msg string, args ...interface{}) {
	for _, logger := range l.loggers {
		if logger != nil {
			logger.Errorf(msg, args...)
		}
	}
}

//type GlogWrapper struct {
//	*zap.SugaredLogger
//}
//
//func (g GlogWrapper) Infof(msg string, args ...interface{}) {
//	if common.Glog != nil {
//		common.Glog.Info(msg, args...)
//	}
//	g.SugaredLogger.Infof(msg, args...)
//}
//
//func (g GlogWrapper) Debugf(msg string, args ...interface{}) {
//	if common.Glog != nil {
//		common.Glog.Debug(msg, args...)
//	}
//	var topFrame Frame = 0
//	stack := Callers(3)
//	for _, f := range *stack {
//		topFrame = Frame(f)
//		break
//	}
//	g.SugaredLogger.Debugf(fmt.Sprintf("%s %s", topFrame, msg), args...)
//}
//
//func (g GlogWrapper) Warnf(msg string, args ...interface{}) {
//	if common.Glog != nil {
//		common.Glog.Warn(msg, args...)
//	}
//	g.SugaredLogger.Warnf(msg, args...)
//}
//
//func (g GlogWrapper) Errorf(msg string, args ...interface{}) {
//	if common.Glog != nil {
//		common.Glog.Error(msg, args...)
//	}
//	g.SugaredLogger.Errorf(msg, args...)
//}

type GinLogger struct {
	*Logger
}

type FxLogger struct {
	*Logger
}

type GormLogger struct {
	*Logger
	gormlogger.Config
}

func NewLogger(loggers ...BasicLogger) func(config LoggerConfig) *Logger {
	return func(config LoggerConfig) *Logger {
		return &Logger{loggers: loggers}
	}
}

// NewLogger get the logger
//func NewLogger(config LoggerConfig) *Logger {
//	if os.Getenv("IS_MOBIUS") == "true" || os.Getenv("MOBIUS_LOCAL") == "1" {
//		return &Logger{WithGlog(config)}
//	}
//	return &Logger{WithZap(config)}
//}

// GetGinLogger get the gin logger
func (l *Logger) GetGinLogger() GinLogger {
	return GinLogger{
		Logger: l,
	}
}

// GetFxLogger get the fx logger
func (l *Logger) GetFxLogger() *FxLogger {
	return &FxLogger{
		Logger: l,
	}
}

// GetGormLogger gets the gorm framework logger
func (l *Logger) GetGormLogger() *GormLogger {
	return &GormLogger{
		Logger: l,
		Config: gormlogger.Config{
			LogLevel: gormlogger.Info,
		},
	}
}

func WithZap(c LoggerConfig) BasicLogger {
	config := zap.NewDevelopmentConfig()
	if c.Debug() {
		config.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
	} else {
		config.Level = zap.NewAtomicLevelAt(zap.InfoLevel)
	}
	zapLogger, _ := config.Build()

	return zapLogger.
		WithOptions(
			zap.WithCaller(true),
			zap.AddStacktrace(zapcore.WarnLevel),
		).
		Named(c.ServerName()).
		Sugar()
}

// Write interface implementation for gin-framework
func (l GinLogger) Write(p []byte) (n int, err error) {
	l.Infof(string(p))
	return len(p), nil
}

// Printf prints go-fx logs
func (l FxLogger) Printf(str string, args ...interface{}) {
	l.Infof(str, args)
}

func (l FxLogger) LogEvent(event fxevent.Event) {
	l.Debugf("fx event: %s", event)
}

// GORM Framework logger Interfac e Implementations
// ---- START ----

// LogMode set log mode
func (l *GormLogger) LogMode(level gormlogger.LogLevel) gormlogger.Interface {
	newlogger := *l
	newlogger.LogLevel = level
	return &newlogger
}

// Info prints info
func (l GormLogger) Info(ctx context.Context, str string, args ...interface{}) {
	if l.LogLevel >= gormlogger.Info {
		l.Infof(str, args...)
	}
}

// Warn prints warn messages
func (l GormLogger) Warn(ctx context.Context, str string, args ...interface{}) {
	if l.LogLevel >= gormlogger.Warn {
		l.Warnf(str, args...)
	}
}

// Error prints error messages
func (l GormLogger) Error(ctx context.Context, str string, args ...interface{}) {
	if l.LogLevel >= gormlogger.Error {
		l.Errorf(str, args...)
	}
}

// Trace prints trace messages
func (l GormLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	if l.LogLevel <= 0 {
		return
	}
	elapsed := time.Since(begin)
	if l.LogLevel >= gormlogger.Info {
		sql, rows := fc()
		l.Debugf("[%d ms, %d rows] sql -> %s", elapsed.Milliseconds(), rows, sql)
		return
	}

	if l.LogLevel >= gormlogger.Warn {
		sql, rows := fc()
		l.Warnf("[%d ms, %d rows] sql -> %s", elapsed.Milliseconds(), rows, sql)
		return
	}

	if l.LogLevel >= gormlogger.Error {
		sql, rows := fc()
		l.Errorf("[ %d ms, %d rows] sql -> %s", elapsed.Milliseconds(), rows, sql)
		return
	}
}

type logfunc func(string, ...interface{})

type logMethod struct {
	Errorf,
	Warnf,
	Infof,
	Debugf logfunc
}

func (l Logger) WithContext(ctx context.Context) *logMethod {
	if ctx == nil {
		ctx = cecontext.NewContext()
	}
	// 将上下文相关信息打印在最前面
	contextMsg := fmt.Sprintf("%s", ctx)
	return &logMethod{
		Debugf: func(format string, args ...interface{}) {
			format = contextMsg + format
			l.Debugf(format, args...)
		},
		Infof: func(format string, args ...interface{}) {
			format = contextMsg + format
			l.Infof(format, args...)
		},
		Warnf: func(format string, args ...interface{}) {
			format = contextMsg + format
			l.Warnf(format, args...)
		},
		Errorf: func(format string, args ...interface{}) {
			format = contextMsg + format
			l.Errorf(format, args...)
		},
	}
}
