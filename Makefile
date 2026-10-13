.PHONY: build

BIN := bin/oracle-mcp
ifeq ($(OS),Windows_NT)
BIN := bin/oracle-mcp.exe
endif

build:
	go build -trimpath -o $(BIN) ./cmd/server
