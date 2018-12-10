ARTIFACTS_DIR=artifacts

default: test

all: test docker

test: artifacts integration_test.sh
	cd cmd/out && go test
	./integration_test.sh

artifacts: $(ARTIFACTS_DIR)/check $(ARTIFACTS_DIR)/in $(ARTIFACTS_DIR)/out

$(ARTIFACTS_DIR)/check: cmd/check/check.go
	go get -d ./cmd/$(shell basename $@)
	go build -o $@ ./cmd/$(shell basename $@)

$(ARTIFACTS_DIR)/in: cmd/in/in.go
	go get -d ./cmd/$(shell basename $@)
	go build -o $@ ./cmd/$(shell basename $@)

$(ARTIFACTS_DIR)/out: cmd/out/out.go
	go get -d ./cmd/$(shell basename $@)
	go build -o $@ ./cmd/$(shell basename $@)

docker: Dockerfile
	docker build -t flavorjones/irc-notification-resource .

docker-push: docker
	docker push flavorjones/irc-notification-resource

clean:
	rm -rf $(ARTIFACTS_DIR)

.PHONY: default all test artifacts docker docker-push clean
