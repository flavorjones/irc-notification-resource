ARTIFACTS_DIR=artifacts

default: test

all: test docker

test: unit integration

unit: artifacts
	ginkgo -r

integration: artifacts integration_test.sh
	./integration_test.sh

artifacts: $(ARTIFACTS_DIR)/check $(ARTIFACTS_DIR)/in $(ARTIFACTS_DIR)/out

$(ARTIFACTS_DIR)/check: cmd/check/check.go
	go get ./cmd/$(shell basename $@)
	go build -o $@ ./cmd/$(shell basename $@)

$(ARTIFACTS_DIR)/in: cmd/in/in.go
	go get ./cmd/$(shell basename $@)
	go build -o $@ ./cmd/$(shell basename $@)

$(ARTIFACTS_DIR)/out: cmd/out/out.go pkg/irc/irc.go
	go get ./cmd/$(shell basename $@)
	go build -o $@ ./cmd/$(shell basename $@)

docker: Dockerfile
	docker build -t flavorjones/irc-notification-resource .

docker-push: docker
	docker push flavorjones/irc-notification-resource

clean:
	rm -rf $(ARTIFACTS_DIR)

.PHONY: default all test artifacts docker docker-push clean unit
