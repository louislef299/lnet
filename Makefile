.DEFAULT_GOAL := default
.PHONY: docs

default: lint test
	@echo "Building binary for your machine..."
	@go build

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

update:
	go get -u
	go mod tidy
	go mod vendor

clean:
	@rm -rf lnet