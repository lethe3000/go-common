package log

import (
	"testing"
	"time"

	"go.uber.org/zap"
)

func TestLogger(t *testing.T) {
	Logger().Infof("info: %s", "info message")
	Logger().Infow("infow:", "Url", "http://fiture.com", "retry", 3, "backoff", time.Second)
	Logger().Info("info:", zap.String("url", "http://fiture.com"))
}
