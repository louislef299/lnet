.DEFAULT_GOAL := default
.PHONY: docs clean assembly

BUILD_DIR= .build
BINARY_NAME= lnet

GOBIN = ${HOME}/go/bin
GOTRACEBACK = 'crash'
GOVERSION= $(shell go version | awk '{print $$3}')
GOFLAGS= -s -w -X 'github.com/louislef299/lnet/pkg/version.Version=$(shell cat version.txt)' \
-X 'github.com/louislef299/lnet/pkg/version.BuildOS=$(shell go env GOOS)' \
-X 'github.com/louislef299/lnet/pkg/version.BuildArch=$(shell go env GOARCH)' \
-X 'github.com/louislef299/lnet/pkg/version.GoVersion=$(GOVERSION)' \
-X 'github.com/louislef299/lnet/pkg/version.BuildTime=$(shell date)' \
-X 'github.com/louislef299/lnet/pkg/version.CommitHash=$(shell git rev-parse --short HEAD)'

default: lint test $(BINARY_NAME)
	@echo "Run './$(BINARY_NAME) -h' to get started"

local: lint test $(BINARY_NAME)
	@echo "GOVERSION: $(GOVERSION)"
	@echo "Moving binary to $(GOBIN)"
	@mv lnet $(GOBIN)

$(BINARY_NAME):
	@echo "Building $(BINARY_NAME) binary for your machine..."
	@go build -mod vendor -ldflags="$(GOFLAGS)" -o $(BINARY_NAME)

# creates the command documentation
docs:
	@echo "Generating command documentation in docs/cmd"
	@go run main.go docs --dir docs/cmds

test:
	@echo "Running tests..."
	@go test -v -race -cover ./...

lint:
	@echo "Linting..."
	@golangci-lint run

update:
	go mod tidy
	go mod vendor

login:
	@gh auth status || gh auth login --git-protocol https -w -s repo,repo_deployment,workflow

release: lint test login
	@goreleaser check
	@GITHUB_TOKEN=$(shell gh auth token) GOVERSION=$(GOVERSION) \
	 goreleaser release --clean

container:
	docker buildx build -f ./Dockerfile -t lnet .

$(BUILD_DIR):
	@mkdir -p $(BUILD_DIR)

assembly: $(BINARY_NAME) $(BUILD_DIR)
	@echo "Dumping assembly output to $(BUILD_DIR)/$(BINARY_NAME).asm..."
	@go tool objdump $(BINARY_NAME) > $(BUILD_DIR)/$(BINARY_NAME).asm

clean:
	@rm -rf lnet $(BUILD_DIR) dist