// SPDX-License-Identifier: Apache-2.0

package main

import (
	"reflect"
	"testing"

	"github.com/go-vela/types/library"
)

func TestDownstream_Repo_Parse(t *testing.T) {
	// setup types
	r := &Repo{
		Names: []string{"go-vela/hello-world@test", "go-vela/hello-world"},
	}

	r1 := new(library.Repo)
	r1.SetOrg("go-vela")
	r1.SetName("hello-world")
	r1.SetFullName("go-vela/hello-world")
	r1.SetBranch("test")

	r2 := new(library.Repo)
	r2.SetOrg("go-vela")
	r2.SetName("hello-world")
	r2.SetFullName("go-vela/hello-world")
	r2.SetBranch("master")

	want := []*library.Repo{r1, r2}

	// run test
	got, err := r.Parse("master")
	if err != nil {
		t.Errorf("Parse returned err: %v", err)
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("Parse is %v, want %v", got, want)
	}
}

func TestDownstream_Repo_Parse_MultipleSlashes(t *testing.T) {
	// setup types
	r := &Repo{
		Names: []string{"go-vela/hello-world/master"},
	}

	// run test
	got, err := r.Parse("master")
	if err == nil {
		t.Errorf("Parse should have returned err")
	}

	if got != nil {
		t.Errorf("Parse is %v, want nil", got)
	}
}

func TestDownstream_Repo_Parse_MultipleAtSigns(t *testing.T) {
	// setup types
	r := &Repo{
		Names: []string{"go-vela/hello-world@master@"},
	}

	// run test
	got, err := r.Parse("master")
	if err == nil {
		t.Errorf("Parse should have returned err")
	}

	if got != nil {
		t.Errorf("Parse is %v, want nil", got)
	}
}

func TestDownstream_Repo_Validate(t *testing.T) {
	// setup types
	r := &Repo{
		Names: []string{"go-vela/hello-world@master"},
	}

	err := r.Validate()
	if err != nil {
		t.Errorf("Validate returned err: %v", err)
	}
}

func TestDownstream_Repo_Validate_NoNames(t *testing.T) {
	// setup types
	r := &Repo{}

	// run test
	err := r.Validate()
	if err == nil {
		t.Errorf("Validate should have returned err")
	}
}

func TestDownstream_Repo_Validate_NoSlash(t *testing.T) {
	// setup types
	r := &Repo{
		Names: []string{"go-vela_hello-world"},
	}

	// run test
	err := r.Validate()
	if err == nil {
		t.Errorf("Validate should have returned err")
	}
}

func TestDownstream_Repo_Validate_MultipleSlashes(t *testing.T) {
	// setup types
	r := &Repo{
		Names: []string{"go-vela/hello-world/master"},
	}

	// run test
	err := r.Validate()
	if err == nil {
		t.Errorf("Validate should have returned err")
	}
}
