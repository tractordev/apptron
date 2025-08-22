VSCODE_URL=https://github.com/progrium/vsclone/releases/download/v0.2/vscode-web.zip

all: assets/vscode router/node_modules extension/dist
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
.PHONY: clean

assets/vscode:
	curl -sL $(VSCODE_URL) -o assets/vscode.zip
	unzip assets/vscode.zip -d assets
	rm assets/vscode.zip

extension/dist: extension/node_modules
	cd extension && npm run compile-web

extension/node_modules:
	cd extension && npm ci

router/node_modules:
	cd router && npm ci

