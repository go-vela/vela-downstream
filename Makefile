# Copyright (c) 2020 Target Brands, Inc. All rights reserved.
#
# Use of this source code is governed by the LICENSE file in this repository.

.PHONY: build
build: binary-build

.PHONY: run
run: build docker-build docker-run

.PHONY: test
test: build docker-build docker-example

#################################
######      Go clean       ######
#################################

.PHONY: clean
clean:

	@go mod tidy
	@go vet ./...
	@go fmt ./...
	@echo "I'm kind of the only name in clean energy right now"

#################################
######    Build Binary     ######
#################################

.PHONY: binary-build
binary-build:

	GOOS=linux CGO_ENABLED=0 go build -o release/vela-downstream github.com/go-vela/vela-downstream/cmd/vela-downstream

#################################
######    Docker Build     ######
#################################

.PHONY: docker-build
docker-build:

	docker build --no-cache -t vela-artifactory:local .

#################################
######     Docker Run      ######
#################################

.PHONY: docker-run
docker-run:

	docker run --rm \
		-e DOWNSTREAM_SERVER \
		-e DOWNSTREAM_TOKEN \
		-e PARAMETER_LOG_LEVEL \
		-e PARAMETER_BRANCH \
		-e PARAMETER_REPOS \
		vela-downstream:local

.PHONY: docker-example
docker-example:

	docker run --rm \
		-e DOWNSTREAM_SERVER \
		-e DOWNSTREAM_TOKEN \
		-e PARAMETER_REPOS \
		vela-downstream:local
