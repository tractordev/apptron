ROOT_DIR := $(dir $(realpath $(lastword $(MAKEFILE_LIST))))
RS_FILES := $(shell find lib/hostbridge/src -name "*.rs") 
GO_FILES := $(shell find . -name "*.go") 

hostbridge: lib/libhostbridge.a clientjs/dist/client.js $(GO_FILES)
	CGO_LDFLAGS="./lib/libhostbridge.a -ldl -framework Carbon -framework Cocoa -framework CoreFoundation -framework CoreVideo -framework IOKit -framework WebKit" \
	go build -a -o ./hostbridge ./cmd/hostbridge/main.go

debug-ffi: lib/libhostbridge.a $(GO_FILES)
	CGO_LDFLAGS="./lib/libhostbridge.a -ldl -framework Carbon -framework Cocoa -framework CoreFoundation -framework CoreVideo -framework IOKit -framework WebKit" \
	go build -tags ffi -a -o ./debug-ffi ./cmd/debug

debug-rpc: lib/libhostbridge.a $(GO_FILES)
	CGO_LDFLAGS="./lib/libhostbridge.a -ldl -framework Carbon -framework Cocoa -framework CoreFoundation -framework CoreVideo -framework IOKit -framework WebKit" \
	go build -tags rpc -a -o ./debug-rpc ./cmd/debug

debug-cmd: hostbridge $(GO_FILES)
	go build -tags cmd -o ./debug-cmd ./cmd/debug

clientjs/dist/client.js: clientjs/lib/*.js clientjs/src/*.ts
	make -C clientjs build

lib/libhostbridge.a: $(RS_FILES) lib/hostbridge/Cargo.toml
	cd lib/hostbridge && cargo build --release
	cp lib/hostbridge/target/release/libhostbridge.a lib/

.PHONY: clean
clean:
	rm -rf ./debug-ffi ./debug-rpc ./debug-cmd ./hostbridge ./lib/libhostbridge.a ./lib/hostbridge/target
