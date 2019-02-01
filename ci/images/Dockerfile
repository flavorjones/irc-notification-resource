FROM golang:1.11
LABEL maintainer=mike.dalessio@gmail.com

COPY go.mod irc-notification-resource/
COPY go.sum irc-notification-resource/

RUN cd irc-notification-resource && go get -u
RUN cd irc-notification-resource && go get github.com/onsi/ginkgo/ginkgo
