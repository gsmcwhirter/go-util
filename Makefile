PROJECT := github.com/gsmcwhirter/go-util

# can specify V=1 on the line with `make` to get verbose output
V ?= 0
Q = $(if $(filter 1,$V),,@)

GOPROXY ?= https://proxy.golang.org

.DEFAULT_GOAL := help

deps:  ## download dependencies
	$Q GOPROXY=$(GOPROXY) go mod download
	$Q GOPROXY=$(GOPROXY) go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.52.2
	$Q GOPROXY=$(GOPROXY) go install golang.org/x/tools/cmd/goimports

generate:  ## run a go generate
	$Q GOPROXY=$(GOPROXY) go generate ./...

test:  ## run go test
	$Q GOPROXY=$(GOPROXY) go test ./...

test-coverage:
	$Q GOPROXY=$(GOPROXY) go test -coverprofile=coverage.out ./...

benchmark-json:
	$Q GOPROXY=$(GOPROXY) go test -bench=. -benchmem -cpuprofile=profile_cpu.out -memprofile=profile_mem.out ./json/...
	$Q go tool pprof -svg profile_cpu.out > profile_cpu.svg
	$Q go tool pprof -svg profile_mem.out > profile_mem.svg

vet:  deps ## run various linters and vetters
	$Q bash -c 'for d in $$(go list -f {{.Dir}} ./...); do gofmt -s -w $$d/*.go; done'
	$Q bash -c 'for d in $$(go list -f {{.Dir}} ./...); do goimports -w -local $(PROJECT) $$d/*.go; done'
	$Q golangci-lint run -c .golangci.yml -E revive,gosimple,staticcheck ./...
	$Q golangci-lint run -c .golangci.yml -E asciicheck,contextcheck,depguard,durationcheck,errcheck,errname,gocritic,gofumpt,goimports,gosec,govet,ineffassign,nakedret,paralleltest,prealloc,predeclared,typecheck,unconvert,unused,whitespace ./...
	$Q golangci-lint run -c .golangci.yml -E godox ./... || true

help:  ## Show the help message
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}' ./Makefile

# golangci-lint run --build-tags="$BUILD_TAGS" -E depguard,errcheck,gocritic,gofmt,goimports,gosec,govet,ineffassign,nakedret,prealloc,typecheck,unconvert,unused "$project"