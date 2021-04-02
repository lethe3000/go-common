package log

import (
	"go.uber.org/zap"
	"testing"
	"time"
)

func TestLogger(t *testing.T) {
	initLogger()
	Logger().Infof("info: %s", "info message")
	Logger().Infow("infow:", "Url", "http://fiture.com", "retry", 3, "backoff", time.Second)
	Logger().Info("info:", zap.String("url", "http://fiture.com"))
}
