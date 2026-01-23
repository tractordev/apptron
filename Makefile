VSCODE_URL	?= https://github.com/progrium/vscode-web/releases/download/v1/vscode-web-1.103.2.zip
DOCKER_CMD 	?= $(shell command -v podman || command -v docker)

VERSION		?= $(shell cat version.txt)
COMMIT		?= $(shell git rev-parse --short HEAD)
DATE		?= $(shell date +"%Y%m%d%H%M%S")

all: assets/vscode worker/node_modules extension/system/dist assets/wanix.min.js assets/wanix.wasm
.PHONY: all

dev: all .env.local
	docker rm -f $(shell docker ps -a --filter "name=^workerd-apptron-Session" --format "{{.ID}}") > /dev/null 2>&1 || true
	wrangler dev --port=8788
.PHONY: dev

deploy: all
	wrangler deploy
.PHONY: deploy

wasm: boot.go
	GOOS=js GOARCH=wasm go build -ldflags="-X 'main.Version=$(VERSION)-$(DATE)-$(COMMIT)'" -o ./assets/wanix.wasm .
.PHONY: wasm

ext:
	cd extension/system && npm run compile-web
.PHONY: ext

live-install: ./assets/wanix.wasm
	@if [ -z "$$ENV_UUID" ]; then \
		echo "ERROR: This is expected to run in an Apptron environment"; \
		exit 1; \
	fi
	@if [ -d /web/caches/assets/localhost:8788 ]; then \
		cp ./assets/wanix.wasm /web/caches/assets/localhost:8788/wanix.wasm; \
	else \
		cp ./assets/wanix.wasm /web/caches/assets/$(ENV_UUID).apptron.dev/wanix.wasm; \
	fi
.PHONY: live-install

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
	make ext

extension/system/node_modules:
	cd extension/system && npm ci

worker/node_modules: worker/package.json
	cd worker && npm ci

assets/wanix.wasm:
	make wasm

assets/wanix.min.js:
	$(DOCKER_CMD) rm -f apptron-wanix
	$(DOCKER_CMD) pull --platform linux/amd64 ghcr.io/tractordev/wanix:runtime
	$(DOCKER_CMD) create --name apptron-wanix --platform linux/amd64 ghcr.io/tractordev/wanix:runtime
	$(DOCKER_CMD) cp apptron-wanix:/wanix.min.js assets/wanix.min.js
	$(DOCKER_CMD) cp apptron-wanix:/wanix.js assets/wanix.js