// Copyright (c) 2020 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package main

import (
	"testing"
)

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
