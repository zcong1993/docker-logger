generate:
	@go generate ./...
.PHONY: generate

install:
	@go get -u -v github.com/golang/dep/cmd/dep
	@dep ensure

build: generate
	@echo "====> Build docker logger cli"
	@go build -o ./bin/docker-logger cmd/main.go
.PHONY: build

release:
	@echo "====> Build and release"
	@go get github.com/goreleaser/goreleaser
	@goreleaser
.PHONY: release

test:
	@go test ./...
.PHONY: test

test.cov:
	@go test ./... -coverprofile=coverage.txt -covermode=atomic
.PHONY: test.cov
