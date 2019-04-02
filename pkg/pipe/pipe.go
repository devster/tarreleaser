package pipe

import "github.com/devster/tarreleaser/pkg/context"

type NoDefault struct {}

func (NoDefault) Default(ctx *context.Context) error {
	return nil
}