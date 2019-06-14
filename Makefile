
# can specify V=1 on the line with `make` to get verbose output
V ?= 0
Q = $(if $(filter 1,$V),,@)

GOPROXY ?= https://proxy.golang.org

.DEFAULT_GOAL := help

deps:  ## download dependencies
	$Q GOPROXY=$(GOPROXY) go mod download
	$Q GOPROXY=$(GOPROXY) go get github.com/golangci/golangci-lint/cmd/golangci-lint@v1.17.1
	$Q GOPROXY=$(GOPROXY) go get golang.org/x/tools/cmd/goimports

generate:  ## run a go generate
	$Q GOPROXY=$(GOPROXY) go generate ./...

test:  ## run go test
	$Q GOPROXY=$(GOPROXY) go test ./...

vet:  deps ## run various linters and vetters
	$Q goimports -w -local github.com/gsmcwhirter/discord-bot-lib  .
	$Q gofmt -s -w .
	$Q golangci-lint run -E golint,gosimple,staticcheck ./...
	$Q golangci-lint run -E deadcode,depguard,errcheck,gocritic,gofmt,goimports,gosec,govet,ineffassign,nakedret,prealloc,structcheck,typecheck,unconvert,varcheck ./...

help:  ## Show the help message
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}' ./Makefile