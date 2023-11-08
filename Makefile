go_files=$(wildcard *.go)

get_wasm_exec:
	cp "$$(go env GOROOT)/misc/wasm/wasm_exec.js" ./dist

build: get_wasm_exec $(go_files) ./dist/main.wasm ./dist/index.html ./dist/docs.html
	GOOS=js GOARCH=wasm go build -o ./dist/main.wasm

server:
	python3 -m http.server --directory dist

run: build ./dist/wasm_exec.js server