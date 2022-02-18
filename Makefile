ROOT_DIR := $(dir $(realpath $(lastword $(MAKEFILE_LIST))))
SRC_FILES := $(shell find lib/hostbridge/src -name "*.rs") 

.PHONY: ffi-static
ffi-static: lib/libhostbridge.a
	CGO_LDFLAGS="./lib/libhostbridge.a -ldl -framework Carbon -framework Cocoa -framework CoreFoundation -framework CoreVideo -framework IOKit -framework WebKit" \
	go build -a -o ./ffi-debug ./cmd/ffi-debug/main_static.go

.PHONY: ffi-shared
ffi-shared: lib/libhostbridge.dylib
	go build -a -o ./ffi-debug -ldflags="-r $(ROOT_DIR)lib" ./cmd/ffi-debug/main_shared.go

lib/libhostbridge.dylib: $(SRC_FILES) lib/hostbridge/Cargo.toml
	cd lib/hostbridge && cargo build --release
	cp lib/hostbridge/target/release/libhostbridge.dylib lib/

lib/libhostbridge.a: $(SRC_FILES) lib/hostbridge/Cargo.toml
	cd lib/hostbridge && cargo build --release
	cp lib/hostbridge/target/release/libhostbridge.a lib/

.PHONY: clean
clean:
	rm -rf ffi-debug lib/libhostbridge.dylib lib/libhostbridge.a lib/hostbridge/target
