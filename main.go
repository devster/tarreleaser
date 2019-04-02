package main

import (
	"fmt"
	"github.com/devster/tarreleaser/pkg/config"
	"github.com/devster/tarreleaser/pkg/context"
	pkglog "github.com/devster/tarreleaser/pkg/log"
	"github.com/devster/tarreleaser/pkg/pipeline"
	"github.com/devster/tarreleaser/pkg/static"
	log "github.com/sirupsen/logrus"
	"gopkg.in/alecthomas/kingpin.v2"
	"io/ioutil"
	"os"
)

var (
	version = "dev"
)

const defaultConfigFile = ".tarreleaser.yml"

type releaseOptions struct {
	Config       string
	SkipPublish  bool
	RmDist       bool
}

func init() {
	log.SetFormatter(pkglog.TextFormatter)
}

func main() {
	//defer fmt.Println()

	// Cli app
	app := kingpin.New("tarreleaser", "Build and publish your app as tarball")
	app.Version(fmt.Sprintf("%v", version))
	app.HelpFlag.Short('h')
	debug := app.Flag("debug", "Enable debug mode").Bool()
	quiet := app.Flag("quiet", "Enable silent mode (only display errors)").Short('q').Bool()

	// Release cli command
	rOptions := releaseOptions{}
	releaseCmd := app.Command("release", "Releases the current project").Default()
	releaseCmd.Flag("config", "Load configuration from file").Short('c').Default(defaultConfigFile).StringVar(&rOptions.Config)
	releaseCmd.Flag("skip-publish", "Skips publishing artifacts").Short('s').BoolVar(&rOptions.SkipPublish)

	// Init config file cli command
	initCmd := app.Command("init", fmt.Sprintf("Generate a %v file", defaultConfigFile))

	cmd := kingpin.MustParse(app.Parse(os.Args[1:]))

	if *quiet {
		log.SetLevel(log.ErrorLevel)
	}

	if *debug {
		log.SetLevel(log.DebugLevel)
	}

	switch cmd {
	case releaseCmd.FullCommand():
		log.WithFields(log.Fields{
			"config": rOptions.Config,
			"skip-publish": rOptions.SkipPublish,
		}).Infof("releasing using tarreleaser %s...", version)

		if err := releaseProject(rOptions); err != nil {
			log.WithError(err).Fatal("release failed")
		}

		log.Info("release succeeded")

	case initCmd.FullCommand():
		if err := initProject(defaultConfigFile); err != nil {
			log.WithError(err).Fatal("failed to init project")
		}

		log.WithField("file", defaultConfigFile).Info("config created; please edit accordingly to your needs")
	}
}

func releaseProject(options releaseOptions) error {
	cfg, err := config.Load(options.Config)
	if err != nil {
		log.WithError(err).Fatal("failed to load config")
	}

	ctx := context.New(cfg)

	if err := pipeline.Run(ctx); err != nil {
		return err
	}
	// 1. Gathering data from project: hash commit, last tag, last commit message etc
	// 2. Build the tar
	// https://gist.github.com/jonmorehouse/9060515
	// 3. Generate a checksum file
	// 4. Publish

	return nil
}

func initProject(filename string) error {
	if _, err := os.Stat(filename); !os.IsNotExist(err) {
		if err != nil {
			return err
		}
		return fmt.Errorf("%s already exists", filename)
	}

	log.Infof("Generating example %v file", filename)
	return ioutil.WriteFile(filename, []byte(static.ExampleConfig), 0644)
}
