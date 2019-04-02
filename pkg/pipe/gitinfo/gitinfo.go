package gitinfo

import (
	"fmt"
	"github.com/devster/tarreleaser/pkg/context"
	"github.com/devster/tarreleaser/pkg/git"
	"github.com/devster/tarreleaser/pkg/pipe"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

type Pipe struct {pipe.NoDefault}

func (Pipe) String() string {
	return "gitinfo"
}

func (Pipe) Run(ctx *context.Context) error {
	log.Info("gathering git repository infos")

	if !git.HasGit() {
		log.Warn("git not found in PATH, skipping")
		return nil
	}

	if !git.IsRepo() {
		log.Warn("current dir is not a git repository, skipping")
		return nil
	}

	gitinfo, err := getInfo()
	if err != nil {
		return err
	}

	ctx.Git = gitinfo
	log.WithFields(log.Fields{
		"tag": gitinfo.CurrentTag,
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
		log.WithError(err).Warn("Unable to retrieve the current tag")
	}

	gitinfo := context.GitInfo{
		ShortCommit: short,
		FullCommit: full,
		CurrentTag: tag,
		Branch: branch,
		Commit: context.GitInfoCommit{
			Author: author,
			Message: msg,
		},
	}

	return gitinfo, nil
}

func getShortCommit() (string, error) {
	return git.Run("show", "--format='%h'", "HEAD", "-q")
}

func getFullCommit() (string, error) {
	return git.Run("show", "--format='%H'", "HEAD", "-q")
}

func getAuthor() (string, error) {
	return git.Run("show", "--format='%aN'", "HEAD", "-q")
}

func getMessage() (string, error) {
	return git.Run("show", "--format='%s'", "HEAD", "-q")
}

func getTag() (string, error) {
	return git.Run("describe", "--tags", "--abbrev=0")
}

func getBranch() (string, error) {
	return git.Run("rev-parse", "--abbrev-ref", "HEAD")
}