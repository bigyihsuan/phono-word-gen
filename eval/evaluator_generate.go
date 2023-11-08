package eval

import (
	"math/rand"
	"phono-word-gen/parts"
	"phono-word-gen/util"
	"strings"
	"unicode"
)

func (e *Evaluator) createSentences() {
	sentences := []string{}
	for i := 0; i < e.sentenceCount; i++ {
		sentence := e.generateSentence()
		if len(sentence) > 0 {
			sentences = append(sentences, sentence)
		}
	}
	e.displaySentences(sentences)
}

// generate a single sentence
func (e *Evaluator) generateSentence() string {
	wordCount := 1 + util.PeakedPowerLaw(15, 5, 50)
	sentenceWords := []string{}
	words := e.generateWords(wordCount * 2)
	words = e.syllabizeWords(words)
	words = e.rejectWords(words)
	for len(words) < wordCount {
		words = e.generateWords(wordCount * 2)
		words = e.syllabizeWords(words)
		words = e.rejectWords(words)
		rand.Shuffle(len(words), func(i, j int) {
			words[i], words[j] = words[j], words[i]
		})
		if len(words) >= wordCount {
			words = words[:wordCount]
		}
	}
	if len(words) >= wordCount {
		words = words[:wordCount]
	}
	// TODO: replacements
	// words = e.replaceWords(words)
	for i, w := range words {
		word, _ := w.Join()
		if i == 0 {
			runes := []rune(word)
			word = string(unicode.ToTitle(runes[0])) + string(runes[1:])
		}
		sentenceWords = append(sentenceWords, word)
	}
	sentence := strings.Join(sentenceWords, " ") + "."
	return sentence
}

// generate a `wordCount` number of words.
func (e *Evaluator) generateWords(wordCount int) (words []Word) {
	for i := 0; i < wordCount; i++ {
		syllableCount := min(e.minSylCount+util.PowerLaw(e.maxSylCount, 50), e.maxSylCount)
		words = append(words, e.generateWord(syllableCount))
	}
	e.generatedCount += e.wordCount
	return
}

func (e *Evaluator) generateWord(syllableCount int) Word {
	syllables := []*parts.Syllable{}
	for i := 0; i < syllableCount; i++ {
		syllable := e.syllables[rand.Intn(min(len(e.syllables)))]
		syllables = append(syllables, syllable)
	}
	return NewWord(syllables...)
}
