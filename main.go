package main

import (
	"phono-word-gen/parts"

	"github.com/mroth/weightedrand/v2"
	"honnef.co/go/js/dom/v2"
)

func main() {
	document := dom.GetWindow().Document()
	outputTextElement := document.GetElementByID("outputText").(*dom.HTMLTextAreaElement)

	submitButton := document.GetElementByID("submit").(*dom.HTMLButtonElement)
	submitButton.AddEventListener("click", false, func(event dom.Event) {
		text := ""
		elements := []weightedrand.Choice[parts.CategoryElement, int]{
			weightedrand.NewChoice(parts.NewPhoneme("p"), 1),
			weightedrand.NewChoice(parts.NewPhoneme("t"), 1),
			weightedrand.NewChoice(parts.NewPhoneme("k"), 1),
		}
		category := parts.NewCategory("C", elements)
		for i := 0; i < 10; i++ {
			text += category.Get(make(map[string]parts.Category))
		}
		outputTextElement.SetValue(text)
	})

	// keep the go program alive
	select {}
}
