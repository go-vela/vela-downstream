// SPDX-License-Identifier: Apache-2.0

package main

import (
	"fmt"
	"net/url"

	"github.com/go-vela/sdk-go/vela"

	"github.com/sirupsen/logrus"
)

// Config represents the plugin configuration for Config information.
type Config struct {
	// Vela server to interact with
	Server string
	// user token to authenticate with the Vela server
	Token string
	// depth of builds search in downstream repo
	Depth int
	// the app name utilizing this config
	AppName string
	// the app version utilizing this config
	AppVersion string
}

// New creates a Vela client for triggering builds.
func (c *Config) New() (*vela.Client, error) {
	logrus.Trace("creating new Vela client from plugin configuration")

	// create the app string
	appID := fmt.Sprintf("%s; %s", c.AppName, c.AppVersion)

	// create Vela client from configuration
	client, err := vela.NewClient(c.Server, appID, nil)
	if err != nil {
		return nil, err
	}

	// check if a token is provided for authentication
	if len(c.Token) > 0 {
		logrus.Debugf("setting authentication token for Vela")

		// set the token for authentication in the Vela client
		client.Authentication.SetPersonalAccessTokenAuth(c.Token)
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

	// check to make sure it's a valid url
	_, err := url.ParseRequestURI(c.Server)
	if err != nil {
		return fmt.Errorf("%s is not a valid url", c.Server)
	}

	// verify token is provided
	if len(c.Token) == 0 {
		return fmt.Errorf("no config token provided")
	}

	return nil
}
