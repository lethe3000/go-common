package rpc

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	testUrl     = "https://jsonplaceholder.typicode.com/todos/1"
	postUrl     = "https://jsonplaceholder.typicode.com/posts"
	valid       = "2"
	validSecret = map[string]string{
		"secret": valid,
	}
)

func TestGet(t *testing.T) {
	t.Run("valid", func(t *testing.T) {
		var resp map[string]interface{}
		err := Get(testUrl, validSecret, validSecret, &resp)
		assert.Nil(t, err)
	})

	t.Run("v2", func(t *testing.T) {
		var result map[string]interface{}
		resp, _, err := GetV2(testUrl, nil, nil, nil, &result)
		assert.Nil(t, err)
		fmt.Println(resp)

		resp, _, err = PostV2(postUrl, map[string]interface{}{
			"title":  "foo",
			"body":   "bar",
			"userId": 1,
		}, nil, nil, nil, &result)
		assert.Nil(t, err)
		fmt.Println(resp)
	})
}

func TestPost(t *testing.T) {
	t.Run("error", func(t *testing.T) {
		var resp map[string]interface{}
		err := Post(testUrl, map[string]interface{}{}, nil, nil, &resp)
		assert.NotNil(t, err)
	})
}
