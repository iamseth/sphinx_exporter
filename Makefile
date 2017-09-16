VERSION := 0.0.1
LDFLAGS := -X main.Version=$(VERSION)
GOFLAGS := -ldflags "$(LDFLAGS) -s -w"
GOOS ?= $(shell uname | tr A-Z a-z)
GOARCH ?= $(subst x86_64,amd64,$(patsubst i%86,386,$(shell uname -m)))
SUFFIX ?= $(GOOS)-$(GOARCH)
BINARY := sphinx_exporter.$(SUFFIX)


all: clean test build

deps:
	@go get -u github.com/golang/dep/cmd/dep
	@dep ensure
build:
	@mkdir -p ./dist
	@go build $(GOFLAGS) -o ./dist/${BINARY}

test:
	@go test $$(go list ./... | grep -v /vendor/)

clean:
	@rm -rf ./dist


.PHONY: all deps build test clean
