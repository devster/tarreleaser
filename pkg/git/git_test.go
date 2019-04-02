package git

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestHasGit(t *testing.T) {
	assert.True(t, HasGit())
}

func TestRun(t *testing.T) {
	assert := assert.New(t)

	out, err := Run("status")
	assert.NoError(err)
	assert.NotEmpty(out)

	out, err = Run("command-unknown")
	assert.Error(err)
	assert.Empty(out)
}

func TestIsRepo(t *testing.T) {
	assert.True(t, IsRepo(), "current folder should be a git repo")

	assert.NoError(t, os.Chdir(os.TempDir()))
	assert.False(t, IsRepo(), os.TempDir()+" folder should not be a git repo")
}