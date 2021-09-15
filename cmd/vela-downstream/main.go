// Copyright (c) 2021 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

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
	app.Copyright = "Copyright (c) 2021 Target Brands, Inc. All rights reserved."
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

		&cli.StringFlag{
			EnvVars:  []string{"PARAMETER_BRANCH", "DOWNSTREAM_BRANCH"},
			FilePath: "/vela/parameters/downstream/branch,/vela/secrets/downstream/branch",
			Name:     "repo.branch",
			Usage:    "default branch to trigger builds for the repo",
			Value:    "master",
		},
		&cli.StringSliceFlag{
			EnvVars:  []string{"PARAMETER_REPOS", "DOWNSTREAM_REPOS"},
			FilePath: "/vela/parameters/downstream/repos,/vela/secrets/downstream/repos",
			Name:     "repo.names",
			Usage:    "list of <org>/<repo> names to trigger",
		},
		&cli.StringSliceFlag{
			EnvVars:  []string{"PARAMETER_STATUS", "DOWNSTREAM_STATUS"},
			FilePath: "/vela/parameters/downstream/status,/vela/secrets/downstream/status",
			Name:     "status",
			Usage:    "status of last build to trigger - example: (error|failure|running|success)",
			Value:    cli.NewStringSlice(constants.StatusSuccess),
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
		"docs":     "https://go-vela.github.io/docs/plugins/registry/downstream",
		"registry": "https://hub.docker.com/r/target/vela-downstream",
	}).Info("Vela Downstream Plugin")

	// create the plugin
	p := &Plugin{
		// config configuration
		Config: &Config{
			Server:     c.String("config.server"),
			Token:      c.String("config.token"),
			AppName:    c.App.Name,
			AppVersion: c.App.Version,
		},
		// repo configuration
		Repo: &Repo{
			Branch: c.String("repo.branch"),
			Names:  c.StringSlice("repo.names"),
		},
		Status: &Status{"status"},
	}

	// validate the plugin
	err := p.Validate()
	if err != nil {
		return err
	}

	// execute the plugin
	return p.Exec()
}
