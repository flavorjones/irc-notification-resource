FROM golang:alpine AS builder
RUN apk --no-cache add bash ca-certificates git make

COPY . /go/src/github.com/flavorjones/irc-notification-resource
RUN cd /go/src/github.com/flavorjones/irc-notification-resource && make clean
RUN cd /go/src/github.com/flavorjones/irc-notification-resource && make test

FROM alpine:edge AS resource
COPY --from=builder /go/src/github.com/flavorjones/irc-notification-resource/artifacts /opt/resource
