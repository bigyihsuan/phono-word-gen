package sample

import "embed"

var ExampleToFilename = map[string]string{
	"example":  "example.txt",
	"chinese":  "chinese-ish.txt",
	"japanese": "japanese-ish.txt",
	"spanish":  "spanish-ish.txt",
}

//go:embed *.txt
var Examples embed.FS
