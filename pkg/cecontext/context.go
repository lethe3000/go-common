package cecontext

import (
	"context"
	"fmt"

	"gorm.io/gorm"

	"github.com/google/uuid"
)

type Context struct {
	context.Context
	TraceId string
	User    interface{}
	DbTrx   *gorm.DB
}

func NewContext() *Context {
	return &Context{
		Context: context.TODO(),
		TraceId: uuid.New().String(),
	}
}

func (c *Context) SetUser(user interface{}) {
	c.User = user
}

func (c *Context) SetDbTrx(tr *gorm.DB) {
	c.DbTrx = tr
}

func (c *Context) String() string {
	return fmt.Sprintf("trace-id=%s User=%s", c.TraceId, c.User)
}
