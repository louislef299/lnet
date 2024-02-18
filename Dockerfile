# syntax=docker/dockerfile:1

FROM golang:1.22 as builder

ARG commithash="notset"
ENV COMMIT_HASH=$commithash

# Copy local code to the container image.
WORKDIR /go/src/lnet
COPY . .

RUN make lnet

FROM debian:buster
COPY --from=builder /go/src/lnet/lnet /lnet

ENV DEBIAN_FRONTEND=noninteractive

RUN apt-get -y update && apt-get clean all \
    && rm -rf /var/lib/apt/lists/*

EXPOSE 8080

CMD ["/lnet"]
