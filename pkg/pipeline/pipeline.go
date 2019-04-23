package pipeline

import (
	"fmt"
	"github.com/apex/log"
	"github.com/apex/log/handlers/cli"
	"github.com/devster/tarreleaser/pkg/context"
	"github.com/devster/tarreleaser/pkg/pipe"
	"github.com/devster/tarreleaser/pkg/pipe/archive"
	"github.com/devster/tarreleaser/pkg/pipe/dist"
	"github.com/devster/tarreleaser/pkg/pipe/gitinfo"
	"github.com/devster/tarreleaser/pkg/pipe/output"
	"github.com/devster/tarreleaser/pkg/pipe/s3"
	"github.com/fatih/color"
	"strings"
)

var Pipes = []Pipe{
	gitinfo.Pipe{},
	dist.Pipe{},
	archive.Pipe{},
	s3.Pipe{},
	output.Pipe{},
}

type Pipe interface {
	fmt.Stringer

	Default(ctx *context.Context) error
	Run(ctx *context.Context) error
}

func Run(ctx *context.Context) error {
	// Run Defaults
	for _, pipe := range Pipes {
		if err := pipe.Default(ctx); err != nil {
			return err
		}
	}

	// Run pipes!
	for _, p := range Pipes {
		if p.String() != "" {
			log.Info(color.New(color.Bold).Sprintf(strings.ToUpper(p.String())))
		}
		cli.Default.Padding = 6
		if err := p.Run(ctx); err != nil {
			if !pipe.IsSkip(err) {
				return err
			}

			log.WithError(err).Warn("pipe skipped")
		}
		cli.Default.Padding = 3
	}

	return nil
}
