VERSION=$(shell go run main.go --version)
ARCHIVE_PREFIX="goclockify-$(VERSION)"

.PHONY: default
default: build-all

.PHONY: build-linux
build-linux-binary:
	@echo "+ $@"
	@GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -mod=vendor -a -installsuffix cgo -o ./build/$(ARCHIVE_PREFIX)-linux-amd64 .

.PHONY: build-mac-binary
build-mac-binary:
	@echo "+ $@"
	@GOOS=darwin GOARCH=amd64 CGO_ENABLED=0 go build -mod=vendor -a -installsuffix cgo -o ./build/$(ARCHIVE_PREFIX)-darwin-amd64 .

.PHONY: build-rpm
build-rpm: build-linux-binary
	@echo "+ $@"
	@docker run --rm \
    	-v "$(PWD):/tmp" \
    	-e "VERSION=$(VERSION)" \
    	goreleaser/nfpm pkg \
    		--config /tmp/nfpm.yaml \
    		--target /tmp/build/$(ARCHIVE_PREFIX)-linux-amd64.rpm

.PHONY: build-deb
build-deb: build-linux-binary
	@echo "+ $@"
	@docker run --rm \
    	-v "$(PWD):/tmp" \
    	-e "VERSION=$(VERSION)" \
    	goreleaser/nfpm pkg \
    		--config /tmp/nfpm.yaml \
    		--target /tmp/build/$(ARCHIVE_PREFIX)-linux-amd64.deb

.PHONY: clean
clean:
	@echo "+ $@"
	@rm -rf "build/"

.PHONY: modules
modules:
	@echo "+ $@"
	@go mod tidy
	@go mod vendor

.PHONY: build-all
build-all: build-rpm build-deb build-mac-binary
