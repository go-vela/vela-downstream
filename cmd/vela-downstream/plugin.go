// Copyright (c) 2021 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package main

import (
	"fmt"

	"github.com/go-vela/sdk-go/vela"
	"github.com/go-vela/types/library"

	"github.com/sirupsen/logrus"
)

// Plugin represents the configuration loaded for the plugin.
type Plugin struct {
	// build arguments loaded for the plugin
	Build *Build
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
	repos, err := p.Repo.Parse(p.Build.Branch)
	if err != nil {
		return err
	}

	// iterate through each repo from provided configuration
	for _, repo := range repos {
		// create new build type to store last successful build
		build := library.Build{}

		logrus.Infof("listing last 500 builds for %s", repo.GetFullName())

		// create options for listing builds
		opts := &vela.ListOptions{
			// set the default starting page for options
			Page: 1,
			// set the max per page for options
			PerPage: 100,
		}

		// create new slice of builds to store API results
		builds := []library.Build{}

		// loop to capture *ALL* the builds
		for {
			// send API call to capture a list of builds for the repo
			b, resp, err := client.Build.GetAll(repo.GetOrg(), repo.GetName(), opts)
			if err != nil {
				return fmt.Errorf("unable to list builds for %s: %w", repo.GetFullName(), err)
			}

			// add the results to the list of builds
			builds = append(builds, *b...)

			// break the loop if there is no more results
			// to page through or after 5 pages of results
			// giving us up to a total of 500 builds
			if resp.NextPage == 0 || resp.NextPage > 5 {
				break
			}

			// update the options for listing builds
			// to point at the next page
			opts.Page = resp.NextPage
		}

		logrus.Debugf("searching for latest %s build on branch %s with status %s",
			p.Build.Event,
			repo.GetBranch(),
			p.Build.Status,
		)

		// iterate through list of builds for the repo
		for _, b := range builds {
			// check if the build branch, event and status match
			if b.GetBranch() == repo.GetBranch() &&
				b.GetEvent() == p.Build.Event &&
				contains(p.Build.Status, b.GetStatus()) {
				// update the build object to the current build
				build = b

				// break out of the loop
				break
			}
		}

		// check if we found a build to restart
		if build.GetNumber() == 0 {
			return fmt.Errorf("no %s build on branch %s with status %s found for %s",
				p.Build.Event,
				repo.GetBranch(),
				p.Build.Status,
				repo.GetFullName(),
			)
		}

		logrus.Infof("Restarting build %s/%d", repo.GetFullName(), build.GetNumber())

		// send API call to restart the latest build for the repo
		b, _, err := client.Build.Restart(repo.GetOrg(), repo.GetName(), build.GetNumber())
		if err != nil {
			return fmt.Errorf("unable to restart build %s/%d", repo.GetFullName(), build.GetNumber())
		}

		logrus.Infof("New build created %s/%d", repo.GetFullName(), b.GetNumber())
	}

	return nil
}

// Validate verifies the plugin is properly configured.
func (p *Plugin) Validate() error {
	logrus.Debug("validating plugin configuration")

	// validate build configuration
	err := p.Build.Validate()
	if err != nil {
		return err
	}

	// validate config configuration
	err = p.Config.Validate()
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
