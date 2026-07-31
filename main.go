package main

import (
	"fmt"
	"phono-word-gen/web"
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
	_, err := web.New()
	if err != nil {
		fmt.Println(err)
	}
	// keep the go program alive
	select {}
}
