go_files=$(wildcard *.go)

get_wasm_exec:
	cp "$$(go env GOROOT)/misc/wasm/wasm_exec.js" .

build: $(go_files) main.wasm
	GOOS=js GOARCH=wasm go build -o main.wasm

server:
	python3 -m http.server

run: build ./wasm_exec.js server