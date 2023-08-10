package main

import (
	dom "honnef.co/go/js/dom/v2"
)

func main() {
	document := dom.GetWindow().Document()
	outputTextElement := document.GetElementByID("outputText").(*dom.HTMLTextAreaElement)
	outputTextElement.SetValue("hello world!")
}
