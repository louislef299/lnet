# syntax=docker/dockerfile:1

FROM fedora:38

RUN dnf -y update && dnf clean all
RUN dnf -y install go-1.20.8-1.fc38 curl-8.0.1-4.fc38 git-2.41.0-1.fc38 cmake-3.27.4-7.fc38 && dnf clean all

# Install golangci-lint
SHELL ["/bin/bash", "-o", "pipefail", "-c"]
RUN curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh \
    | sh -s -- -b "$(go env GOPATH)"/bin v1.54.2

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
