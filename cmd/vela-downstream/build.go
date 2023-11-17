// SPDX-License-Identifier: Apache-2.0

package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/go-vela/types/constants"
	"github.com/sirupsen/logrus"
)

// Build represents the plugin configuration for Build information.
type Build struct {
	// branch to trigger a build for the repo
	Branch string
	// event to trigger a build for the repo
	Event string
	// status to trigger a build for the repo
	Status []string
	// report determines whether to report back the build statuses
	Report bool
	// target status for triggered builds
	TargetStatus []string
	// timeout for waiting on triggered builds
	Timeout time.Duration
	// continue through repo list if build is not found to restart
	Continue bool
}

// Validate verifies the Build is properly configured.
func (b *Build) Validate() error {
	logrus.Trace("validating build configuration")

	// verify build branch is provided
	if len(b.Branch) == 0 {
		logrus.Debug("no build branch provided for filtering")
	}

	// verify build event is provided
	if len(b.Event) == 0 {
		return fmt.Errorf("no build event provided")
	}

	// set timeout
	if b.Timeout > (90 * time.Minute) {
		logrus.Info("timeout set too high. Using 90 minutes...")

		b.Timeout = 90 * time.Minute
	}

	// create a list of valid events for a build
	validEvents := []string{
		constants.EventComment,
		constants.EventDeploy,
		constants.EventPull,
		constants.EventPush,
		"schedule",
		constants.EventTag,
	}

	// verify the build event provided is valid
	if !contains(validEvents, b.Event) {
		return fmt.Errorf("invalid build event provided: %s", b.Event)
	}

	// verify build status is provided
	if len(b.Status) == 0 {
		return fmt.Errorf("no build status provided")
	}

	// create a list of valid statuses for a build
	validStatuses := []string{
		constants.StatusCanceled,
		constants.StatusError,
		constants.StatusFailure,
		constants.StatusKilled,
		constants.StatusPending,
		constants.StatusRunning,
		constants.StatusSuccess,
	}

	// iterate through the build statuses provided
	for _, status := range b.Status {
		// verify the build status provided is valid
		if !contains(validStatuses, status) {
			return fmt.Errorf("invalid build status provided: %s", status)
		}
	}

	return nil
}

// contains checks if the provided input string is found in the given list of
// strings. If the input string is not found, then the function returns false.
func contains(list []string, input string) bool {
	// iterate through the list of strings
	for _, item := range list {
		// check if the item matches the input
		if strings.EqualFold(item, input) {
			return true
		}
	}

	return false
}
