package main

import (
	"fmt"
	"strings"

	"github.com/go-vela/types/constants"
	"github.com/sirupsen/logrus"
)

type Status []string

// Validate verifies the Config is properly configured.
func (s *Status) Validate() error {
	logrus.Trace("validating config configuration")

	acceptableStatus := Status{
		constants.StatusError,
		constants.StatusFailure,
		constants.StatusKilled,
		constants.StatusCanceled,
		constants.StatusPending,
		constants.StatusRunning,
		constants.StatusSuccess,
	}

	for _, v := range *s {
		if !contains(acceptableStatus, v) {
			return fmt.Errorf("Status %s not of acceptable type %s", s, acceptableStatus)
		}
	}

	return nil
}

// contains checks to see if a []string contains a string.
func contains(acceptableStatus Status, s string) bool {
	for _, v := range acceptableStatus {
		if strings.EqualFold(v, s) {
			return true
		}
	}
	return false
}
