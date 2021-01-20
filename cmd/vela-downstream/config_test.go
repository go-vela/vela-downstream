// Copyright (c) 2021 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package main

import (
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

	want, err := vela.NewClient(c.Server, nil)
	if err != nil {
		t.Errorf("Unable to create new Vela client: %v", err)
	}

	want.Authentication.SetTokenAuth(c.Token)

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
