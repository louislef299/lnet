.DEFAULT_GOAL := default
.PHONY: docs

default:
	go build

# creates the command documentation
docs:
	go run main.go docs --dir docs/cmds

test:
	go test -v -race ./...

lint:
	golangci-lint run

update:
	go get -u
	go mod tidy
	go mod vendor

clean:
	rm -rf lnet