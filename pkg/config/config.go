package config

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type Archive struct {
	Name string `yaml:",omitempty"`
	CompressionLevel int `yaml:"compression_level,omitempty"`
	IncludeFiles []string `yaml:"includes,omitempty"`
	ExcludeFiles []string `yaml:"excludes,omitempty"`
}

type Project struct {
	Dist string `yaml:",omitempty"`
	Archive Archive `yaml:",omitempty"`
}

func Load(file string) (config Project, err error) {
	cfgYml, err := ioutil.ReadFile(file)
	if err != nil {
		return config, err
	}

	err = yaml.Unmarshal(cfgYml, &config)
	log.WithField("config", fmt.Sprintf("%+v", config)).Debug("loaded config file")
	return config, err
}
