package config

import (
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"os"
	"testing"
)

func TestLoad(t *testing.T) {
	cfgFile, err := ioutil.TempFile("", "")
	assert.NoError(t, err)
	_, err = cfgFile.Write([]byte(`
  archive:
    name: "release.tar.gz"
    excludes:
      - "myfile"
`))
	assert.NoError(t, err)
	assert.NoError(t, cfgFile.Close())
	defer os.Remove(cfgFile.Name())

	cfg, err := Load(cfgFile.Name())

	assert.NoError(t, err)
	assert.Equal(t, "release.tar.gz", cfg.Archive.Name)
	assert.Equal(t, "myfile", cfg.Archive.ExcludeFiles[0])
}

func TestLoad_ConfigNotFound(t *testing.T) {
	_, err := Load("unknown.file")

	assert.EqualError(t, err, "open unknown.file: no such file or directory")
}
