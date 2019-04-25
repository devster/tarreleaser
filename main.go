package main

import (
	"fmt"
	"github.com/apex/log"
	"github.com/apex/log/handlers/cli"
	"github.com/caarlos0/ctrlc"
	"github.com/devster/tarreleaser/pkg/config"
	"github.com/devster/tarreleaser/pkg/context"
	"github.com/devster/tarreleaser/pkg/pipeline"
	"github.com/devster/tarreleaser/pkg/static"
	"github.com/fatih/color"
	"gopkg.in/alecthomas/kingpin.v2"
	"io/ioutil"
	"os"
	"time"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

const defaultConfigFile = ".tarreleaser.yml"

type releaseOptions struct {
	Config       string
	SkipPublish  bool
	Timeout      time.Duration
	OutputFormat string
}

func main() {
	// enable colored output on travis/circleci
	if os.Getenv("CI") != "" {
		color.NoColor = false
	}
	log.SetHandler(cli.Default)

	// Cli app
	app := kingpin.New("tarreleaser", "Build and publish your app as tarball")
	app.Version(fmt.Sprintf("%v, commit %v, built at %v", version, commit, date))
	app.HelpFlag.Short('h')
	debug := app.Flag("debug", "Enable debug mode").Bool()
	quiet := app.Flag("quiet", "Enable silent mode (only display errors)").Short('q').Bool()

	// Release cli command
	rOptions := releaseOptions{}
	releaseCmd := app.Command("release", "Releases the current project").Default()
	releaseCmd.Flag("config", "Load configuration from file").Short('c').Default(defaultConfigFile).StringVar(&rOptions.Config)
	releaseCmd.Flag("skip-publish", "Skips publishing artifacts").Short('s').BoolVar(&rOptions.SkipPublish)
	releaseCmd.Flag("timeout", "Timeout to the entire release process").Default("30m").DurationVar(&rOptions.Timeout)
	releaseCmd.Flag("output", "Format the output. Ex: -o '{{.Archive.Name}}'").Short('o').StringVar(&rOptions.OutputFormat)

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
		start := time.Now()
		log.Infof(color.New(color.Bold).Sprintf("releasing using tarreleaser %s...", version))

		if err := releaseProject(rOptions); err != nil {
			log.WithError(err).Fatalf(color.New(color.Bold).Sprintf("release failed after %0.2fs", time.Since(start).Seconds()))
		}
		log.Infof(color.New(color.Bold).Sprintf("release succeeded after %0.2fs", time.Since(start).Seconds()))

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

	ctx, cancel := context.NewWithTimeout(cfg, options.Timeout)
	defer cancel()
	ctx.SkipPublish = options.SkipPublish
	ctx.OutputFormat = options.OutputFormat

	return ctrlc.Default.Run(ctx, func() error {
		return pipeline.Run(ctx)
	})
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
