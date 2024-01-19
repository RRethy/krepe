.PHONY: build test test-krepe test-jsonpatch clean tidy

BINARY_NAME=krepe
BINARY_PATH=./krepe/$(BINARY_NAME)
GOPATH=$(shell go env GOPATH)
GOCACHE=$(GOPATH)/pkg/mod
GOBUILD=go build
GOCLEAN=go clean
GOTEST=go test
GOMOD=go mod

build:
	$(GOBUILD) -C krepe -o $(BINARY_NAME) -v

test:
	$(GOTEST) -C krepe ./...
	$(GOTEST) -C jsonpatch ./...

test-krepe:
	$(GOTEST) -C krepe ./...

test-jsonpatch:
	$(GOTEST) -C jsonpatch ./...

clean:
	$(GOCLEAN)
	rm -f $(BINARY_PATH)

tidy:
	$(GOMOD) -C krepe tidy
	$(GOMOD) -C jsonpatch tidy
