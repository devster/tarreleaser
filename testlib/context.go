package testlib

import (
	"github.com/devster/tarreleaser/pkg/config"
	"github.com/devster/tarreleaser/pkg/context"
	"time"
)

func NewContext() *context.Context {
	ctx, _ := context.NewWithTimeout(config.Project{}, 15*time.Minute)
	return ctx
}
