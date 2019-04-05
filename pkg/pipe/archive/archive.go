package archive

import (
	"compress/gzip"
	"fmt"
	"github.com/devster/tarreleaser/pkg/archive/targz"
	"github.com/devster/tarreleaser/pkg/context"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"os"
	"path/filepath"
	"strings"
	"github.com/mattn/go-zglob"
	"github.com/campoy/unique"
)

type Pipe struct{}

func (Pipe) String() string {
	return "archive"
}

func (Pipe) Default(ctx *context.Context) error {
	if ctx.Config.Archive.Name == "" {
		ctx.Config.Archive.Name = "release.tar.gz"
	}

	if ctx.Config.Archive.CompressionLevel == 0 {
		ctx.Config.Archive.CompressionLevel = gzip.DefaultCompression
	}

	if len(ctx.Config.Archive.IncludeFiles) == 0 && len(ctx.Config.Archive.ExcludeFiles) == 0 {
		ctx.Config.Archive.ExcludeFiles = []string{
			".git",
		}
	}

	// Exclude the dist directory
	ctx.Config.Archive.ExcludeFiles = append(ctx.Config.Archive.ExcludeFiles, ctx.Config.Dist)

	if len(ctx.Config.Archive.IncludeFiles) == 0 {
		ctx.Config.Archive.IncludeFiles = []string{
			"./**/*",
		}
	}

	ctx.Config.Archive.IncludeFiles = betterGlob(ctx.Config.Archive.IncludeFiles)
	ctx.Config.Archive.ExcludeFiles = betterGlob(ctx.Config.Archive.ExcludeFiles)

	return nil
}

func (Pipe) Run(ctx *context.Context) error {
	archivePath := filepath.Join(ctx.Config.Dist, ctx.Config.Archive.Name)
	archiveFile, err := os.Create(archivePath)
	if err != nil {
		return errors.Wrapf(err, "failed to create archive file: %s", archivePath)
	}
	defer archiveFile.Close()

	log.WithFields(log.Fields{
		"archive": archivePath,
		"gzip_lvl": ctx.Config.Archive.CompressionLevel,
	}).Info("creating archive")

	a, err := targz.New(archiveFile, ctx.Config.Archive.CompressionLevel)
	if err != nil {
		return errors.Wrap(err, "failed to create archive")
	}
	defer a.Close()

	files, err := findFiles(ctx)
	if err != nil {
		return fmt.Errorf("failed to find files to archive: %s", err.Error())
	}

	for _, f := range files {
		log.Debugf("adding file: %s", f)

		if err = a.Add(f, f); err != nil {
			return fmt.Errorf("failed to add %s to the archive: %s", f, err.Error())
		}
	}

	log.Info("archive created with success")

	return nil
}

func findFiles(ctx *context.Context) (result []string, err error) {
	for _, glob := range ctx.Config.Archive.IncludeFiles {
		files, err := zglob.Glob(glob)
		if err != nil {
			return result, fmt.Errorf("include globbing failed for pattern %s: %s", glob, err.Error())
		}
		// excluding file that matches exclude glob pattern
		for _, f := range files {
			ok, err := isFileExcluded(ctx.Config.Archive.ExcludeFiles, f)
			if err != nil {
				return result, err
			}

			if !ok {
				result = append(result, f)
			}
		}
	}
	// remove duplicates
	unique.Slice(&result, func(i, j int) bool {
		return strings.Compare(result[i], result[j]) < 0
	})
	return
}

func isFileExcluded(patterns []string, file string) (bool, error) {
	for _, pattern := range patterns {
		ok, err := zglob.Match(pattern, file)
		if err != nil {
			return false, fmt.Errorf("exclude globbing failed for pattern %s: %s", pattern, err.Error())
		}
		if ok {
			return true, nil
		}
	}

	return false, nil
}

// Convert existing dirs into glob pattern to have the same effect that git does with the .gitignore file
func betterGlob(patterns []string) (result []string) {
	for _, glob := range patterns {
		if info, err := os.Stat(glob); err == nil && info.IsDir() {
			result = append(result, info.Name(), filepath.Join(info.Name(), "**/*"))
			continue
		}

		result = append(result,  glob)
	}

	return
}
