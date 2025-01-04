VERSION=$(shell head -n 1 version)
GIT_REVISION=$(shell git log --pretty=format:'%h' -n 1)

TOP_DIR=$(dir $(realpath $(firstword $(MAKEFILE_LIST))))
PKG_SRC=$(shell find $(TOP_DIR)pkg -type f -name '*.go')
CMD_CLI_SRC=$(shell find $(TOP_DIR)cmd/cli -type f -name '*.go')

GO_BIN ?= go

all: converter-cli converter-cli-silicon converter-cli-linux converter-cli-linux-arm32 converter-cli-linux-arm64

converter-cli: $(PKG_SRC) $(CMD_CLI_SRC) go.mod go.sum
	$(GO_BIN) build $(LDFLAGS) -o $(TOP_DIR)bin/$@ $(TOP_DIR)cmd/cli

converter-cli-silicon: $(PKG_SRC) $(CMD_CLI_SRC) go.mod go.sum
	GOOS=darwin GOARCH=arm64 $(GO_BIN) build $(LDFLAGS) -o $(TOP_DIR)bin/converter-cli-silicon $(TOP_DIR)cmd/cli

converter-cli-linux: $(PKG_SRC) $(CMD_CLI_SRC) go.mod go.sum
	GOOS=linux GOARCH=amd64 $(GO_BIN) build $(LDFLAGS) -o $(TOP_DIR)bin/converter-cli-linux $(TOP_DIR)cmd/cli

converter-cli-linux-arm32: $(PKG_SRC) $(CMD_CLI_SRC) go.mod go.sum
	GOOS=linux GOARCH=arm $(GO_BIN) build $(LDFLAGS) -o $(TOP_DIR)bin/converter-cli-linux-arm32 $(TOP_DIR)cmd/cli

converter-cli-linux-arm64: $(PKG_SRC) $(CMD_CLI_SRC) go.mod go.sum
	GOOS=linux GOARCH=arm64 $(GO_BIN) build $(LDFLAGS) -o $(TOP_DIR)bin/converter-cli-linux-arm64 $(TOP_DIR)cmd/cli

dep:
	$(GO_BIN) mod download
	$(GO_BIN) mod tidy

clean:
	rm -rf bin/

lint: fmt vet

fmt:
	gofmt -w -s .

vet:
	$(GO_BIN) vet ./...
