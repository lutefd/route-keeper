.PHONY: build install test clean

VERSION ?= $(shell git describe --tags --abbrev=0 2>/dev/null || echo "v0.1.0")
COMMIT := $(shell git rev-parse --short HEAD 2>/dev/null || echo "none")
BUILD_DATE := $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
BUILT_BY := $(shell whoami)
GO_VERSION := $(shell go version | awk '{print $$3}')

LDFLAGS = -X 'main.Version=$(VERSION)' \
          -X 'main.Commit=$(COMMIT)' \
          -X 'main.BuildDate=$(BUILD_DATE)' \
          -X 'main.BuiltBy=$(BUILT_BY)' \
          -X 'main.GoVersion=$(GO_VERSION)' \
          -s -w

build:
	@echo "Building route-keeper $(VERSION) ($(COMMIT))..."
	go build -ldflags="$(LDFLAGS)" -o bin/route-keeper ./cmd/route-keeper

install:
	@echo "Installing route-keeper $(VERSION)..."
	go install -ldflags="$(LDFLAGS)" ./cmd/route-keeper

test:
	go test -v ./...

clean:
	rm -rf bin/ dist/

release: clean
	@echo "Building release binaries for route-keeper $(VERSION)..."
	@mkdir -p dist
	
	# Linux
	GOOS=linux GOARCH=amd64 go build -ldflags="$(LDFLAGS)" -o dist/route-keeper-linux-amd64 ./cmd/route-keeper
	GOOS=linux GOARCH=arm64 go build -ldflags="$(LDFLAGS)" -o dist/route-keeper-linux-arm64 ./cmd/route-keeper
	
	# macOS
	GOOS=darwin GOARCH=amd64 go build -ldflags="$(LDFLAGS)" -o dist/route-keeper-darwin-amd64 ./cmd/route-keeper
	GOOS=darwin GOARCH=arm64 go build -ldflags="$(LDFLAGS)" -o dist/route-keeper-darwin-arm64 ./cmd/route-keeper
	
	# Windows
	GOOS=windows GOARCH=amd64 go build -ldflags="$(LDFLAGS)" -o dist/route-keeper-windows-amd64.exe ./cmd/route-keeper

	@echo "\nRelease binaries created in the dist/ directory:"
	@ls -la dist/
