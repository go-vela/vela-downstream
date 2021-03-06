// Copyright (c) 2021 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package main

import (
	"testing"
)

func TestDownstream_Plugin_Exec_Error(t *testing.T) {
	// setup types
	p := &Plugin{
		Config: &Config{
			Server: "http://vela.localhost.com",
			Token:  "superSecretVelaToken",
		},
		Repo: &Repo{
			Branch: "master",
			Names:  []string{"go-vela/hello-world@master"},
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
		Config: &Config{
			Server: "http://vela.localhost.com",
			Token:  "superSecretVelaToken",
		},
		Repo: &Repo{
			Branch: "master",
			Names:  []string{"go-vela/hello-world@master"},
		},
	}

	err := p.Validate()
	if err != nil {
		t.Errorf("Validate returned err: %v", err)
	}
}

func TestDownstream_Plugin_Validate_NoConfig(t *testing.T) {
	// setup types
	p := &Plugin{
		Config: &Config{},
		Repo: &Repo{
			Branch: "master",
			Names:  []string{"go-vela/hello-world@master"},
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
