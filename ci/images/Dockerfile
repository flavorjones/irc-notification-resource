FROM golang:1.11
LABEL maintainer=mike.dalessio@gmail.com

RUN go get github.com/onsi/ginkgo/ginkgo
RUN go get github.com/onsi/gomega/...

COPY go.mod irc-notification-resource/
COPY go.sum irc-notification-resource/
RUN cd irc-notification-resource && go mod download
