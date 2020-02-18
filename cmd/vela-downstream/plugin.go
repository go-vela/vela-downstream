// Copyright (c) 2020 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package main

import (
	"fmt"

	"github.com/go-vela/types/constants"
	"github.com/go-vela/types/library"

	"github.com/sirupsen/logrus"
)

// Plugin represents the configuration loaded for the plugin.
type Plugin struct {
	// config arguments loaded for the plugin
	Config *Config
	// repo arguments loaded for the plugin
	Repo *Repo
}

// Exec formats and runs the commands for triggering builds in Vela.
func (p *Plugin) Exec() error {
	logrus.Debug("running plugin with provided configuration")

	// create new Vela client from config configuration
	client, err := p.Config.New()
	if err != nil {
		return err
	}

	// parse list of repos to trigger builds on
	repos, err := p.Repo.Parse()
	if err != nil {
		return err
	}

	// iterate through each repo from provided configuration
	for _, repo := range repos {
		// create new build type to store last successful build
		build := new(library.Build)

		// send API call to capture a list of builds for the repo
		builds, _, err := client.Build.GetAll(repo.GetOrg(), repo.GetName())
		if err != nil {
			return fmt.Errorf("unable to capture latest build for %s: %w", repo.GetFullName(), err)
		}

		// iterate through list of builds for the repo
		for _, b := range *builds {
			// check if the build branch matches and was successful
			if b.GetBranch() == repo.GetBranch() && b.GetStatus() == constants.StatusSuccess {
				build = &b
				break
			}
		}

		// send API call to restart the latest build for the repo
		_, _, err = client.Build.Restart(repo.GetOrg(), repo.GetName(), build.GetNumber())
		if err != nil {
			return fmt.Errorf("unable to restart build %s/%d", repo.GetFullName(), build.GetNumber())
		}
	}

	return nil
}

// Validate verifies the plugin is properly configured.
func (p *Plugin) Validate() error {
	logrus.Debug("validating plugin configuration")

	// validate config configuration
	err := p.Config.Validate()
	if err != nil {
		return err
	}

	// validate repo configuration
	err = p.Repo.Validate()
	if err != nil {
		return err
	}

	return nil
}
