package output

import (
	"fmt"
	"github.com/devster/tarreleaser/pkg/context"
	"github.com/devster/tarreleaser/pkg/tmpl"
)

type Pipe struct{}

func (Pipe) String() string {
	return ""
}

func (Pipe) Default(ctx *context.Context) error {
	return nil
}

func (Pipe) Run(ctx *context.Context) error {
	if ctx.OutputFormat == "" {
		return nil
	}

	output, err := tmpl.New(ctx).Apply(ctx.OutputFormat)
	if err != nil {
		return fmt.Errorf("output format error: %s", err.Error())
	}

	fmt.Print(output)

	return nil
}
