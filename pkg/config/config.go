package config

import (
	"fmt"
	"github.com/apex/log"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type ReleaseInfoFile struct {
	Name    string `yaml:",omitempty"`
	Content string `yaml:",omitempty"`
}

type Archive struct {
	Name             string          `yaml:",omitempty"`
	CompressionLevel int             `yaml:"compression_level,omitempty"`
	IncludeFiles     []string        `yaml:"includes,omitempty"`
	ExcludeFiles     []string        `yaml:"excludes,omitempty"`
	WrapInDirectory  string          `yaml:"wrap_in_directory,omitempty"`
	InfoFile         ReleaseInfoFile `yaml:"info_file,omitempty"`
}

type S3 struct {
	Folder  string
	Bucket  string
	Region  string
	ACL     string
	Profile string
}

type Publish struct {
	S3 S3 `yaml:"s3,omitempty"`
}

type Project struct {
	Dist    string  `yaml:",omitempty"`
	Archive Archive `yaml:",omitempty"`
	Publish Publish `yaml:",omitempty"`
}

func Load(file string) (config Project, err error) {
	log.WithField("file", file).Info("loading config file")

	cfgYml, err := ioutil.ReadFile(file)
	if err != nil {
		return config, err
	}

	err = yaml.Unmarshal(cfgYml, &config)
	log.WithField("config", fmt.Sprintf("%+v", config)).Debug("loaded config file")
	return config, err
}
