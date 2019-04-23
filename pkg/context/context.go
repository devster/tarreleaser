package context

import (
	ctx "context"
	"github.com/devster/tarreleaser/pkg/config"
	"os"
	"strings"
	"time"
)

type GitInfoCommit struct {
	Message string
	Author  string
}

type GitInfo struct {
	Branch      string
	CurrentTag  string
	Commit      GitInfoCommit
	ShortCommit string
	FullCommit  string
}

type Archive struct {
	Path string
	Name string
}

type Context struct {
	ctx.Context
	Config       config.Project
	Env          map[string]string
	Git          GitInfo
	Date         time.Time
	Archive      Archive
	SkipPublish  bool
	OutputFormat string
}

func NewWithTimeout(cfg config.Project, timeout time.Duration) (*Context, ctx.CancelFunc) {
	ctx, cancel := ctx.WithTimeout(ctx.Background(), timeout)
	return &Context{
		Context: ctx,
		Config:  cfg,
		Env:     splitEnv(os.Environ()),
		Date:    time.Now().UTC(),
	}, cancel
}

func splitEnv(env []string) map[string]string {
	r := map[string]string{}
	for _, e := range env {
		p := strings.SplitN(e, "=", 2)
		r[p[0]] = p[1]
	}
	return r
}
