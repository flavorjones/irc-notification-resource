ARTIFACTS_DIR=artifacts

clean:
	rm -rf $(ARTIFACTS_DIR)

all: check in out

%: cmd/%
	go get -d ./cmd/$@
	go build -o $(ARTIFACTS_DIR)/$@ ./cmd/$@
