BINARY     := helm-cm-delete
CMD        := ./cmd/helm-cm-delete
VERSION    := $(shell grep "^version:" plugin.yaml | awk '{print $$2}')
REVISION   := $(shell git rev-parse --short HEAD 2>/dev/null || echo "dev")
LDFLAGS    := -X main.version=$(VERSION) -X main.revision=$(REVISION)

export GO111MODULE := on
export CGO_ENABLED := 0

.PHONY: all build build-linux-amd64 build-linux-arm64 build-darwin-amd64 build-darwin-arm64 build-windows-amd64 \
        clean test lint fmt install remove

all: build

build:
	go build -ldflags "$(LDFLAGS)" -o bin/$(BINARY) $(CMD)

build-linux-amd64:
	GOOS=linux GOARCH=amd64 go build -ldflags "$(LDFLAGS)" -o bin/$(BINARY)_linux_amd64 $(CMD)

build-linux-arm64:
	GOOS=linux GOARCH=arm64 go build -ldflags "$(LDFLAGS)" -o bin/$(BINARY)_linux_arm64 $(CMD)

build-darwin-amd64:
	GOOS=darwin GOARCH=amd64 go build -ldflags "$(LDFLAGS)" -o bin/$(BINARY)_darwin_amd64 $(CMD)

build-darwin-arm64:
	GOOS=darwin GOARCH=arm64 go build -ldflags "$(LDFLAGS)" -o bin/$(BINARY)_darwin_arm64 $(CMD)

build-windows-amd64:
	GOOS=windows GOARCH=amd64 go build -ldflags "$(LDFLAGS)" -o bin/$(BINARY)_windows_amd64.exe $(CMD)

build-all: build-linux-amd64 build-linux-arm64 build-darwin-amd64 build-darwin-arm64 build-windows-amd64

clean:
	rm -rf bin/

test:
	go test -v ./...

lint:
	golangci-lint run ./...

fmt:
	go fmt ./...

install: build
	HELM_CM_DELETE_PLUGIN_NO_INSTALL_HOOK=1 helm plugin install .

reinstall: remove install

remove:
	helm plugin remove cm-delete
