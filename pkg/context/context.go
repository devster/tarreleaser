package context

import (
	ctx "context"
	"github.com/devster/tarreleaser/pkg/config"
	"time"
)

type GitInfoCommit struct {
	Message string
	Author string
}

type GitInfo struct {
	Branch string
	CurrentTag  string
	Commit      GitInfoCommit
	ShortCommit string
	FullCommit  string
}

type Context struct {
	ctx.Context
	Config config.Project
	Git GitInfo
	Date time.Time
}

func New(cfg config.Project) *Context {
	return &Context{
		Context: ctx.Background(),
		Config: cfg,
		Date: time.Now().UTC(),
	}
}
