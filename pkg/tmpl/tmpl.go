package tmpl

import (
	"bytes"
	"github.com/devster/tarreleaser/pkg/context"
	"text/template"
	"time"
)

// Template holds data that can be applied to a template string
type Template struct {
	fields fields
}

type fields map[string]interface{}

func New(ctx *context.Context) *Template {
	return &Template{
		fields: fields{
			"Tag":         ctx.Git.CurrentTag,
			"ShortCommit": ctx.Git.ShortCommit,
			"FullCommit":  ctx.Git.FullCommit,
			"Branch":      ctx.Git.Branch,
			"Commit":      ctx.Git.Commit,
			"Date":        ctx.Date.Format(time.RFC3339),
			"Timestamp":   ctx.Date.Unix(),
			"Env":         ctx.Env,
		},
	}
}

func (t *Template) Apply(s string) (string, error) {
	var out bytes.Buffer
	tmpl, err := template.New("tmpl").Option("missingkey=error").Parse(s)
	if err != nil {
		return "", err
	}

	err = tmpl.Execute(&out, t.fields)
	return out.String(), err
}
