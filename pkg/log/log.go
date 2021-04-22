package log

import (
	"sync"

	"go.uber.org/zap"
)

var logger_ *zap.SugaredLogger
var once sync.Once

func init() {
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
