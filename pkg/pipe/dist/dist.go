package dist

import (
	"github.com/devster/tarreleaser/pkg/context"
	log "github.com/sirupsen/logrus"
	"os"
)

type Pipe struct {}

func (Pipe) String() string {
	return "dist"
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
		err = os.MkdirAll(ctx.Config.Dist, 0755)
	}

	return err
}