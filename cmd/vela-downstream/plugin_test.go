// SPDX-License-Identifier: Apache-2.0

package main

import (
	"testing"

	"github.com/go-vela/types/constants"
)

func TestDownstream_Plugin_Exec_Error(t *testing.T) {
	// setup types
	p := &Plugin{
		Build: &Build{
			Branch: "main",
			Event:  constants.EventPush,
			Status: []string{constants.StatusSuccess},
		},
		Config: &Config{
			Server: "http://vela.localhost.com",
			Token:  "superSecretVelaToken",
		},
		Repo: &Repo{
			Names: []string{"go-vela/hello-world@main"},
		},
	}

	err := p.Exec()
	if err == nil {
		t.Errorf("Exec should have returned err")
	}
}

func TestDownstream_Plugin_Validate(t *testing.T) {
	// setup types
	p := &Plugin{
		Build: &Build{
			Branch: "main",
			Event:  constants.EventPush,
			Status: []string{constants.StatusSuccess},
		},
		Config: &Config{
			Server: "http://vela.localhost.com",
			Token:  "superSecretVelaToken",
		},
		Repo: &Repo{
			Names: []string{"go-vela/hello-world@main"},
		},
	}

	err := p.Validate()
	if err != nil {
		t.Errorf("Validate returned err: %v", err)
	}
}

func TestDownstream_Plugin_Validate_NoBuild(t *testing.T) {
	// setup types
	p := &Plugin{
		Build: &Build{},
		Config: &Config{
			Server: "http://vela.localhost.com",
			Token:  "superSecretVelaToken",
		},
		Repo: &Repo{
			Names: []string{"go-vela/hello-world@main"},
		},
	}

	err := p.Validate()
	if err == nil {
		t.Errorf("Validate should have returned err")
	}
}

func TestDownstream_Plugin_Validate_NoConfig(t *testing.T) {
	// setup types
	p := &Plugin{
		Build: &Build{
			Branch: "main",
			Event:  constants.EventPush,
			Status: []string{constants.StatusSuccess},
		},
		Config: &Config{},
		Repo: &Repo{
			Names: []string{"go-vela/hello-world@main"},
		},
	}

	err := p.Validate()
	if err == nil {
		t.Errorf("Validate should have returned err")
	}
}

func TestDownstream_Plugin_Validate_NoRepo(t *testing.T) {
	// setup types
	p := &Plugin{
		Build: &Build{
			Branch: "main",
			Event:  constants.EventPush,
			Status: []string{constants.StatusSuccess},
		},
		Config: &Config{
			Server: "http://vela.localhost.com",
			Token:  "superSecretVelaToken",
		},
		Repo: &Repo{},
	}

	err := p.Validate()
	if err == nil {
		t.Errorf("Validate should have returned err")
	}
}
