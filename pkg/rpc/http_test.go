package rpc

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGet(t *testing.T) {
	var resp map[string]interface{}
	err := Get("http://inf-feishu-app-manager.nt.dev.fiture.com/api/v1/corps/1", nil, map[string]string{"secret": "1"}, &resp)
	assert.NotNil(t, err)
	assert.True(t, errors.Is(err, HttpStatusCodeBadErr))
}
