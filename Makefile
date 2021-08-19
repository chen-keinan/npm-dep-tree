SHELL := /bin/bash

GOCMD=go
GOMOD=$(GOCMD) mod
GOBUILD=$(GOCMD) build
GOLINT=${GOPATH}/bin/golangci-lint
GORELEASER=/usr/local/bin/goreleaser
GOIMPI=${GOPATH}/bin/impi
GOTEST=$(GOCMD) test

all:
	$(info  "completed running make file for npm dependency resolver")
fmt:
	@go fmt ./...
lint:
	./lint.sh
tidy:
	$(GOMOD) tidy -v
test:
	@go get github.com/golang/mock/mockgen@latest
	@go install -v github.com/golang/mock/mockgen && export PATH=$GOPATH/bin:$PATH;
	@go generate ./...
	$(GOTEST) ./... -coverprofile cp.out
build:
	$(GOBUILD) -v

build_docker:
	docker build -t gcr.io/snyk-main/homebase-upstream-consumer:${CIRCLE_SHA1} .
	docker push gcr.io/snyk-main/homebase-upstream-consumer:${CIRCLE_SHA1}


.PHONY: install-req fmt test lint build ci build-binaries tidy imports
