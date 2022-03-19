
build:
	mkdir -p dist
	deno bundle -c tsconfig.json ./src/client.ts > ./dist/client.js

docs:
	npx typedoc src/client.ts
.PHONY: docs