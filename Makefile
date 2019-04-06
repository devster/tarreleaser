export GO111MODULE := on
export GOPROXY := https://gocenter.io

.DEFAULT_GOAL := build

# Build a dev version
build:
	go build
.PHONY: build

# Run all the tests
test:
	go test -failfast -race -coverpkg=./... -covermode=atomic -coverprofile=coverage.txt ./...
.PHONY: test

# Run all the tests and opens the coverage report
cover: test
	go tool cover -html=coverage.txt
.PHONY: cover

fmt:
	go fmt ./...
.PHONY: fmt
