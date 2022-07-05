package ginex

import (
	"fmt"
	"io"
	"net/http"
	"reflect"

	"github.com/gin-gonic/gin"
	"github.com/lethe3000/go-common/pkg/cecontext"
)

var (
	badRequestHandlerCode = 5001
)

func fail(c *gin.Context, msg string) {
	c.AbortWithStatusJSON(http.StatusOK, gin.H{
		"code": badRequestHandlerCode,
		"msg":  msg,
	})
}

// SessionInject 为handler函数自动注入session参数
// handler函数需要两个参数: *gin.Context, *cecontext.Context)
func SessionInject(handler interface{}) gin.HandlerFunc {
	return func(c *gin.Context) {
		var (
			handlerType = reflect.TypeOf(handler)
		)
		if handlerType.Kind() != reflect.Func {
			fail(c, "request handler参数不是函数类型")
			return
		}
		if handlerType.NumIn() != 2 {
			fail(c, "request handler必须为2个参数")
			return
		}

		var (
			ctx, ok = c.Get(Session)
			session *cecontext.Context
		)
		if ok {
			session = ctx.(*cecontext.Context)
		} else {
			session = cecontext.NewContext()
		}

		defer func() {
			if rec := recover(); rec != nil {
				fail(c, fmt.Sprintf("%v", rec))
			}
		}()

		handlerFunc := reflect.ValueOf(handler)
		params := []reflect.Value{reflect.ValueOf(c), reflect.ValueOf(session)}
		handlerFunc.Call(params)
	}
}

// RequestBodyBinding 针对post请求自动解析body结构体，并传递给handler函数
// value必须为指针类型
// handler函数需要接收3个参数: *gin.Context, *cecontext.Context, value类型的指针
// 具体使用参考request_test.go中的案例
func RequestBodyBinding(value interface{}, handler interface{}) gin.HandlerFunc {
	return func(c *gin.Context) {
		var (
			handlerType   = reflect.TypeOf(handler)
			valueType     = reflect.TypeOf(value)
			valueReceiver interface{}
		)

		if valueType.Kind() == reflect.Ptr {
			valueReceiver = reflect.New(valueType.Elem()).Interface()
		} else {
			fail(c, "post body参数不为指针类型")
			return
		}

		if handlerType.Kind() != reflect.Func {
			fail(c, "request handler参数不是函数类型")
			return
		}
		if handlerType.NumIn() != 3 {
			fail(c, "request handler必须为3个参数")
			return
		}

		var (
			ctx, ok = c.Get(Session)
			session *cecontext.Context
		)
		if ok {
			session = ctx.(*cecontext.Context)
		} else {
			session = cecontext.NewContext()
		}

		defer func() {
			if rec := recover(); rec != nil {
				fail(c, fmt.Sprintf("%v", rec))
			}
		}()

		handlerFunc := reflect.ValueOf(handler)
		if err := c.ShouldBindJSON(valueReceiver); err != nil {
			fail(c, err.Error())
			return
		}
		params := []reflect.Value{reflect.ValueOf(c), reflect.ValueOf(session), reflect.ValueOf(valueReceiver)}
		handlerFunc.Call(params)
	}
}

// RequestHandler function
type RequestHandler struct {
	Gin *gin.Engine
}

// NewRequestHandler creates a new request handler
func NewRequestHandler(logger io.Writer) *RequestHandler {
	gin.DefaultWriter = logger
	engine := gin.New()
	_ = engine.SetTrustedProxies([]string{"*"})
	return &RequestHandler{Gin: engine}
}
