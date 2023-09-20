# syntax=docker/dockerfile:1

FROM fedora:38

RUN dnf -y update && dnf clean all
RUN dnf -y install go curl git make && dnf clean all

# Install golangci-lint
RUN "curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh \
    | sh -s -- -b $(go env GOPATH)/bin v1.54.2"

WORKDIR /lnet

# Copy dependencies and download
COPY go.mod go.sum ./
COPY vendor ./vendor

# Copy program files
COPY cmd cmd/
COPY pkg pkg/
COPY *.go Makefile version.txt ./

ARG commithash="notset"
ENV COMMIT_HASH=$commithash

EXPOSE 8080

# Build the lnet binary
RUN make lnet
CMD ["./lnet"]
