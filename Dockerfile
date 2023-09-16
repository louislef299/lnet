# syntax=docker/dockerfile:1

FROM fedora:38

MAINTAINER Louis Lefebvre <louislefebvre1999@gmail.com>

RUN dnf -y update && dnf clean all
RUN dnf -y install go curl git make && dnf clean all

# Install golangci-lint
RUN curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh \
    | sh -s -- -b $(go env GOPATH)/bin v1.54.2 

WORKDIR /lnet

# Copy dependencies and download
COPY go.mod go.sum ./
RUN go mod download

# Copy program files
COPY cmd/* cmd/
COPY pkg/* pkg/
COPY .git .git
COPY *.go Makefile version.txt ./

CMD ["sleep"]

# RUN make lnet

# CMD ["/lnet"]