.DEFAULT_GOAL := default

default:
	go build

test:
	go test -v -race ./...

lint:
	golangci-lint run

clean:
	rm -rf lnet