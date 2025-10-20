VSCODE_URL	?= https://github.com/progrium/vscode-web/releases/download/v1/vscode-web-1.103.2.zip
DOCKER_CMD 	?= $(shell command -v podman || command -v docker)

all: assets/vscode router/node_modules extension/dist session/bundle.tgz assets/wanix.min.js assets/wanix.wasm
.PHONY: all

dev: all .env.local
	wrangler dev --port=8788
.PHONY: dev

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

.env.local:
	cp .env.example .env.local

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

router/node_modules: router/package.json
	cd router && npm ci

session/bundle.tgz:
	cd bundle && make

assets/wanix.wasm:
	cd wanix && GOOS=js GOARCH=wasm go build -o ../assets/wanix.wasm

assets/wanix.min.js:
	$(DOCKER_CMD) rm -f apptron-wanix
	$(DOCKER_CMD) pull --platform linux/amd64 ghcr.io/tractordev/wanix:runtime
	$(DOCKER_CMD) create --name apptron-wanix --platform linux/amd64 ghcr.io/tractordev/wanix:runtime
	$(DOCKER_CMD) cp apptron-wanix:/wanix.min.js assets/wanix.min.js
	$(DOCKER_CMD) cp apptron-wanix:/wanix.js assets/wanix.js
#$(DOCKER_CMD) cp apptron-wanix:/wanix.debug.wasm assets/wanix.debug.wasm
#$(DOCKER_CMD) cp apptron-wanix:/wanix.wasm assets/wanix.wasm