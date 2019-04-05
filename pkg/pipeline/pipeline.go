package pipeline

import (
	"fmt"
	"github.com/devster/tarreleaser/pkg/context"
	pkglog "github.com/devster/tarreleaser/pkg/log"
	"github.com/devster/tarreleaser/pkg/pipe/archive"
	"github.com/devster/tarreleaser/pkg/pipe/dist"
	"github.com/devster/tarreleaser/pkg/pipe/gitinfo"
)

var Pipes = []Pipe{
	gitinfo.Pipe{},
	dist.Pipe{},
	archive.Pipe{},
}

type Pipe interface {
	fmt.Stringer

	Default(ctx *context.Context) error
	Run(ctx *context.Context) error
}

func Run(ctx *context.Context) error {
	// Run Defaults
	for _, pipe := range Pipes {
		pkglog.TextFormatter.Prefix = pipe.String()

		if err := pipe.Default(ctx); err != nil {
			return err
		}
	}

	// Run pipes!
	for _, pipe := range Pipes {
		pkglog.TextFormatter.Prefix = pipe.String()

		if err := pipe.Run(ctx); err != nil {
			return err
		}
	}

	pkglog.TextFormatter.Prefix = ""

	return nil
}
