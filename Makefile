export GOPROXY = https://proxy.golang.org

GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean

CMD_DIR=./cmd
BIN_DIR=./bin

WEB3SAFE_BINARY_NAME=web3safe
WEB3SAFE_CMD_DIR=$(CMD_DIR)/$(WEB3SAFE_BINARY_NAME)

build: clean
	$(GOBUILD) -o $(BIN_DIR)/$(WEB3SAFE_BINARY_NAME) $(WEB3SAFE_CMD_DIR)/...

clean:
	$(GOCLEAN)
	rm -f $(WEB3SAFE_BINARY_NAME)

deps:
	$(GOCMD) mod tidy

test:
	$(GOCMD) test ./...

.PHONY: build clean deps test
