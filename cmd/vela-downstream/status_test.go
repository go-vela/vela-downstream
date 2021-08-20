// Copyright (c) 2021 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package main

import (
	"testing"

	"github.com/go-vela/types/constants"
)

func TestDownstream_Status_Validate(t *testing.T) {
	// setup types
	s := &Status{constants.StatusSuccess}

	err := s.Validate()
	if err != nil {
		t.Errorf("Validate returned err: %v", err)
	}
}

func TestDownstream_Status_Validate_Multiple(t *testing.T) {
	// setup types
	s := &Status{constants.StatusSuccess, constants.StatusCanceled}

	err := s.Validate()
	if err != nil {
		t.Errorf("Validate returned err: %v", err)
	}
}

func TestDownstream_Status_Validate_MultipleOneInvalid(t *testing.T) {
	// setup types
	s := &Status{constants.StatusSuccess, "foo"}

	err := s.Validate()
	if err == nil {
		t.Errorf("Validate should have returned err")
	}
}

func TestDownstream_Status_Validate_InvalidStatus(t *testing.T) {
	// setup types
	s := &Status{"foo"}

	err := s.Validate()
	if err == nil {
		t.Errorf("Validate should have returned err")
	}
}
