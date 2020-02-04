include versions.mk
ROOT := $(PWD)
GOPATH ?= $(ROOT)/../..

.PHONY: allall
allall: all docker-build-server docker-publish-server

.PHONY: all
all: clean-client clean-server lint build test

.PHONY: build
build: linux-client linux-server

linux-%: %
	GOPATH=$(GOPATH) GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build $([ -w $(go env GOROOT) ] && echo "-i") -o $(ROOT)/docker/filestore-$</filestore-$<  filestore/$<

darwin-%: %
	GOPATH=$(GOPATH) GOOS=darwin GOARCH=amd64 CGO_ENABLED=0 go build $([ -w $(go env GOROOT) ] && echo "-i") -o $(ROOT)/docker/filestore-$</filestore-$<  filestore/$<

.PHONY: darwin
darwin: darwin-client darwin-server

deps-lint:
	@which golangci-lint > /dev/null || curl -sSL https://install.goreleaser.com/github.com/golangci/golangci-lint.sh | sh -s -- -b ${GOPATH}/bin v1.23.1

lint: deps-lint
	golangci-lint run

clean-%: %
	@rm -f $(ROOT)/docker/filestore-$</filestore-$<

.PHONY: test
test: test-client test-server

test-%: %
	GOPATH=$(GOPATH) go test -cover -v filestore/$(<F)/...

.PHONY: docker-build
docker-build: docker-build-server docker-build-client

docker-build-%: %
	docker build -t $(IMAGE_REPO)/filestore-$<:$($<_VERSION) --force-rm=true --no-cache=true --pull=true -f $(ROOT)/docker/filestore-$(<F)/Dockerfile $(ROOT)/docker/filestore-$(<F)

docker-publish-%: %
	docker push $(IMAGE_REPO)/filestore-$<:$($<_VERSION)