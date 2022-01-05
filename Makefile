ROOT_DIR := $(dir $(realpath $(lastword $(MAKEFILE_LIST))))

.PHONY: ffi-static
ffi-static:
	cd lib/hostbridge && cargo build --release
	cp lib/hostbridge/target/release/libhostbridge.a lib/
	go build -o ./ffi-debug ./cmd/ffi-debug/main_static.go

.PHONY: ffi-shared
ffi-shared:
	cd lib/hostbridge && cargo build --release
	cp lib/hostbridge/target/release/libhostbridge.dylib lib/
	go build -o ./ffi-debug -ldflags="-r $(ROOT_DIR)lib" ./cmd/ffi-debug/main_shared.go

.PHONY: clean
clean:
	rm -rf ffi-debug lib/libhostbridge.dylib lib/libhostbridge.a lib/hostbridge/target