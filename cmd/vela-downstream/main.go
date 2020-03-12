// Copyright (c) 2020 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package main

import (
	"os"

	"time"

	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"

	_ "github.com/joho/godotenv/autoload"
)

func main() {
	app := cli.NewApp()

	// Plugin Information

	app.Name = "vela-downstream"
	app.HelpName = "vela-downstream"
	app.Usage = "Vela Downstream plugin for triggering builds in other repos"
	app.Copyright = "Copyright (c) 2020 Target Brands, Inc. All rights reserved."
	app.Authors = []cli.Author{
		{
			Name:  "Vela Admins",
			Email: "vela@target.com",
		},
	}

	// Plugin Metadata

	app.Compiled = time.Now()
	app.Action = run

	// Plugin Flags

	app.Flags = []cli.Flag{

		cli.StringFlag{
			EnvVar: "PARAMETER_LOG_LEVEL,VELA_LOG_LEVEL,DOWNSTREAM_LOG_LEVEL",
			Name:   "log.level",
			Usage:  "set log level - options: (trace|debug|info|warn|error|fatal|panic)",
			Value:  "info",
		},

		// Config Flags

		cli.StringFlag{
			EnvVar: "PARAMETER_SERVER,CONFIG_SERVER,VELA_SERVER,DOWNSTREAM_SERVER",
			Name:   "config.server",
			Usage:  "Vela server to authenticate with",
		},
		cli.StringFlag{
			EnvVar: "PARAMETER_TOKEN,CONFIG_TOKEN,VELA_TOKEN,DOWNSTREAM_TOKEN",
			Name:   "config.token",
			Usage:  "user token to authenticate with the Vela server",
		},

		// Repo Flags

		cli.StringFlag{
			EnvVar: "PARAMETER_BRANCH,REPO_BRANCH",
			Name:   "repo.branch",
			Usage:  "default branch to trigger builds for the repo",
			Value:  "master",
		},
		cli.StringSliceFlag{
			EnvVar: "PARAMETER_REPOS,REPO_NAMES,DOWNSTREAM_REPOS",
			Name:   "repo.names",
			Usage:  "list of <org>/<repo> names to trigger",
		},
	}

	err := app.Run(os.Args)
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
		"docs":     "https://go-vela.github.io/docs/plugins/registry/downstream",
		"registry": "https://hub.docker.com/r/target/vela-downstream",
	}).Info("Vela Downstream Plugin")

	// create the plugin
	p := &Plugin{
		// config configuration
		Config: &Config{
			Server: c.String("config.server"),
			Token:  c.String("config.token"),
		},
		// repo configuration
		Repo: &Repo{
			Branch: c.String("repo.branch"),
			Names:  c.StringSlice("repo.names"),
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
