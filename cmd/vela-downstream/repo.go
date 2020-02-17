// Copyright (c) 2020 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package main

import (
	"fmt"

	"github.com/sirupsen/logrus"
)

// Repo represents the plugin configuration for Repo information.
type Repo struct {
	// list of Vela repos to trigger a build for
	Names []string
}

// Validate verifies the Repo is properly configured.
func (r *Repo) Validate() error {
	logrus.Trace("validating repo configuration")

	// verify repos are provided
	if len(r.Names) == 0 {
		return fmt.Errorf("no repos provided")
	}

	return nil
}
