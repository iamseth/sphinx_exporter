VERSION := 0.0.2
LDFLAGS := -X main.Version=$(VERSION)
GOFLAGS := -ldflags "$(LDFLAGS) -s -w"
GOARCH ?= $(subst x86_64,amd64,$(patsubst i%86,386,$(shell uname -m)))


all: clean test build

deps:
	@go get -u github.com/golang/dep/cmd/dep
	@dep ensure
build:
	@mkdir -p ./dist
	@export CGO_ENABLED=0; export GOOS=linux; go build $(GOFLAGS) -o ./dist/sphinx_exporter.linux-${GOARCH}
	@export CGO_ENABLED=0; export GOOS=darwin; go build $(GOFLAGS) -o ./dist/sphinx_exporter.darwin-${GOARCH}

test:
	@go test $$(go list ./... | grep -v /vendor/)

clean:
	@rm -rf ./dist

docker: all
	@docker build -t "iamseth/sphinx_exporter:${VERSION}" .
	@docker tag iamseth/sphinx_exporter:${VERSION} iamseth/sphinx_exporter:latest

.PHONY: all deps build test clean docker

