// SPDX-License-Identifier: Apache-2.0

package main

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/go-vela/sdk-go/vela"
	api "github.com/go-vela/server/api/types"
	"github.com/go-vela/server/constants"
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

	rBMap := make(map[*api.Repo]int)

	// parse list of repos to trigger builds on
	repos, err := p.Repo.Parse(p.Build.Branch)
	if err != nil {
		return err
	}

	// iterate through each repo from provided configuration
	for _, repo := range repos {
		// create new build type to store last successful build
		build := api.Build{}

		logrus.Infof("searching last %d %s builds with branch %s for %s", p.Config.Depth, p.Build.Event, repo.GetBranch(), repo.GetFullName())

		// create options for listing builds
		//
		// https://pkg.go.dev/github.com/go-vela/sdk-go/vela#BuildListOptions
		opts := &vela.BuildListOptions{
			Branch: repo.GetBranch(),
			Event:  p.Build.Event,
			// https://pkg.go.dev/github.com/go-vela/sdk-go/vela#ListOptions
			ListOptions: vela.ListOptions{
				// set the default starting page for options
				Page: 1,
				// set the max per page for options
				PerPage: 10,
			},
		}

		// loop to capture *ALL* the builds
		for {
			// send API call to capture a list of builds for the repo
			//
			// https://pkg.go.dev/github.com/go-vela/sdk-go/vela#BuildService.GetAll
			builds, resp, err := client.Build.GetAll(repo.GetOrg(), repo.GetName(), opts)
			if err != nil {
				return fmt.Errorf("unable to list builds for %s: %w", repo.GetFullName(), err)
			}

			// iterate through list of builds for the repo
			for _, b := range *builds {
				// check if the build branch, event and status match
				if contains(p.Build.Status, b.GetStatus()) || contains(p.Build.Status, "any") {
					// update the build object to the current build
					build = b

					logrus.Infof("found %s build %s/%d on branch %s with status %s", p.Build.Event, repo.GetFullName(), build.GetNumber(), repo.GetBranch(), build.GetStatus())

					// break out of the loop
					break
				}
			}

			// break the loop if there is no more results
			// to page through or after 50 pages of results
			// giving us up to a total of 500 builds
			if resp.NextPage == 0 || resp.NextPage > 50 {
				break
			}

			// update the options for listing builds
			// to point at the next page
			opts.ListOptions.Page = resp.NextPage
		}

		// check if we found a build to restart
		if build.GetNumber() == 0 {
			msg := fmt.Sprintf("no %s build on branch %s with status %s found for %s",
				p.Build.Event,
				repo.GetBranch(),
				p.Build.Status,
				repo.GetFullName(),
			)

			if p.Build.Continue {
				logrus.Warn(msg)

				continue
			}

			return errors.New(msg)
		}

		logrus.Infof("restarting build %s/%d", repo.GetFullName(), build.GetNumber())

		// send API call to restart the latest build for the repo
		//
		// https://pkg.go.dev/github.com/go-vela/sdk-go/vela#BuildService.Restart
		b, _, err := client.Build.Restart(repo.GetOrg(), repo.GetName(), build.GetNumber())
		if err != nil {
			return fmt.Errorf("unable to restart build %s/%d: %w", repo.GetFullName(), build.GetNumber(), err)
		}

		// set map value for status checking
		rBMap[repo] = b.GetNumber()

		logrus.Infof("new build created %s/%d", repo.GetFullName(), b.GetNumber())
	}

	// early exit if reporting back is not enabled
	if !p.Build.Report || len(rBMap) == 0 {
		return nil
	}

	err = p.Report(client, rBMap)
	if err != nil {
		return err
	}

	return nil
}

// Report is a plugin method that checks the build statuses of all the builds kicked off from the plugin.
// It will continue to check the statuses on 30 second intervals until the timeout is reached.
func (p *Plugin) Report(client *vela.Client, rBMap map[*api.Repo]int) error {
	logrus.Info("waiting for 30 seconds to check status of downstream builds...")
	// sleep to allow for all restart processing
	time.Sleep(30 * time.Second)

	// set timeout
	timeout := time.Now().Add(p.Build.Timeout).Unix()

	successMap := make(map[int]bool)

	for time.Now().Unix() < timeout {
		logrus.Debug("checking build statuses of downstream builds...")

		for r, num := range rBMap {
			if successMap[num] {
				continue
			}

			build, _, err := client.Build.Get(r.GetOrg(), r.GetName(), num)
			if err != nil {
				return fmt.Errorf("unable to get build %s/%d: %w", r.GetFullName(), num, err)
			}

			if contains(p.Build.TargetStatus, build.GetStatus()) {
				successMap[num] = true
			} else if strings.EqualFold(build.GetStatus(), constants.StatusRunning) || strings.EqualFold(build.GetStatus(), constants.StatusPending) {
				continue
			} else {
				return fmt.Errorf("triggered build %s/%d returned %s status, exiting", r.GetFullName(), num, build.GetStatus())
			}
		}

		if len(successMap) == len(rBMap) {
			logrus.Info("all builds matched desired status")
			return nil
		}

		logrus.Info("sleeping for 30 seconds to check build statuses...")

		time.Sleep(30 * time.Second)
	}

	return fmt.Errorf("timeout while awaiting downstream build statuses")
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
