VERSION=$(shell go run cmd/goclockify/main.go --version)
ARCHIVE_PREFIX="goclockify-$(VERSION)"

.PHONY: default
default: build-all

.PHONY: build-linux-binary
build-linux-binary:
	@echo "+ $@"
	@GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -mod=vendor -a -installsuffix cgo -o ./build/output/$(ARCHIVE_PREFIX)-linux-amd64 ./cmd/goclockify

.PHONY: build-mac-binary
build-mac-binary:
	@echo "+ $@"
	@GOOS=darwin GOARCH=amd64 CGO_ENABLED=0 go build -mod=vendor -a -installsuffix cgo -o ./build/output/$(ARCHIVE_PREFIX)-darwin-amd64 ./cmd/goclockify

.PHONY: build-rpm
build-rpm: build-linux-binary
	@echo "+ $@"
	@docker run --rm \
    	-v "$(PWD):/tmp" \
    	-e "VERSION=$(VERSION)" \
    	goreleaser/nfpm pkg \
    		--config /tmp/build/package/nfpm.yaml \
    		--target /tmp/build/output/$(ARCHIVE_PREFIX)-linux-amd64.rpm

.PHONY: build-deb
build-deb: build-linux-binary
	@echo "+ $@"
	@docker run --rm \
    	-v "$(PWD):/tmp" \
    	-e "VERSION=$(VERSION)" \
    	goreleaser/nfpm pkg \
    		--config /tmp/build/package/nfpm.yaml \
    		--target /tmp/build/output/$(ARCHIVE_PREFIX)-linux-amd64.deb

.PHONY: clean
clean:
	@echo "+ $@"
	@sudo rm -rf "build/output/"

.PHONY: modules
modules:
	@echo "+ $@"
	@go mod tidy
	@go mod vendor

.PHONY: build-all
build-all: build-rpm build-deb build-mac-binary
