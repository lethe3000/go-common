package ginex

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/lethe3000/go-common/pkg/cecontext"
	"github.com/stretchr/testify/assert"
)

func TestApi(t *testing.T) {
	createRequest := func(body []byte) *gin.Context {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		var requestBody io.Reader
		marsh := []byte(`{"foo": "f", "bar": 100}`)
		requestBody = bytes.NewReader(marsh)
		req, _ := http.NewRequest("GET", "", requestBody)
		c.Request = req
		return c
	}

	c := createRequest([]byte(`{"foo": "f", "bar": 100}`))

	echoHandler := func(c *gin.Context, session *cecontext.Context, body *testReq) {
		assert.Equal(t, body.Foo, "f")
		assert.Equal(t, body.Bar, 100)
	}
	RequestBodyBinding(&testReq{}, echoHandler)(c)
	assert.False(t, c.IsAborted())

	c = createRequest([]byte(`{"foo": "f", "bar": 100}`))

	RequestBodyBinding(testReq{}, echoHandler)(c)
	assert.True(t, c.IsAborted())

	invalidHandler := func(c *gin.Context) {}
	RequestBodyBinding(&testReq{}, invalidHandler)(c)
	assert.True(t, c.IsAborted())
}

type testReq struct {
	Foo string `json:"foo"`
	Bar int    `json:"bar"`
}
