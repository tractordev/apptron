ROOT_DIR := $(dir $(realpath $(lastword $(MAKEFILE_LIST))))
GO_FILES := $(shell find . -name "*.go")
ifeq ($(OS),Windows_NT)
	EXE := ./dist/apptron.exe
else
	EXE := ./dist/apptron
endif
version=$(shell cat version)

apptron: $(GO_FILES)
	CGO_CFLAGS="-w" go build -ldflags $(ldflags) -o $(EXE) ./cmd/apptron/main.go

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

version=$(shell cat version)
branch=$(shell git branch --show-current)
ldflags="-X main.Version=$(version:dev=$(branch))"