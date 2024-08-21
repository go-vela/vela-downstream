// SPDX-License-Identifier: Apache-2.0

package main

import (
	"reflect"
	"testing"

	api "github.com/go-vela/server/api/types"
)

func TestDownstream_Repo_Parse(t *testing.T) {
	// setup types
	r := &Repo{
		Names: []string{"go-vela/hello-world@test", "go-vela/hello-world"},
	}

	r1 := new(api.Repo)
	r1.SetOrg("go-vela")
	r1.SetName("hello-world")
	r1.SetFullName("go-vela/hello-world")
	r1.SetBranch("test")

	r2 := new(api.Repo)
	r2.SetOrg("go-vela")
	r2.SetName("hello-world")
	r2.SetFullName("go-vela/hello-world")
	r2.SetBranch("main")

	want := []*api.Repo{r1, r2}

	// run test
	got, err := r.Parse("main")
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
		Names: []string{"go-vela/hello-world/main"},
	}

	// run test
	got, err := r.Parse("main")
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
		Names: []string{"go-vela/hello-world@main@"},
	}

	// run test
	got, err := r.Parse("main")
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
		Names: []string{"go-vela/hello-world@main"},
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
		Names: []string{"go-vela/hello-world/main"},
	}

	// run test
	err := r.Validate()
	if err == nil {
		t.Errorf("Validate should have returned err")
	}
}
