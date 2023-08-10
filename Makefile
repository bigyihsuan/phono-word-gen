get_wasm_exec:
	cp "$$(go env GOROOT)/misc/wasm/wasm_exec.js" .

build: get_wasm_exec
	GOOS=js GOARCH=wasm go build -o main.wasm

run: build
	goexec 'http.ListenAndServe(`:8080`, http.FileServer(http.Dir(`.`)))'