ROOT_DIR := $(dir $(realpath $(lastword $(MAKEFILE_LIST))))
SRC_FILES := $(shell find lib/hostbridge/src -name "*.rs") 

ffi-demo: lib/libhostbridge.a
	CGO_LDFLAGS="./lib/libhostbridge.a -ldl -framework Carbon -framework Cocoa -framework CoreFoundation -framework CoreVideo -framework IOKit -framework WebKit" \
	go build -a -o ./ffi-demo ./cmd/ffi-demo/main.go

rpc-demo: lib/libhostbridge.a
	CGO_LDFLAGS="./lib/libhostbridge.a -ldl -framework Carbon -framework Cocoa -framework CoreFoundation -framework CoreVideo -framework IOKit -framework WebKit" \
	go build -a -o ./rpc-demo ./cmd/rpc-demo/main.go

hostbridge: lib/libhostbridge.a
	CGO_LDFLAGS="./lib/libhostbridge.a -ldl -framework Carbon -framework Cocoa -framework CoreFoundation -framework CoreVideo -framework IOKit -framework WebKit" \
	go build -a -o ./hostbridge ./cmd/hostbridge/main.go

hostbridge-demo: hostbridge
	go build -o ./hostbridge-demo ./cmd/hostbridge-demo/main.go

lib/libhostbridge.a: $(SRC_FILES) lib/hostbridge/Cargo.toml
	cd lib/hostbridge && cargo build --release
	cp lib/hostbridge/target/release/libhostbridge.a lib/

.PHONY: clean
clean:
	rm -rf ./ffi-demo ./rpc-demo ./hostbridge ./hostbridge-demo ./lib/libhostbridge.a ./lib/hostbridge/target
