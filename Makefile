ROOT_DIR := $(dir $(realpath $(lastword $(MAKEFILE_LIST))))
GO_FILES := $(shell find . -name "*.go") 

#CGO_LDFLAGS="./lib/libhostbridge.a -ldl -framework Carbon -framework Cocoa -framework CoreFoundation -framework CoreVideo -framework IOKit -framework WebKit" \

hostbridge: clientjs/dist/client.js $(GO_FILES)
	go build -o ./hostbridge ./cmd/hostbridge/main.go

debug-pkg: $(GO_FILES)
	go build -tags pkg -o ./debug-pkg ./cmd/debug

debug-rpc: $(GO_FILES)
	go build -tags rpc -o ./debug-rpc ./cmd/debug

debug-cmd: hostbridge $(GO_FILES)
	go build -tags cmd -o ./debug-cmd ./cmd/debug

clientjs/dist/client.js: clientjs/lib/*.js clientjs/src/*.ts
	make -C clientjs build

.PHONY: clean
clean:
	rm -rf ./debug-pkg ./debug-rpc ./debug-cmd ./hostbridge
