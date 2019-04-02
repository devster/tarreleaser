package context

import (
	"github.com/devster/tarreleaser/pkg/config"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNew(t *testing.T) {
	ctx := New(config.Project{})
	assert.NotEmpty(t, ctx.Date)
}
