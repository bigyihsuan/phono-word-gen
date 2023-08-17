package main

import (
	"fmt"
	"phono-word-gen/eval"
)

func main() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered from", r)
		}
	}()
	evaluator, err := eval.New()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(evaluator)

	// keep the go program alive
	select {}
}
