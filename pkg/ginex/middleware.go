package ginex

import (
	"github.com/gin-gonic/gin"
	"github.com/lethe3000/go-common/pkg/cecontext"
	"github.com/lethe3000/go-common/pkg/db"
)

var corsHeaders = map[string]string{
	"Access-Control-Allow-Origin":      "*",
	"Access-Control-Allow-Credentials": "true",
	"Access-Control-Allow-Methods":     "POST, GET, OPTIONS, PUT, DELETE, UPDATE",
	"Access-Control-Allow-Headers":     "Authorization, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With",
}

func Cors(c *gin.Context) {
	for h, v := range corsHeaders {
		c.Writer.Header().Set(h, v)
	}

	if c.Request.Method == "OPTIONS" {
		c.AbortWithStatus(204)
		return
	}
	c.Next()
}

func Trace(c *gin.Context) {
	ctx := cecontext.NewContext()
	c.Set(Session, ctx)

	c.Next()
}

type GinLogger interface {
	Debugf(f string, args ...interface{})
	Warnf(f string, args ...interface{})
}

type DatabaseTrxMiddleware struct {
	handler *RequestHandler
	logger  GinLogger
	db      db.Database
}

func (e *DatabaseTrxMiddleware) Serve(c *gin.Context) {
	txHandle := e.db.DB.Begin()
	e.logger.Debugf("beginning database transaction: %v", c.Request.URL.Path)
	defer func() {
		// 这里捕获到内部错误，需要回滚并告知前端
		if r := recover(); r != nil {
			e.logger.Warnf("database transaction panic: %v %v", r, txHandle.Error)
			txHandle.Rollback()
			Fail(c, ErrUncaughtErr, errorToString(r))
		}
	}()
	c.Set(DBTransactionKey, txHandle)
	c.Next()

	rollback := c.GetBool(DBTransactionRollbackKey)
	if rollback {
		e.logger.Debugf("rolling back database transaction: %v", c.Request.URL.Path)
		if err := txHandle.Rollback().Error; err != nil {
			e.logger.Debugf("rolling back database transaction err: %v", err)
		}
	} else {
		e.logger.Debugf("committing database transaction: %v", c.Request.URL.Path)
		if err := txHandle.Commit().Error; err != nil {
			e.logger.Debugf("committing database transaction error: %v", err)
		}
	}
}

func errorToString(e interface{}) string {
	switch v := e.(type) {
	case error:
		return v.Error()
	default:
		return e.(string)
	}
}
