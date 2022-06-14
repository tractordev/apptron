ROOT_DIR := $(dir $(realpath $(lastword $(MAKEFILE_LIST))))
GO_FILES := $(shell find . -name "*.go")
ifeq ($(OS),Windows_NT)
	EXE := ./apptron.exe
else
	EXE := ./apptron
endif

apptron: clientjs/dist/client.js $(GO_FILES)
	CGO_CFLAGS="-w" go build -o $(EXE) ./cmd/apptron/main.go

debug-pkg: $(GO_FILES)
	CGO_CFLAGS="-w" go build -tags pkg -o ./debug-pkg ./cmd/debug

debug-app: clientjs/dist/client.js cmd/debug/index.html $(GO_FILES) 
	go build -tags app -o ./debug-app ./cmd/debug

debug-cmd: apptron $(GO_FILES)
	go build -tags cmd -o ./debug-cmd ./cmd/debug

clientjs/dist/client.js: clientjs/lib/*.js clientjs/src/*.ts
	make -C clientjs build

.PHONY: install
install:
	CGO_CFLAGS="-w" go install ./cmd/apptron

.PHONY: clean
clean:
	rm -rf ./debug-pkg ./debug-app ./debug-cmd ./apptron
