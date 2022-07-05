package ginex

import (
	"github.com/gin-gonic/gin"
)

var (
	Session                  = "session"
	DBTransactionRollbackKey = "db_trx_rollback"
	DBTransactionKey         = "db_trx"
)

type ApiError interface {
	Code() int
	Message() string
}

func Success(c *gin.Context, data interface{}) {
	c.JSON(200, gin.H{
		"code":    0,
		"content": data,
	})
}

func Fail(c *gin.Context, apiError ApiError, msg string) {
	FailWithStatus(c, apiError, msg, 200)
}

func FailWithStatus(c *gin.Context, apiError ApiError, msg string, code int) {
	c.Set(DBTransactionRollbackKey, true)
	c.AbortWithStatusJSON(code, gin.H{
		"code":   apiError.Code(),
		"msg":    apiError.Message(),
		"detail": msg,
	})
}

func ListContent(total int64, content interface{}) map[string]interface{} {
	return map[string]interface{}{
		"total": total,
		"list":  content,
	}
}
