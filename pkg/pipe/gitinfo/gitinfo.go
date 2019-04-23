package gitinfo

import (
	"fmt"
	"github.com/apex/log"
	"github.com/devster/tarreleaser/pkg/context"
	"github.com/devster/tarreleaser/pkg/git"
	"github.com/devster/tarreleaser/pkg/pipe"
	"github.com/pkg/errors"
	"os"
)

type Pipe struct{}

func (Pipe) String() string {
	return "gathering git info"
}

func (Pipe) Default(ctx *context.Context) error {
	return nil
}

func (Pipe) Run(ctx *context.Context) error {
	if !git.HasGit() {
		return pipe.Skip("git not found in PATH, skipping")
	}

	if !git.IsRepo() {
		return pipe.Skip("current dir is not a git repository, skipping")
	}

	gitinfo, err := getInfo()
	if err != nil {
		return err
	}

	ctx.Git = gitinfo
	log.WithFields(log.Fields{
		"tag":    gitinfo.CurrentTag,
		"branch": gitinfo.Branch,
		"commit": fmt.Sprintf("%s %s - %s", gitinfo.ShortCommit, gitinfo.Commit.Author, gitinfo.Commit.Message),
	}).Info("git info")

	return nil
}

func getInfo() (context.GitInfo, error) {
	short, err := getShortCommit()
	if err != nil {
		return context.GitInfo{}, errors.Wrap(err, "couldn't get current commit")
	}

	full, err := getFullCommit()
	if err != nil {
		return context.GitInfo{}, errors.Wrap(err, "couldn't get current commit")
	}

	branch, err := getBranch()
	if err != nil {
		return context.GitInfo{}, errors.Wrap(err, "couldn't get current branch")
	}

	author, err := getAuthor()
	if err != nil {
		return context.GitInfo{}, errors.Wrap(err, "couldn't get author name")
	}

	msg, err := getMessage()
	if err != nil {
		return context.GitInfo{}, errors.Wrap(err, "couldn't get commit message")
	}

	tag, err := getTag()
	if err != nil {
		log.WithError(err).Warn("No tags defined")
	}

	gitinfo := context.GitInfo{
		ShortCommit: short,
		FullCommit:  full,
		CurrentTag:  tag,
		Branch:      branch,
		Commit: context.GitInfoCommit{
			Author:  author,
			Message: msg,
		},
	}

	if err = validateTag(gitinfo.CurrentTag); err != nil {
		log.Warnf("git tag %v was not made against commit %v, skipping tag", gitinfo.CurrentTag, gitinfo.FullCommit)
		gitinfo.CurrentTag = ""
	}

	return gitinfo, nil
}

func getShortCommit() (string, error) {
	return git.Run("show", "--format=%h", "HEAD", "-q")
}

func getFullCommit() (string, error) {
	return git.Run("show", "--format=%H", "HEAD", "-q")
}

func getAuthor() (string, error) {
	return git.Run("show", "--format=%aN", "HEAD", "-q")
}

func getMessage() (string, error) {
	return git.Run("show", "--format=%s", "HEAD", "-q")
}

func getTag() (string, error) {
	return git.Run("describe", "--tags", "--abbrev=0")
}

func getBranch() (string, error) {
	branch, err := git.Run("rev-parse", "--abbrev-ref", "HEAD")
	if err != nil {
		return "", err
	}

	// Provides consistent branch name on Travis
	travisBranch := os.Getenv("TRAVIS_BRANCH")
	if branch == "HEAD" && travisBranch != "" {
		branch = travisBranch
	}

	return branch, nil
}

func validateTag(tag string) (err error) {
	if tag == "" {
		return
	}
	_, err = git.Run("describe", "--exact-match", "--tags", "--match", tag)
	return
}
