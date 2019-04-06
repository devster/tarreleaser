// Package targz implements the Archive interface providing tar.gz archiving and compression.
package targz

import (
	"archive/tar"
	"compress/gzip"
	"io"
	"os"
	"time"
)

// Archive as tar.gz
type Archive struct {
	gw *gzip.Writer
	tw *tar.Writer
}

// Close all closeables
func (a Archive) Close() error {
	if err := a.tw.Close(); err != nil {
		return err
	}
	return a.gw.Close()
}

// New tar.gz archive
func New(target io.Writer, level int) (*Archive, error) {
	gw, err := gzip.NewWriterLevel(target, level)
	if err != nil {
		return &Archive{}, err
	}
	tw := tar.NewWriter(gw)
	return &Archive{
		gw: gw,
		tw: tw,
	}, nil
}

// Add file to the archive
func (a Archive) Add(name, path string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()
	info, err := file.Stat()
	if err != nil {
		return err
	}
	header, err := tar.FileInfoHeader(info, name)
	if err != nil {
		return err
	}
	header.Name = name
	if err = a.tw.WriteHeader(header); err != nil {
		return err
	}
	if info.IsDir() {
		return nil
	}
	_, err = io.Copy(a.tw, file)
	return err
}

// Add file from content to the archive
func (a Archive) AddFromString(name, content string) error {
	header := new(tar.Header)
	header.Name = name
	header.Size = int64(len(content))
	header.Mode = 0644
	header.ModTime = time.Now()

	if err := a.tw.WriteHeader(header); err != nil {
		return err
	}

	_, err := a.tw.Write([]byte(content))
	return err
}
