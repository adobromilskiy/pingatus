.DEFAULT_GOAL := build

export GOBIN = $(shell pwd)/bin
export GOFLAGS = -mod=vendor -race

.PHONY: fmt
fmt:
	@gofumpt -l -w .

.PHONY: vet
vet: fmt
	@go vet ./...
	@staticcheck ./...
	@shadow ./...

.PHONY: lint
lint: vet
	@deadcode ./...
	@golangci-lint run ./...

.PHONY: build
build: vet
	@go install ./cmd/... && echo "Build successful"

.PHONY: test
test: vet
	@go test -cover -v -coverprofile cover.out ./...

.PHONY: cover
cover: test
	@go tool cover -html cover.out
