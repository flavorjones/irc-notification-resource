ARTIFACTS_DIR=artifacts

default: artifacts

all: artifacts docker

artifacts: check in out

%: cmd/%
	go get -d ./cmd/$@
	go build -o $(ARTIFACTS_DIR)/$@ ./cmd/$@

docker: Dockerfile
	docker build -t irc-notifications-resource .

clean:
	rm -rf $(ARTIFACTS_DIR)
