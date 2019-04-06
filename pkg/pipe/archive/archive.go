package archive

import (
	"compress/gzip"
	"fmt"
	"github.com/apex/log"
	"github.com/campoy/unique"
	"github.com/devster/tarreleaser/pkg/archive/targz"
	"github.com/devster/tarreleaser/pkg/context"
	"github.com/devster/tarreleaser/pkg/static"
	"github.com/devster/tarreleaser/pkg/tmpl"
	"github.com/dustin/go-humanize"
	"github.com/mattn/go-zglob"
	"github.com/pkg/errors"
	"os"
	"path/filepath"
	"strings"
)

type Pipe struct{}

func (Pipe) String() string {
	return "archiving"
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

	// Release info file defaults
	if ctx.Config.Archive.InfoFile.Name != "" || ctx.Config.Archive.InfoFile.Content != "" {
		if ctx.Config.Archive.InfoFile.Name == "" {
			ctx.Config.Archive.InfoFile.Name = "release.txt"
		}

		if ctx.Config.Archive.InfoFile.Content == "" {
			ctx.Config.Archive.InfoFile.Content = static.DefaultReleaseFileContent
		}
	}

	return nil
}

func (Pipe) Run(ctx *context.Context) error {
	ctx.Config.Archive.IncludeFiles = betterGlob(ctx.Config.Archive.IncludeFiles)
	ctx.Config.Archive.ExcludeFiles = betterGlob(ctx.Config.Archive.ExcludeFiles)

	t := tmpl.New(ctx)

	archiveName, err := t.Apply(ctx.Config.Archive.Name)
	if err != nil {
		return err
	}
	archivePath := filepath.Join(ctx.Config.Dist, archiveName)
	archiveFile, err := os.Create(archivePath)
	if err != nil {
		return errors.Wrapf(err, "failed to create archive file: %s", archivePath)
	}
	defer archiveFile.Close()

	wrap, err := t.Apply(ctx.Config.Archive.WrapInDirectory)
	if err != nil {
		return err
	}

	log.WithFields(log.Fields{
		"archive":  archivePath,
		"gzip_lvl": ctx.Config.Archive.CompressionLevel,
		"wrap":     wrap,
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
		name := filepath.Join(wrap, f)
		log.Debugf("adding file: %s", f)

		if err = a.Add(name, f); err != nil {
			return fmt.Errorf("failed to add %s to the archive: %s", f, err.Error())
		}
	}

	if err := addReleaseInfoFile(a, ctx, wrap); err != nil {
		return errors.Wrap(err, "failed to add release file info")
	}

	archiveInfo, err := archiveFile.Stat()
	if err != nil {
		return errors.Wrap(err, "unable to retrieve stat on archive file")
	}

	log.WithFields(log.Fields{
		"files": len(files),
		"size":  humanize.Bytes(uint64(archiveInfo.Size())),
	}).Info("archive created with success")

	ctx.Archive.Path = archivePath
	ctx.Archive.Name = archiveName

	return nil
}

func findFiles(ctx *context.Context) (result []string, err error) {
	for _, glob := range ctx.Config.Archive.IncludeFiles {
		files, err := zglob.Glob(glob)
		if err != nil {
			return result, fmt.Errorf("include globbing failed for pattern %s: %s", glob, err.Error())
		}
		// excluding files that matches exclude glob pattern
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

		result = append(result, glob)
	}

	return
}

func addReleaseInfoFile(a *targz.Archive, ctx *context.Context, wrap string) error {
	if ctx.Config.Archive.InfoFile.Name == "" {
		return nil
	}

	t := tmpl.New(ctx)

	name, err := t.Apply(ctx.Config.Archive.InfoFile.Name)
	if err != nil {
		return err
	}
	name = filepath.Join(wrap, name)

	log.WithField("file", name).Info("adding release info file")

	content, err := t.Apply(ctx.Config.Archive.InfoFile.Content)
	if err != nil {
		return err
	}

	return a.AddFromString(name, content)
}
