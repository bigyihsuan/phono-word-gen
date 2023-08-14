package main

import (
	"phono-word-gen/eval"

	"github.com/mroth/weightedrand/v2"
	"honnef.co/go/js/dom/v2"
)

func main() {
	document := dom.GetWindow().Document()
	outputTextElement := document.GetElementByID("outputText").(*dom.HTMLTextAreaElement)

	submitButton := document.GetElementByID("submit").(*dom.HTMLButtonElement)
	submitButton.AddEventListener("click", false, func(event dom.Event) {
		text := ""
		elements := []weightedrand.Choice[eval.CategoryElement, int]{
			weightedrand.NewChoice(eval.NewPhoneme("p"), 1),
			weightedrand.NewChoice(eval.NewPhoneme("t"), 1),
			weightedrand.NewChoice(eval.NewPhoneme("k"), 1),
		}
		category := eval.NewCategory("C", elements)
		for i := 0; i < 10; i++ {
			text += category.Get(make(map[string]eval.Category))
		}
		outputTextElement.SetValue(text)
	})

	// keep the go program alive
	select {}
}
