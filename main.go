package main

import (
	"fmt"
	"math/rand"
	"phono-word-gen/parts"
	"strconv"
	"strings"

	"github.com/mroth/weightedrand/v2"
	"honnef.co/go/js/dom/v2"
)

var document = dom.GetWindow().Document()
var outputTextElement = document.GetElementByID("outputText").(*dom.HTMLTextAreaElement)
var submitButton = document.GetElementByID("submit").(*dom.HTMLButtonElement)
var minSylCountElement = document.GetElementByID("minSylCount").(*dom.HTMLInputElement)
var maxSylCountElement = document.GetElementByID("maxSylCount").(*dom.HTMLInputElement)
var wordCountElement = document.GetElementByID("wordCount").(*dom.HTMLInputElement)

func main() {

	submitButton.AddEventListener("click", false, func(event dom.Event) {
		text := ""
		c, _ := parts.NewCategory(
			weightedrand.NewChoice(parts.NewPhoneme("p"), 1),
			weightedrand.NewChoice(parts.NewPhoneme("t"), 1),
			weightedrand.NewChoice(parts.NewPhoneme("k"), 1),
		)
		v, _ := parts.NewCategory(
			weightedrand.NewChoice(parts.NewPhoneme("a"), 1),
			weightedrand.NewChoice(parts.NewPhoneme("i"), 1),
			weightedrand.NewChoice(parts.NewPhoneme("u"), 1),
		)
		categories := map[string]parts.Category{"C": c, "V": v}
		syllable := parts.NewSyllable(
			parts.NewReference("C").(parts.SyllableElement),
			parts.NewReference("V").(parts.SyllableElement),
		)

		minCount, _ := strconv.Atoi(minSylCountElement.Value())
		maxCount, _ := strconv.Atoi(maxSylCountElement.Value())

		wordCount, _ := strconv.Atoi(wordCountElement.Value())

		words := [][]string{}

		for i := 0; i < wordCount; i++ {
			syllables := []string{}

			syllableCount := min(minCount+powerLaw(maxCount, 50), maxCount)
			for i := 0; i < syllableCount; i++ {
				syllables = append(syllables, syllable.Get(categories))
			}
			syllables = append([]string{fmt.Sprintf("%d ", syllableCount)}, syllables...)
			words = append(words, syllables)
		}
		wordStrings := []string{}
		for _, word := range words {
			wordStrings = append(wordStrings, strings.Join(word, ""))
		}
		text += strings.Join(wordStrings, "\n")
		outputTextElement.SetValue(text)
	})

	// keep the go program alive
	select {}
}

func powerLaw(max, percentage int) int {
	for r := 0; ; r = (r + 1) % max {
		if randomPercentage() < percentage {
			return r
		}
	}
}

func randomPercentage() int {
	return rand.Intn(101) + 1
}
