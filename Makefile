go_files=$(wildcard *.go)
pages=$(wildcard ./pages/*)

# https://go.dev/wiki/WebAssembly#getting-started
get_wasm_exec:
	cp "$$(go env GOROOT)/lib/wasm/wasm_exec.js" ./dist

build: get_wasm_exec $(go_files) ./dist/main.wasm $(pages)
	GOOS=js GOARCH=wasm go build -o ./dist/main.wasm
	cp $(pages) ./dist

server:
	python3 -m http.server --directory dist

run: build ./dist/wasm_exec.js server