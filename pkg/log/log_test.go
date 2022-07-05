package log

import (
	"testing"
)

type logconfig struct {
}

func (l logconfig) Debug() bool {
	return true
}

func (l logconfig) ServerName() string {
	return "go-common"
}

func TestLogger(t *testing.T) {
	c := logconfig{}
	logger := NewLogger(WithZap(c))(c)
	logger.Infof("info: %s", "info message")
	logger.Debugf("debug: %s", "debug")
}
