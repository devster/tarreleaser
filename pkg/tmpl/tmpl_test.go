package tmpl

import (
	"github.com/devster/tarreleaser/testlib"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNew(t *testing.T) {
	ctx := testlib.NewContext()
	ctx.Env = map[string]string{
		"FOO": "BAR",
	}
	ctx.Git.CurrentTag = "v1.2.3"
	ctx.Git.Commit.Message = "message"
	ctx.Git.Commit.Author = "john"
	ctx.Git.Branch = "master"
	ctx.Git.ShortCommit = "shortcommit"
	ctx.Git.FullCommit = "fullcommit"

	var tests = map[string]string{
		"v1.2.3":         "{{ .Tag }}",
		"BAR":            "{{ .Env.FOO }}",
		"message - john": "{{ .Commit.Message }} - {{ .Commit.Author }}",
		"master":         "{{ .Branch }}",
		"shortcommit":    "{{ .ShortCommit }}",
		"fullcommit":     "{{ .FullCommit }}",
	}

	for expect, tmpl := range tests {
		t.Run(expect, func(t *testing.T) {
			result, err := New(ctx).Apply(tmpl)
			assert.NoError(t, err)
			assert.Equal(t, expect, result)
		})
	}

	date, err := New(ctx).Apply("{{ .Date }}")
	assert.NoError(t, err)
	assert.NotEmpty(t, date)

	ts, err := New(ctx).Apply("{{ .Timestamp }}")
	assert.NoError(t, err)
	assert.NotEmpty(t, ts)

	_, err = New(ctx).Apply("{{ .Unknown }}")
	assert.Error(t, err)

	_, err = New(ctx).Apply("{{ .Tag }")
	assert.Error(t, err)
}
