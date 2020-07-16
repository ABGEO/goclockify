VERSION="1.0.0" # TODO: Get from app.

.PHONY: default
default: build-all

.PHONY: build-linux
build-linux-binary:
	@echo "+ $@"
	@GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -mod=vendor -a -installsuffix cgo -o ./build/binary/goclockify-linux-amd64 .

.PHONY: build-mac
build-mac-binary:
	@echo "+ $@"
	@GOOS=darwin GOARCH=amd64 CGO_ENABLED=0 go build -mod=vendor -a -installsuffix cgo -o ./build/binary/goclockify-darwin-amd64 .

.PHONY: build-rpm
build-rpm: build-linux-binary
	@echo "+ $@"
	@mkdir "build/rpm" -p
	@nfpm pkg --config nfpm.yaml --target ./build/rpm/goclockify-linux-amd64.rpm

.PHONY: build-deb
build-deb: build-linux-binary
	@echo "+ $@"
	@mkdir "build/deb" -p
	@nfpm pkg --config nfpm.yaml --target ./build/deb/goclockify-linux-amd64.deb

.PHONY: clean
clean:
	@echo "+ $@"
	@rm -rf "build/"

.PHONY: build-all
build-all: build-rpm build-deb build-mac-binary
