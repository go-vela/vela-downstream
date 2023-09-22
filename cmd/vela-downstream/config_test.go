// SPDX-License-Identifier: Apache-2.0

package main

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/go-vela/sdk-go/vela"
)

func TestDownstream_Config_New(t *testing.T) {
	// setup types
	c := &Config{
		Server: "http://vela.localhost.com",
		Token:  "superSecretVelaToken",
	}

	appID := fmt.Sprintf("%s; %s", c.AppName, c.AppVersion)

	want, err := vela.NewClient(c.Server, appID, nil)
	if err != nil {
		t.Errorf("Unable to create new Vela client: %v", err)
	}

	want.Authentication.SetPersonalAccessTokenAuth(c.Token)

	got, err := c.New()
	if err != nil {
		t.Errorf("New returned err: %v", err)
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("New is %v, want %v", got, want)
	}
}

func TestDownstream_Config_New_NoConfig(t *testing.T) {
	// setup types
	c := &Config{}

	got, err := c.New()
	if err == nil {
		t.Errorf("New should have returned err")
	}

	if got != nil {
		t.Errorf("New is %v, want nil", got)
	}
}

func TestDownstream_Config_Validate(t *testing.T) {
	// setup types
	c := &Config{
		Server: "http://vela.localhost.com",
		Token:  "superSecretVelaToken",
	}

	err := c.Validate()
	if err != nil {
		t.Errorf("Validate returned err: %v", err)
	}
}

func TestDownstream_Config_Validate_NoServer(t *testing.T) {
	// setup types
	c := &Config{
		Token: "superSecretVelaToken",
	}

	err := c.Validate()
	if err == nil {
		t.Errorf("Validate should have returned err")
	}
}

func TestDownstream_Config_Validate_InvalidServer(t *testing.T) {
	// setup types
	c := &Config{
		Server: "vela.server",
		Token:  "superSecretVelaToken",
	}

	err := c.Validate()
	if err == nil {
		t.Errorf("Validate should have returned err")
	}
}

func TestDownstream_Config_Validate_NoToken(t *testing.T) {
	// setup types
	c := &Config{
		Server: "http://vela.localhost.com",
	}

	err := c.Validate()
	if err == nil {
		t.Errorf("Validate should have returned err")
	}
}
