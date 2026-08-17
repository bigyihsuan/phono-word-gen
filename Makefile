go_files=$(wildcard *.go)

# https://go.dev/wiki/WebAssembly#getting-started
get_wasm_exec:
	mkdir -p ./build
	cp "$$(go env GOROOT)/lib/wasm/wasm_exec.js" ./build

compile $(go_files):
	GOOS=js GOARCH=wasm go build -C ./backend -o ../build/main.wasm
	cd ./..

build: get_wasm_exec compile
	cp ./build/* ./dist

server:
	python3 -m http.server --directory build

run: build ./build/wasm_exec.js server
