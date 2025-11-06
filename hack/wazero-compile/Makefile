
run-demo: gosrc/bin/go gosrc/pkg/tool/wasip1_wasm/compile
	cd demo && go run .
.PHONY: run-demo

gosrc/bin/go: gosrc
	cd gosrc/src && ./make.bash

gosrc/pkg/tool/wasip1_wasm/compile: gosrc
	cd gosrc/src && GOOS=wasip1 GOARCH=wasm ./make.bash

gosrc:
	git clone --depth 1 --branch go1.25.0 https://go.googlesource.com/go $@