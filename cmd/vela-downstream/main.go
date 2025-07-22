// SPDX-License-Identifier: Apache-2.0

package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/mail"
	"os"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v3"

	_ "github.com/joho/godotenv/autoload"

	"github.com/go-vela/server/constants"
	"github.com/go-vela/vela-downstream/version"
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
	// Plugin Information
	cmd := cli.Command{
		Name:      "vela-downstream",
		Usage:     "Vela Downstream plugin for triggering builds in other repos",
		Copyright: "Copyright 2020 Target Brands, Inc. All rights reserved.",
		Authors: []any{
			&mail.Address{
				Name:    "Vela Admins",
				Address: "vela@target.com",
			},
		},
		Version: v.Semantic(),
		Action:  run,
	}

	// Plugin Flags

	cmd.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:  "log.level",
			Usage: "set log level - options: (trace|debug|info|warn|error|fatal|panic)",
			Value: "info",
			Sources: cli.NewValueSourceChain(
				cli.EnvVar("PARAMETER_LOG_LEVEL"),
				cli.EnvVar("DOWNSTREAM_LOG_LEVEL"),
				cli.File("/vela/parameters/downstream/log_level"),
				cli.File("/vela/secrets/downstream/log_level"),
			),
		},

		// Build Flags

		&cli.StringFlag{
			Name:  "build.branch",
			Usage: "branch to trigger a build for the repo",
			Sources: cli.NewValueSourceChain(
				cli.EnvVar("PARAMETER_BRANCH"),
				cli.EnvVar("DOWNSTREAM_BRANCH"),
				cli.File("/vela/parameters/downstream/branch"),
				cli.File("/vela/secrets/downstream/branch"),
			),
		},
		&cli.StringFlag{
			Name:  "build.event",
			Usage: "event to trigger a build for the repo",
			Value: constants.EventPush,
			Sources: cli.NewValueSourceChain(
				cli.EnvVar("PARAMETER_EVENT"),
				cli.EnvVar("DOWNSTREAM_EVENT"),
				cli.File("/vela/parameters/downstream/event"),
				cli.File("/vela/secrets/downstream/event"),
			),
		},
		&cli.StringSliceFlag{
			Name:  "build.status",
			Usage: "list of statuses to trigger a build for the repo",
			Value: []string{constants.StatusSuccess},
			Sources: cli.NewValueSourceChain(
				cli.EnvVar("PARAMETER_STATUS"),
				cli.EnvVar("DOWNSTREAM_STATUS"),
				cli.File("/vela/parameters/downstream/status"),
				cli.File("/vela/secrets/downstream/status"),
			),
		},
		&cli.BoolFlag{
			Name:  "build.continue",
			Usage: "determine whether the downstream plugin should continue through repo list if a build is not found to restart",
			Sources: cli.NewValueSourceChain(
				cli.EnvVar("PARAMETER_CONTINUE_ON_NOT_FOUND"),
				cli.EnvVar("DOWNSTREAM_CONTINUE_ON_NOT_FOUND"),
				cli.File("/vela/parameters/downstream/report_back"),
				cli.File("/vela/secrets/downstream/report_back"),
			),
		},

		// Build Check Flags

		&cli.BoolFlag{
			Name:  "build-check.enabled",
			Usage: "determine whether the downstream plugin should wait for its triggered builds and report their statuses",
			Sources: cli.NewValueSourceChain(
				cli.EnvVar("PARAMETER_REPORT_BACK"),
				cli.EnvVar("DOWNSTREAM_REPORT_BACK"),
				cli.File("/vela/parameters/downstream/report_back"),
				cli.File("/vela/secrets/downstream/report_back"),
			),
		},
		&cli.DurationFlag{
			Name:  "build-check.timeout",
			Usage: "timeout for checking on triggered build statuses",
			Value: 30 * time.Minute,
			Sources: cli.NewValueSourceChain(
				cli.EnvVar("PARAMETER_TIMEOUT"),
				cli.EnvVar("DOWNSTREAM_TIMEOUT"),
				cli.File("/vela/parameters/downstream/timeout"),
				cli.File("/vela/secrets/downstream/timeout"),
			),
		},
		&cli.StringSliceFlag{
			Name:  "build-check.status",
			Usage: "list of statuses that constitute a successful triggered build",
			Value: []string{constants.StatusSuccess},
			Sources: cli.NewValueSourceChain(
				cli.EnvVar("PARAMETER_TARGET_STATUS"),
				cli.EnvVar("DOWNSTREAM_TARGET_STATUS"),
				cli.File("/vela/parameters/downstream/target_status"),
				cli.File("/vela/secrets/downstream/target_status"),
			),
		},

		// Config Flags

		&cli.StringFlag{
			Name:  "config.server",
			Usage: "Vela server to authenticate with",
			Sources: cli.NewValueSourceChain(
				cli.EnvVar("PARAMETER_SERVER"),
				cli.EnvVar("DOWNSTREAM_SERVER"),
				cli.File("/vela/parameters/downstream/server"),
				cli.File("/vela/secrets/downstream/server"),
			),
		},
		&cli.StringFlag{
			Name:  "config.token",
			Usage: "user token to authenticate with the Vela server",
			Sources: cli.NewValueSourceChain(
				cli.EnvVar("PARAMETER_TOKEN"),
				cli.EnvVar("DOWNSTREAM_TOKEN"),
				cli.File("/vela/parameters/downstream/token"),
				cli.File("/vela/secrets/downstream/token"),
			),
		},
		&cli.IntFlag{
			Name:  "config.depth",
			Usage: "number of builds to search for downstream repositories",
			Value: 50,
			Sources: cli.NewValueSourceChain(
				cli.EnvVar("PARAMETER_DEPTH"),
				cli.EnvVar("DOWNSTREAM_DEPTH"),
				cli.File("/vela/parameters/downstream/depth"),
				cli.File("/vela/secrets/downstream/depth"),
			),
		},

		// Repo Flags

		&cli.StringSliceFlag{
			Name:  "repo.names",
			Usage: "list of <org>/<repo> names to trigger",
			Sources: cli.NewValueSourceChain(
				cli.EnvVar("PARAMETER_REPOS"),
				cli.EnvVar("DOWNSTREAM_REPOS"),
				cli.File("/vela/parameters/downstream/repos"),
				cli.File("/vela/secrets/downstream/repos"),
			),
		},
	}

	err = cmd.Run(context.Background(), os.Args)
	if err != nil {
		logrus.Fatal(err)
	}
}

// run executes the plugin based off the configuration provided.
func run(_ context.Context, c *cli.Command) error {
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
			Depth:      c.Int("config.depth"),
			AppName:    c.Name,
			AppVersion: c.Version,
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
