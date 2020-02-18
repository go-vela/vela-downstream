// Copyright (c) 2020 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package main

import (
	"fmt"
	"strings"

	"github.com/go-vela/types/library"

	"github.com/sirupsen/logrus"
)

// Repo represents the plugin configuration for Repo information.
type Repo struct {
	// list of Vela repos to trigger a build for
	Names []string
}

// Parse verifies the Repo is properly configured.
func (r *Repo) Parse() ([]*library.Repo, error) {
	// create new repos type to store parsed repos
	repos := []*library.Repo{}

	for _, name := range r.Names {
		// create new repo type to store parsed repo information
		repo := new(library.Repo)

		// split the repo on / to account for org/repo as input
		parts := strings.Split(name, "/")
		if len(parts) != 2 {
			return nil, fmt.Errorf("unable to parse repo on /: %s", name)
		}

		// set the org field for the repo
		repo.SetOrg(parts[0])
		// set the name field for the repo
		repo.SetName(parts[1])

		// check if a branch was provided with org/repo@branch
		if strings.Contains(parts[1], "@") {
			// split the remaining repo on @ to account for repo@branch as input
			parts = strings.Split(parts[1], "@")
			if len(parts) != 2 {
				return nil, fmt.Errorf("unable to parse repo on @: %s", name)
			}

			repo.SetName(parts[0])
			repo.SetBranch(parts[1])
		}

		// check if a branch was parsed from the input
		if len(repo.GetBranch()) == 0 {
			// set the default branch from the provided input
			repo.SetBranch("master")
		}

		// set the full name for the repo
		repo.SetFullName(
			fmt.Sprintf("%s/%s", repo.GetOrg(), repo.GetName()),
		)

		// add the parsed repo to our list of repos
		repos = append(repos, repo)
	}

	return repos, nil
}

// Validate verifies the Repo is properly configured.
func (r *Repo) Validate() error {
	logrus.Trace("validating repo configuration")

	// verify repo names are provided
	if len(r.Names) == 0 {
		return fmt.Errorf("no repo names provided")
	}

	for _, repo := range r.Names {
		if !strings.Contains(repo, "/") {
			return fmt.Errorf("invalid <org>/<repo> name provided: %s", repo)
		}

		if strings.Count(repo, "/") > 1 {
			return fmt.Errorf("invalid <org>/<repo> name provided: %s", repo)
		}
	}

	return nil
}
