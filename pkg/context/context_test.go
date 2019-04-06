package context

import (
	"github.com/devster/tarreleaser/pkg/config"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
	"time"
)

func TestNewWithTimeout(t *testing.T) {
	assert.NoError(t, os.Setenv("FOO", "BAR"))
	ctx, cancel := NewWithTimeout(config.Project{}, time.Second)
	assert.NotEmpty(t, ctx.Date)
	assert.Equal(t, "BAR", ctx.Env["FOO"])

	cancel()
	<-ctx.Done()
	assert.EqualError(t, ctx.Err(), `context canceled`)
}
