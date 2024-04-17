## Make file for GO lang build(all platforms), run, test and release
.PHONY: build run test release

GO_FILES := $(shell find . -name '*.go')
CONFIG_FILE := config.yaml
VERSION := $(shell git describe --tags --abbrev=0)
BINARIES := bin/linux_amd64 bin/darwin_amd64 bin/windows_amd64.exe bin/linux_arm64 bin/darwin_arm64 bin/windows_arm64.exe
TARS := $(BINARIES:%=%.tar.gz)
OS := $(shell uname -s)
ARCH := $(shell uname -m)


build:	$(GO_FILES)
	@echo "Building for all platforms..."
	GOOS=linux GOARCH=amd64 go build -o bin/linux_amd64
	GOOS=darwin GOARCH=amd64 go build -o bin/darwin_amd64
	GOOS=windows GOARCH=amd64 go build -o bin/windows_amd64.exe
	GOOS=linux GOARCH=arm64 go build -o bin/linux_arm64
	GOOS=darwin GOARCH=arm64 go build -o bin/darwin_arm64
	GOOS=windows GOARCH=arm64 go build -o bin/windows_arm64.exe

run: build
	@echo "Running..."
	./bin/$(OS)_$(ARCH) --config $(CONFIG_FILE)
dev:
	@echo "Running..."
	go run . --config $(CONFIG_FILE)

test: $(GO_FILES)
	@echo "Testing..."
	go test -v ./...

$(BINARIES:%=%.tar.gz): %.tar.gz: %
	@echo "Creating tar..."
	tar -czvf $@ $< $(CONFIG_FILE)

release: test build $(TARS)
	@echo "Creating release..."
	gh release create $(VERSION) $(TARS) -t $(VERSION) -n "Release $(VERSION)"