FROM golang:alpine AS builder
RUN apk --no-cache add bash ca-certificates git

COPY . /go/src/github.com/flavorjones/irc-notification-resource
RUN go get -d github.com/flavorjones/irc-notification-resource/...
RUN go build -o /go/src/github.com/flavorjones/irc-notification-resource/artifacts/check github.com/flavorjones/irc-notification-resource/cmd/check
RUN go build -o /go/src/github.com/flavorjones/irc-notification-resource/artifacts/in github.com/flavorjones/irc-notification-resource/cmd/in
RUN go build -o /go/src/github.com/flavorjones/irc-notification-resource/artifacts/out github.com/flavorjones/irc-notification-resource/cmd/out
RUN cd /go/src/github.com/flavorjones/irc-notification-resource/ && ./cmd/out/test_input.sh

FROM alpine:edge AS resource
COPY --from=builder /go/src/github.com/flavorjones/irc-notification-resource/artifacts /opt/resource

