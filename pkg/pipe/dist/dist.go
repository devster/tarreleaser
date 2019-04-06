package dist

import (
	"fmt"
	"github.com/apex/log"
	"github.com/devster/tarreleaser/pkg/context"
	"github.com/devster/tarreleaser/pkg/pipe"
	"os"
)

type Pipe struct{}

func (Pipe) String() string {
	return "checking dist"
}

func (Pipe) Default(ctx *context.Context) error {
	if ctx.Config.Dist == "" {
		ctx.Config.Dist = "./dist"
	}

	return nil
}

func (Pipe) Run(ctx *context.Context) error {
	_, err := os.Stat(ctx.Config.Dist)
	if os.IsNotExist(err) {
		log.WithField("path", ctx.Config.Dist).Info("creating dist directory")
		return os.MkdirAll(ctx.Config.Dist, 0755)
	}

	return pipe.Skip(fmt.Sprintf("%s already exists", ctx.Config.Dist))
}
