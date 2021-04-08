package rpc

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	testUrl     = "http://inf-feishu-app-manager.nt.dev.fiture.com/api/v1/corps/1"
	valid       = "Y6bKZsE1gVn6thYsEIDIerD-oYrNvCwFUGg0NSapy0M"
	invalid     = "1"
	validSecret = map[string]string{
		"secret": valid,
	}
	invalidSecret = map[string]string{
		"secret": invalid,
	}
)

func TestGet(t *testing.T) {
	t.Run("valid", func(t *testing.T) {
		var resp map[string]interface{}
		err := Get(testUrl, validSecret, validSecret, &resp)
		assert.Nil(t, err)
	})

	t.Run("error", func(t *testing.T) {
		var resp map[string]interface{}
		err := Get(testUrl, nil, invalidSecret, &resp)
		assert.NotNil(t, err)
		assert.True(t, errors.Is(err, HttpStatusCodeBadErr))
	})

}

func TestPost(t *testing.T) {
	t.Run("error", func(t *testing.T) {
		var resp map[string]interface{}
		err := Post(testUrl, map[string]interface{}{}, nil, nil, &resp)
		assert.NotNil(t, err)
	})
}
