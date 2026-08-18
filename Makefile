go_files=$(wildcard *.go)

# https://go.dev/wiki/WebAssembly#getting-started
get_wasm_exec:
	cp "$$(go env GOROOT)/lib/wasm/wasm_exec.js" .

compile $(go_files):
	GOOS=js GOARCH=wasm go build -C ./backend -o ../main.wasm

build: get_wasm_exec compile
	cp ./main.wasm ./wasm_exec.js ./dist

server:
	python3 -m http.server --directory build

run: build ./build/wasm_exec.js server
