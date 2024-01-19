BINARY_NAME=krepe
BINARY_PATH=./krepe/$(BINARY_NAME)
INSTALL_PATH=~/bin/$(BINARY_NAME)
GOBUILD=go build
GOCLEAN=go clean
GOTEST=go test
GOMOD=go mod

.PHONY: build
build: build-krepe build-jsonpatch build-threewaymerge build-deepishcopy

.PHONY: build-%
build-%:
	$(GOBUILD) -C $*

.PHONY: test
test: test-krepe test-jsonpatch test-threewaymerge test-deepishcopy

.PHONY: test-%
test-%:
	$(GOTEST) -C $* ./...

.PHONY: clean
clean:
	$(GOCLEAN)
	rm -f $(BINARY_PATH)

.PHONY: tidy
tidy: tidy-krepe tidy-jsonpatch tidy-threewaymerge tidy-deepishcopy

.PHONY: tidy-%
tidy-%:
	cd $* && $(GOMOD) tidy

.PHONY: install
install: build-krepe
	cp $(BINARY_PATH) $(INSTALL_PATH)
