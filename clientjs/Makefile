
build:
	mkdir -p dist
	deno bundle -c tsconfig.json ./src/client.ts > ./dist/client.js
	npx typedoc src/client.ts