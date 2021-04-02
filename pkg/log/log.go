package log

import (
	"go.uber.org/zap"
	"sync"
)

var logger_ *zap.SugaredLogger
var once sync.Once

func initLogger() {
	once.Do(func() {
		logger, err := zap.NewProduction()
		if err != nil {
			panic(err)
		}
		logger_ = logger.Sugar()
	})
}

func Logger() *zap.SugaredLogger {
	if logger_ == nil {
		panic("log instance not initialized")
	}
	return logger_
}