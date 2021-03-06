package git

import (
	"errors"
	"github.com/apex/log"
	"os/exec"
	"strings"
)

func HasGit() bool {
	if _, err := exec.LookPath("git"); err != nil {
		return false
	}
	return true
}

// IsRepo returns true if current folder is a git repository
func IsRepo() bool {
	out, err := Run("rev-parse", "--is-inside-work-tree")
	return err == nil && strings.TrimSpace(out) == "true"
}

func Run(args ...string) (string, error) {
	var extraArgs = []string{
		"-c", "log.showSignature=false",
	}
	args = append(extraArgs, args...)

	var cmd = exec.Command("git", args...)
	log.WithField("args", args).Debug("running git")
	bts, err := cmd.CombinedOutput()
	out := strings.TrimSuffix(string(bts), "\n")
	log.WithField("output", out).Debug("git result")
	if err != nil {
		return "", errors.New(out)
	}

	return out, nil
}
