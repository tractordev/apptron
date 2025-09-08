VSCODE_URL	?= https://github.com/progrium/vscode-web/releases/download/v1/vscode-web-1.103.2.zip
DOCKER_CMD 	?= $(shell command -v podman || command -v docker)

all: assets/vscode router/node_modules extension/dist session/bundle.tgz assets/wanix.wasm
.PHONY: all

deploy: all
	wrangler deploy
.PHONY: deploy

clean:
	rm -rf assets/vscode
	rm -rf router/node_modules
	rm -rf node_modules
	rm -rf extension/dist
	rm -rf extension/node_modules
	rm -f assets/wanix.wasm
	rm -f assets/wanix.debug.wasm
	rm -f assets/wanix.js
	rm -f assets/wanix.min.js
	make -C bundle clean
.PHONY: clean

assets/vscode:
	curl -sL $(VSCODE_URL) -o assets/vscode.zip
	mkdir -p .tmp
	unzip assets/vscode.zip -d .tmp
	mv .tmp/dist/vscode assets/vscode
	rm -rf .tmp
	rm assets/vscode.zip

extension/dist: extension/node_modules
	cd extension && npm run compile-web

extension/node_modules:
	cd extension && npm ci

router/node_modules:
	cd router && npm ci

session/bundle.tgz:
	cd bundle && make

assets/wanix.wasm:
	$(DOCKER_CMD) rm -f apptron-wanix
	$(DOCKER_CMD) pull --platform linux/amd64 ghcr.io/tractordev/wanix:runtime
	$(DOCKER_CMD) create --name apptron-wanix --platform linux/amd64 ghcr.io/tractordev/wanix:runtime
	$(DOCKER_CMD) cp apptron-wanix:/ assets/