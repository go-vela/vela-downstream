// Copyright (c) 2020 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package main

import (
	"testing"
)

func TestDownstream_Repo_Validate(t *testing.T) {
	// setup types
	r := &Repo{
		Names: []string{"vela/hello-world"},
	}

	err := r.Validate()
	if err != nil {
		t.Errorf("Validate returned err: %v", err)
	}
}

func TestDownstream_Repo_Validate_NoNames(t *testing.T) {
	// setup types
	r := &Repo{}

	err := r.Validate()
	if err == nil {
		t.Errorf("Validate should have returned err")
	}
}
