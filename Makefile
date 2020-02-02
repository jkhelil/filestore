include versions.mk
ROOT := $(PWD)
GOPATH ?= $(ROOT)/../..

.PHONY: allall
allall: all build-docker-server

.PHONY: all
all: clean-client clean-server build test

.PHONY: build
build: linux-client linux-server

linux-%: %
	GOPATH=$(GOPATH) GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build $([ -w $(go env GOROOT) ] && echo "-i") -o $(ROOT)/docker/filestore-$</filestore-$<  filestore/$<

darwin-%: %
	GOPATH=$(GOPATH) GOOS=darwin GOARCH=amd64 CGO_ENABLED=0 go build $([ -w $(go env GOROOT) ] && echo "-i") -o $(ROOT)/docker/filestore-$</filestore-$<  filestore/$<

.PHONY: darwin
darwin: darwin-client darwin-server

clean-%: %
	@rm -f $(ROOT)/docker/filestore-$</filestore-$<

.PHONY: test
test: test-client test-server

test-%: %
	GOPATH=$(GOPATH) go test -cover -v filestore/$(<F)/...

build-docker-%: %
	docker build -t $(IMAGE_REPO)/filestore-$<:$($<_VERSION) --force-rm=true --no-cache=true --pull=true -f $(ROOT)/docker/filestore-$(<F)/Dockerfile $(ROOT)/docker/filestore-$(<F)