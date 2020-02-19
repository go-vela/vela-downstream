// Copyright (c) 2020 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package main

import (
	"fmt"

	"github.com/go-vela/sdk-go/vela"

	"github.com/sirupsen/logrus"
)

// Config represents the plugin configuration for Config information.
type Config struct {
	// Vela server to interact with
	Server string
	// user token to authenticate with the Vela server
	Token string
}

// New creates a Vela client for triggering builds.
func (c *Config) New() (*vela.Client, error) {
	logrus.Trace("creating new Vela client from plugin configuration")

	// create Vela client from configuration
	client, err := vela.NewClient(c.Server, nil)
	if err != nil {
		return nil, err
	}

	// check if a token is provided for authentication
	if len(c.Token) > 0 {
		logrus.Debugf("setting authentication token for Vela")

		// set the token for authentication in the Vela client
		client.Authentication.SetTokenAuth(c.Token)
	}

	return client, nil
}

// Validate verifies the Config is properly configured.
func (c *Config) Validate() error {
	logrus.Trace("validating config configuration")

	// verify server is provided
	if len(c.Server) == 0 {
		return fmt.Errorf("no config server provided")
	}

	// verify token is provided
	if len(c.Token) == 0 {
		return fmt.Errorf("no config token provided")
	}

	return nil
}
