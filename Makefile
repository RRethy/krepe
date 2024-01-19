BINARY_NAME=krepe
BINARY_PATH=./krepe/$(BINARY_NAME)
GOPATH=$(shell go env GOPATH)
GOCACHE=$(GOPATH)/pkg/mod
GOBUILD=go build
GOCLEAN=go clean
GOTEST=go test
GOMOD=go mod

.PHONY: build
build:
	$(GOBUILD) -C krepe -o $(BINARY_NAME) -v
	$(GOBUILD) -C jsonpatch

.PHONY: build-krepe
build-krepe:
	$(GOBUILD) -C krepe -o $(BINARY_NAME) -v

.PHONY: build-jsonpatch
build-jsonpatch:
	$(GOBUILD) -C jsonpatch

.PHONY: test
test:
	$(GOTEST) -C krepe ./...
	$(GOTEST) -C jsonpatch ./...

.PHONY: test-%
test-%:
	$(GOTEST) -C $* ./...

.PHONY: clean
clean:
	$(GOCLEAN)
	rm -f $(BINARY_PATH)

.PHONY: tidy
tidy:
	$(GOMOD) -C krepe tidy
	$(GOMOD) -C jsonpatch tidy

.PHONY: tidy-%
tidy-%:
	$(GOMOD) -C $* tidy
