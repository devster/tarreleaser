package dist

import (
	"github.com/devster/tarreleaser/pkg/config"
	"github.com/devster/tarreleaser/pkg/context"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"path/filepath"
	"testing"
)

func TestPipe_Run(t *testing.T) {
	folder, err := ioutil.TempDir("", "disttest")
	assert.NoError(t, err)
	var dist = filepath.Join(folder, "dist")
	assert.NoError(
		t,
		Pipe{}.Run(
			&context.Context{
				Config: config.Project{
					Dist: dist,
				},
			},
		),
	)

	assert.DirExists(t, dist)
}

func TestPipe_Default(t *testing.T) {
	ctx := &context.Context{
		Config: config.Project{},
	}

	assert.NoError(t, Pipe{}.Default(ctx))
	assert.Equal(t, "./dist", ctx.Config.Dist)
}

func TestPipe_String(t *testing.T) {
	assert.NotEmpty(t, Pipe{}.String())
}
