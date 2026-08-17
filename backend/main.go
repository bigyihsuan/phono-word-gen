package main

import (
	"fmt"
	"phono-word-gen/eval"
	"syscall/js"
)

// func main() {
// 	defer func() {
// 		if r := recover(); r != nil {
// 			fmt.Println("Recovered from", r)
// 		}
// 		initialize()
// 	}()
// 	initialize()
// }

// func initialize() {
// 	_, err := web.New()
// 	if err != nil {
// 		fmt.Println(err)
// 	}
// 	// keep the go program alive
// 	select {}
// }

func main() {
	fmt.Println("before global set")
	js.Global().Set("generate", js.FuncOf(eval.Generate))
	fmt.Println("after global set")
	select {}
}
