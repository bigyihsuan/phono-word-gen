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
		initialize()
	}()
	initialize()
}

func initialize() {
	_, err := eval.New()
	if err != nil {
		fmt.Println(err)
	}
	// keep the go program alive
	select {}
}
