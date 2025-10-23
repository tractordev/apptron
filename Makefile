VSCODE_URL	?= https://github.com/progrium/vscode-web/releases/download/v1/vscode-web-1.103.2.zip
DOCKER_CMD 	?= $(shell command -v podman || command -v docker)

all: assets/vscode worker/node_modules extension/system/dist assets/wanix.min.js assets/wanix.wasm
.PHONY: all

dev: all .env.local
	docker rm -f $(shell docker ps -a --filter "name=^workerd-apptron-Session" --format "{{.ID}}") > /dev/null 2>&1 || true
	wrangler dev --port=8788
.PHONY: dev

deploy: all
	wrangler deploy
.PHONY: deploy

clean:
	rm -rf assets/vscode
	rm -rf worker/node_modules
	rm -rf node_modules
	rm -rf extension/system/dist
	rm -rf extension/system/node_modules
	rm -f assets/wanix.wasm
	rm -f assets/wanix.debug.wasm
	rm -f assets/wanix.js
	rm -f assets/wanix.min.js
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

extension/system/dist: extension/system/node_modules
	cd extension/system && npm run compile-web

extension/system/node_modules:
	cd extension/system && npm ci

worker/node_modules: worker/package.json
	cd worker && npm ci

assets/wanix.wasm:
	cd system && GOOS=js GOARCH=wasm go build -o ../assets/wanix.wasm

assets/wanix.min.js:
	$(DOCKER_CMD) rm -f apptron-wanix
	$(DOCKER_CMD) pull --platform linux/amd64 ghcr.io/tractordev/wanix:runtime
	$(DOCKER_CMD) create --name apptron-wanix --platform linux/amd64 ghcr.io/tractordev/wanix:runtime
	$(DOCKER_CMD) cp apptron-wanix:/wanix.min.js assets/wanix.min.js
	$(DOCKER_CMD) cp apptron-wanix:/wanix.js assets/wanix.js