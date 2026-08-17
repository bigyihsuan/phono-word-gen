go_files=$(wildcard *.go)
pages=$(wildcard ./pages/*)

# https://go.dev/wiki/WebAssembly#getting-started
get_wasm_exec:
	mkdir -p ./dist/
	cp "$$(go env GOROOT)/lib/wasm/wasm_exec.js" ./dist/
	cp "$$(go env GOROOT)/lib/wasm/wasm_exec.js" ./frontend/public/

compile $(go_files):
	GOOS=js GOARCH=wasm go build -C ./backend -o ../dist/main.wasm
	cd ./..

build: get_wasm_exec compile $(pages)
	cp ./dist/main.wasm ./frontend/public
	cp $(pages) ./dist

server:
	python3 -m http.server --directory dist

run: build ./dist/wasm_exec.js server