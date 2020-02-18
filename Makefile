# Copyright (c) 2020 Target Brands, Inc. All rights reserved.
#
# Use of this source code is governed by the LICENSE file in this repository.

build: binary-build

run: build docker-build docker-run

test: build docker-build docker-example

#################################
######      Go clean       ######
#################################

clean:

	@go mod tidy
	@go vet ./...
	@go fmt ./...
	@echo "I'm kind of the only name in clean energy right now"

#################################
######    Build Binary     ######
#################################

binary-build:

	GOOS=linux CGO_ENABLED=0 go build -o release/vela-downstream github.com/go-vela/vela-downstream/cmd/vela-downstream

#################################
######    Docker Build     ######
#################################

docker-build:

	docker build --no-cache -t vela-downstream:local .

#################################
######     Docker Run      ######
#################################

docker-run:

	docker run --rm \
		-e DOWNSTREAM_SERVER \
		-e DOWNSTREAM_TOKEN \
		-e PARAMETER_LOG_LEVEL \
		-e PARAMETER_BRANCH \
		-e PARAMETER_REPOS \
		vela-downstream:local

docker-example:

	docker run --rm \
		-e DOWNSTREAM_SERVER \
		-e DOWNSTREAM_TOKEN \
		-e PARAMETER_REPOS \
		vela-downstream:local
