// SPDX-License-Identifier: Apache-2.0

package main

import (
	"encoding/json"
	"fmt"
	"os"

	"time"

	"github.com/go-vela/types/constants"
	"github.com/go-vela/vela-downstream/version"

	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"

	_ "github.com/joho/godotenv/autoload"
)

func main() {
	// capture application version information
	v := version.New()

	// serialize the version information as pretty JSON
	bytes, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		logrus.Fatal(err)
	}

	// output the version information to stdout
	fmt.Fprintf(os.Stdout, "%s\n", string(bytes))

	// create new CLI application
	app := cli.NewApp()

	// Plugin Information

	app.Name = "vela-downstream"
	app.HelpName = "vela-downstream"
	app.Usage = "Vela Downstream plugin for triggering builds in other repos"
	app.Copyright = "Copyright 2020 Target Brands, Inc. All rights reserved."
	app.Authors = []*cli.Author{
		{
			Name:  "Vela Admins",
			Email: "vela@target.com",
		},
	}

	// Plugin Metadata

	app.Action = run
	app.Compiled = time.Now()
	app.Version = v.Semantic()

	// Plugin Flags

	app.Flags = []cli.Flag{

		&cli.StringFlag{
			EnvVars:  []string{"PARAMETER_LOG_LEVEL", "DOWNSTREAM_LOG_LEVEL"},
			FilePath: "/vela/parameters/downstream/log_level,/vela/secrets/downstream/log_level",
			Name:     "log.level",
			Usage:    "set log level - options: (trace|debug|info|warn|error|fatal|panic)",
			Value:    "info",
		},

		// Build Flags

		&cli.StringFlag{
			EnvVars:  []string{"PARAMETER_BRANCH", "DOWNSTREAM_BRANCH"},
			FilePath: "/vela/parameters/downstream/branch,/vela/secrets/downstream/branch",
			Name:     "build.branch",
			Usage:    "branch to trigger a build for the repo",
		},
		&cli.StringFlag{
			EnvVars:  []string{"PARAMETER_EVENT", "DOWNSTREAM_EVENT"},
			FilePath: "/vela/parameters/downstream/event,/vela/secrets/downstream/event",
			Name:     "build.event",
			Usage:    "event to trigger a build for the repo",
			Value:    constants.EventPush,
		},
		&cli.StringSliceFlag{
			EnvVars:  []string{"PARAMETER_STATUS", "DOWNSTREAM_STATUS"},
			FilePath: "/vela/parameters/downstream/status,/vela/secrets/downstream/status",
			Name:     "build.status",
			Usage:    "list of statuses to trigger a build for the repo",
			Value:    cli.NewStringSlice(constants.StatusSuccess),
		},
		&cli.BoolFlag{
			EnvVars:  []string{"PARAMETER_CONTINUE_ON_NOT_FOUND", "DOWNSTREAM_CONTINUE_ON_NOT_FOUND"},
			FilePath: "/vela/parameters/downstream/report_back,/vela/secrets/downstream/report_back",
			Name:     "build.continue",
			Usage:    "determine whether the downstream plugin should continue through repo list if a build is not found to restart",
		},

		// Build Check Flags

		&cli.BoolFlag{
			EnvVars:  []string{"PARAMETER_REPORT_BACK", "DOWNSTREAM_REPORT_BACK"},
			FilePath: "/vela/parameters/downstream/report_back,/vela/secrets/downstream/report_back",
			Name:     "build-check.enabled",
			Usage:    "determine whether the downstream plugin should wait for its triggered builds and report their statuses",
		},
		&cli.DurationFlag{
			EnvVars:  []string{"PARAMETER_TIMEOUT", "DOWNSTREAM_TIMEOUT"},
			FilePath: "/vela/parameters/downstream/timeout,/vela/secrets/downstream/timeout",
			Name:     "build-check.timeout",
			Usage:    "timeout for checking on triggered build statuses",
			Value:    30 * time.Minute,
		},
		&cli.StringSliceFlag{
			EnvVars:  []string{"PARAMETER_TARGET_STATUS", "DOWNSTREAM_TARGET_STATUS"},
			FilePath: "/vela/parameters/downstream/target_status,/vela/secrets/downstream/target_status",
			Name:     "build-check.status",
			Usage:    "list of statuses that constitute a successful triggered build",
			Value:    cli.NewStringSlice(constants.StatusSuccess),
		},

		// Config Flags

		&cli.StringFlag{
			EnvVars:  []string{"PARAMETER_SERVER", "DOWNSTREAM_SERVER"},
			FilePath: "/vela/parameters/downstream/server,/vela/secrets/downstream/server",
			Name:     "config.server",
			Usage:    "Vela server to authenticate with",
		},
		&cli.StringFlag{
			EnvVars:  []string{"PARAMETER_TOKEN", "DOWNSTREAM_TOKEN"},
			FilePath: "/vela/parameters/downstream/token,/vela/secrets/downstream/token",
			Name:     "config.token",
			Usage:    "user token to authenticate with the Vela server",
		},

		// Repo Flags

		&cli.StringSliceFlag{
			EnvVars:  []string{"PARAMETER_REPOS", "DOWNSTREAM_REPOS"},
			FilePath: "/vela/parameters/downstream/repos,/vela/secrets/downstream/repos",
			Name:     "repo.names",
			Usage:    "list of <org>/<repo> names to trigger",
		},
	}

	err = app.Run(os.Args)
	if err != nil {
		logrus.Fatal(err)
	}
}

// run executes the plugin based off the configuration provided.
func run(c *cli.Context) error {
	// set the log level for the plugin
	switch c.String("log.level") {
	case "t", "trace", "Trace", "TRACE":
		logrus.SetLevel(logrus.TraceLevel)
	case "d", "debug", "Debug", "DEBUG":
		logrus.SetLevel(logrus.DebugLevel)
	case "w", "warn", "Warn", "WARN":
		logrus.SetLevel(logrus.WarnLevel)
	case "e", "error", "Error", "ERROR":
		logrus.SetLevel(logrus.ErrorLevel)
	case "f", "fatal", "Fatal", "FATAL":
		logrus.SetLevel(logrus.FatalLevel)
	case "p", "panic", "Panic", "PANIC":
		logrus.SetLevel(logrus.PanicLevel)
	case "i", "info", "Info", "INFO":
		fallthrough
	default:
		logrus.SetLevel(logrus.InfoLevel)
	}

	logrus.WithFields(logrus.Fields{
		"code":     "https://github.com/go-vela/vela-downstream",
		"docs":     "https://go-vela.github.io/docs/plugins/registry/pipeline/downstream",
		"registry": "https://hub.docker.com/r/target/vela-downstream",
	}).Info("Vela Downstream Plugin")

	// create the plugin
	p := &Plugin{
		// build configuration
		Build: &Build{
			Branch:       c.String("build.branch"),
			Event:        c.String("build.event"),
			Status:       c.StringSlice("build.status"),
			Report:       c.Bool("build-check.enabled"),
			TargetStatus: c.StringSlice("build-check.status"),
			Timeout:      c.Duration("build-check.timeout"),
			Continue:     c.Bool("build.continue"),
		},
		// config configuration
		Config: &Config{
			Server:     c.String("config.server"),
			Token:      c.String("config.token"),
			AppName:    c.App.Name,
			AppVersion: c.App.Version,
		},
		// repo configuration
		Repo: &Repo{
			Names: c.StringSlice("repo.names"),
		},
	}

	// validate the plugin
	err := p.Validate()
	if err != nil {
		return err
	}

	// execute the plugin
	return p.Exec()
}
