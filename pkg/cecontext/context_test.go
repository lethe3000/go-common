package cecontext

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCtx(t *testing.T) {
	ctx := NewContext()
	assert.NotEmpty(t, ctx.TraceId)
}
