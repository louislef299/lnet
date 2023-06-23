.DEFAULT_GOAL := default
.PHONY: docs clean assembly

BUILD_DIR= .build
BINARY_NAME= lnet

GOBIN = ${HOME}/go/bin
GOTRACEBACK = 'crash'
GOFLAGS= -s -w -X 'github.com/louislef299/lnet/pkg/version.Version={{.Version}}' \
-X 'github.com/louislef299/lnet/pkg/version.BuildOS={{.Runtime.Goos}}' \
-X 'github.com/louislef299/lnet/pkg/version.BuildArch={{.Runtime.Goarch}}' \
-X 'github.com/louislef299/lnet/pkg/version.GoVersion={{.Env.GOVERSION}}' \
-X 'github.com/louislef299/lnet/pkg/version.BuildTime={{.Date}}' \
-X 'github.com/louislef299/lnet/pkg/version.CommitHash={{.ShortCommit}}'

default: lint test binary

$(BINARY_NAME):
	@echo "Building $(BINARY_NAME) binary for your machine..."
	@go build -ldflags="$(GOFLAGS)"

# creates the command documentation
docs:
	@echo "Generating command documentation in docs/cmd"
	@go run main.go docs --dir docs/cmds

test:
	@echo "Running tests..."
	@go test -v -race ./...

lint:
	@echo "Linting..."
	@golangci-lint run

local: default
	@echo "Moving binary to $(GOBIN)"
	@mv lnet $(GOBIN)

update:
	go get -u
	go mod tidy
	go mod vendor

$(BUILD_DIR):
	@mkdir -p $(BUILD_DIR)

assembly: $(BINARY_NAME) $(BUILD_DIR)
	@echo "Dumping assembly output to $(BUILD_DIR)/$(BINARY_NAME).asm..."
	@go tool objdump $(BINARY_NAME) > $(BUILD_DIR)/$(BINARY_NAME).asm

clean:
	@rm -rf lnet $(BUILD_DIR)