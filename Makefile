ARTIFACTS_DIR=artifacts

default: artifacts

all: artifacts docker

artifacts: $(ARTIFACTS_DIR)/check $(ARTIFACTS_DIR)/in $(ARTIFACTS_DIR)/out

$(ARTIFACTS_DIR)/%: cmd/%
	go get -d ./cmd/$(shell basename $@)
	go build -o $@ ./cmd/$(shell basename $@)

docker: Dockerfile
	docker build -t irc-notifications-resource .

clean:
	rm -rf $(ARTIFACTS_DIR)
