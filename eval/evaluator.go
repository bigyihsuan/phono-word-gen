package eval

import (
	"errors"
	"fmt"
	"math/rand"
	"phono-word-gen/ast"
	"phono-word-gen/lex"
	"phono-word-gen/par"
	"phono-word-gen/parts"
	"phono-word-gen/util"
	"regexp"
	"sort"
	"strings"
	"unicode"

	"github.com/mroth/weightedrand/v2"
	"golang.org/x/exp/maps"
	"golang.org/x/exp/slices"
	"honnef.co/go/js/dom/v2"
)

type Evaluator struct {
	document dom.Document
	Elements
	Options

	generatedCount, duplicateCount, rejectedCount, replacedCount int

	categories map[string]parts.Category
	syllables  []*parts.Syllable

	wordRejections     *regexp.Regexp
	syllableRejections *regexp.Regexp
	generalRejections  *regexp.Regexp

	replacements []parts.Replacement

	letters      []string
	letterRegexp *regexp.Regexp
}

func New() (*Evaluator, error) {
	evaluator := &Evaluator{}
	evaluator.loadDocument()
	evaluator.setEventListeners()
	return evaluator, nil
}

func (e *Evaluator) loadDocument() {
	e.document = dom.GetWindow().Document()
	e.inputTextElement = e.document.QuerySelector("#phonology").(*dom.HTMLTextAreaElement)
	e.outputTextElement = e.document.QuerySelector("#outputText").(*dom.HTMLTextAreaElement)
	e.submitButton = e.document.QuerySelector("#submit").(*dom.HTMLButtonElement)
	e.minSylCountElement = e.document.QuerySelector("#minSylCount").(*dom.HTMLInputElement)
	e.maxSylCountElement = e.document.QuerySelector("#maxSylCount").(*dom.HTMLInputElement)
	e.wordCountElement = e.document.QuerySelector("#wordCount").(*dom.HTMLInputElement)
	e.sentenceCountElement = e.document.QuerySelector("#sentenceCount").(*dom.HTMLInputElement)
	e.generateSentencesElement = e.document.QuerySelector("#generateSentences").(*dom.HTMLInputElement)
	e.forbidDuplicatesElement = e.document.QuerySelector("#forbidDuplicates").(*dom.HTMLInputElement)
	e.forceWordLimitElement = e.document.QuerySelector("#forceWordLimit").(*dom.HTMLInputElement)
	e.sortOutputElement = e.document.QuerySelector("#sortOutput").(*dom.HTMLInputElement)
	e.markSyllablesElement = e.document.QuerySelector("#markSyllables").(*dom.HTMLInputElement)
	e.applyRejectionsElement = e.document.QuerySelector("#applyRejections").(*dom.HTMLInputElement)
	e.applyReplacementsElement = e.document.QuerySelector("#applyReplacements").(*dom.HTMLInputElement)

	e.generatedAlertElement = e.document.QuerySelector("#generatedAlert").(*dom.HTMLDivElement)
	e.duplicateAlertElement = e.document.QuerySelector("#duplicateAlert").(*dom.HTMLDivElement)
	e.rejectedAlertElement = e.document.QuerySelector("#rejectedAlert").(*dom.HTMLDivElement)
	e.replacedAlertElement = e.document.QuerySelector("#replacedAlert").(*dom.HTMLDivElement)
}

func (evaluator *Evaluator) loadCode(src string) ([]ast.Directive, error) {
	l := lex.New([]rune(src))
	p := par.New(l)
	directives := p.Directives()
	if len(p.Errors()) > 0 {
		return directives, errors.Join(p.Errors()...)
	}
	return directives, nil
}

func (evaluator *Evaluator) checkCategories() (ok bool, err error) {
	// for each name/cat pair...
	for catName, cat := range evaluator.categories {
		// for each element in the cat's elements...
		for _, element := range cat.Elements {
			// if the current element is a reference...
			reference, ok := element.Item.(*parts.Reference)
			if !ok {
				continue
			}
			// if this reference is defined...
			reffedCat, ok := evaluator.categories[reference.Name]
			if !ok {
				return false, parts.UndefinedCategoryError(catName, reference.Name)
			}
			// does it contain the cat?
			if slices.ContainsFunc(reffedCat.Elements, func(c weightedrand.Choice[parts.Element, int]) bool {
				item, ok := c.Item.(*parts.Reference)
				return ok && item.Name == catName
			}) {
				return false, parts.RecursiveCategoryError(catName, reference.Name)
			}
		}
	}
	return true, nil
}

func (e *Evaluator) setEventListeners() {
	e.submitButton.AddEventListener("click", false, e.submitMain)
	e.generateSentencesElement.AddEventListener("click", false, func(event dom.Event) {
		if e.generateSentencesElement.Checked() {
			e.forbidDuplicatesElement.SetDisabled(true)
			e.forceWordLimitElement.SetDisabled(true)
			e.markSyllablesElement.SetDisabled(true)
			e.sortOutputElement.SetDisabled(true)

			e.wordCountElement.SetDisabled(true)
			e.sentenceCountElement.SetDisabled(false)
		} else {
			e.forbidDuplicatesElement.SetDisabled(false)
			e.forceWordLimitElement.SetDisabled(false)
			e.markSyllablesElement.SetDisabled(false)
			e.sortOutputElement.SetDisabled(false)

			e.wordCountElement.SetDisabled(false)
			e.sentenceCountElement.SetDisabled(true)
		}
	})
}

func (e *Evaluator) submitMain(event dom.Event) {
	// get the values of the various options
	e.getOptions()

	// refesh the code input
	directives, err := e.loadCode(e.inputTextElement.Value())
	if err != nil {
		e.displayError(err)
		return
	}
	e.evalDirectives(directives)
	if ok, err := e.checkCategories(); !ok {
		e.displayError(err)
		return
	}

	// don't try to generate if we have no syllables
	if len(e.syllables) < 1 {
		return
	}

	if e.generateSentences {
		e.createSentences()
		return
	}
	// generate N words
	words := e.generateWords(e.wordCount)
	// convert the words to lists of syllables
	words = e.syllabizeWords(words)
	// if on, remove duplicates
	words = e.removeDuplicates(words)

	// if on, apply rejections
	// TODO: allow contexts in the middle of rejection elements
	words = e.rejectWords(words)

	// TODO: if on, apply replacements
	// words = e.replaceWords(words)

	// if on, force generate to wordCount
	// get number of possible syllables, and abort forced gen if possible < wanted
	count := e.choiceCount(e.categories)
	if e.forceWordLimit && count >= e.wordCount {
		for len(words) < e.wordCount {
			words = e.generateWords(e.wordCount * 2)
			words = e.syllabizeWords(words)
			words = e.removeDuplicates(words)
			words = e.rejectWords(words)
			// TODO: apply replacements
			// words = e.replaceWords(words)
		}
		rand.Shuffle(len(words), func(i, j int) {
			words[i], words[j] = words[j], words[i]
		})
		if len(words) >= e.wordCount {
			words = words[:e.wordCount]
		}
	} else if e.forceWordLimit && count < e.wordCount {
		e.displayError(fmt.Errorf("not enough choices to force word count: only %d/%d choices available", count, e.wordCount))
		return
	}

	// if on, sort output
	if e.sortOutput {
		words = e.sort(words)
	}

	syllableSep := ""
	// TODO: if on, display with syllable separators
	if e.markSyllables {
		syllableSep = "."
	}

	// display to the output textbox
	e.displayWords(words, syllableSep)
}

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
		util.Log("adding more words", len(words), wordCount)
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

func (e *Evaluator) syllabizeWords(words []Word) []Word {
	for i, word := range words {
		err := word.GenerateSyllables(e.categories)
		if err != nil {
			util.LogError(err.Error())
			e.outputTextElement.SetValue(err.Error())
			return words
		}
		words[i] = word
	}
	return words
}

func (e *Evaluator) displayWords(words []Word, syllableSep string) {
	wordStrings := []string{}
	text := ""
	for _, word := range words {
		wordStrings = append(wordStrings, strings.Join(word.Syllables, syllableSep))
	}
	text += strings.Join(wordStrings, "\n")
	e.outputTextElement.SetValue(text)
	e.updateAlerts()
}

func (e *Evaluator) displaySentences(sentences []string) {
	text := strings.Join(sentences, " ")
	e.outputTextElement.SetValue(text)
	e.updateAlerts()
}

func (e *Evaluator) displayError(err error) {
	util.LogError(err.Error())
	e.outputTextElement.SetValue(err.Error())
}

func (e *Evaluator) getOptions() {
	e.minSylCount = int(e.minSylCountElement.ValueAsNumber())
	e.maxSylCount = int(e.maxSylCountElement.ValueAsNumber())
	e.wordCount = int(e.wordCountElement.ValueAsNumber())
	e.sentenceCount = int(e.sentenceCountElement.ValueAsNumber())

	e.forbidDuplicates = e.forbidDuplicatesElement.Checked()
	e.forceWordLimit = e.forceWordLimitElement.Checked()
	e.sortOutput = e.sortOutputElement.Checked()
	e.markSyllables = e.markSyllablesElement.Checked()
	e.applyRejections = e.applyRejectionsElement.Checked()
	e.applyReplacements = e.applyReplacementsElement.Checked()
	e.generateSentences = e.generateSentencesElement.Checked()

	e.generatedCount = 0
	e.duplicateCount = 0
	e.rejectedCount = 0
	e.replacedCount = 0
}

func (e *Evaluator) updateAlerts() {
	e.generatedAlertElement.SetInnerHTML(fmt.Sprintf("generated %d words", e.generatedCount))
	e.duplicateAlertElement.SetInnerHTML(fmt.Sprintf("removed %d duplicates", e.duplicateCount))
	e.rejectedAlertElement.SetInnerHTML(fmt.Sprintf("rejected %d words", e.rejectedCount))
	e.replacedAlertElement.SetInnerHTML(fmt.Sprintf("replaced %d words", e.replacedCount))
}

func (e *Evaluator) removeDuplicates(words []Word) (ws []Word) {
	if !e.forbidDuplicates {
		return words
	}

	oldLen := len(words)
	wordSet := make(map[string]Word)
	for i, word := range words {
		joined, _ := word.Join()
		if _, containsWord := wordSet[joined]; !containsWord {
			wordSet[joined] = words[i]
		}
	}
	values := maps.Values(wordSet)
	ws = []Word{}
	for _, v := range values {
		ws = append(ws, v)
	}
	e.duplicateCount = oldLen - len(ws)
	return ws
}

func (e *Evaluator) choiceCount(categories map[string]parts.Category) int {
	count := len(e.syllables)
	for _, s := range e.syllables {
		count *= s.ChoiceCount(categories)
	}
	return count
}

func (e *Evaluator) rejectWords(words []Word) []Word {
	if !e.applyRejections {
		return words
	}

	keptWords := []Word{}

	for i, word := range words {
		w, _ := word.Join()

		matchesWordLevel := len(e.wordRejections.String()) > 0 && e.wordRejections.MatchString(w)

		matchesSyllableLevel := false
		if len(e.syllableRejections.String()) > 0 {
			for _, syl := range word.Syllables {
				if e.syllableRejections.MatchString(syl) {
					matchesSyllableLevel = true
					break
				}
			}
		}

		matchesGeneral := len(e.generalRejections.String()) > 0 && e.generalRejections.MatchString(w)

		if !matchesWordLevel && !matchesSyllableLevel && !matchesGeneral {
			keptWords = append(keptWords, words[i])
		} else {
			e.rejectedCount++
		}
	}
	return keptWords
}

func (e *Evaluator) replaceWords(words []Word) []Word {
	if !e.applyReplacements {
		return words
	}

	replacedWords := []Word{}

	// for each word...
	for _, word := range words {
		w, _ := word.Join()
		// w, syllableIndexes := word.Join()

		// for each possible replacement...
		for _, r := range e.replacements {
			matchesException := false

			// check for exception; if true, don't replace
			exception := r.ExceptionRegexp(e.categories)
			if exception != nil {
				if r.Exception.IsSyllableLevel() {
					matchesException = true
					for _, syl := range word.Syllables {
						if !exception.MatchString(syl) {
							matchesException = false
							break
						}
					}
				} else {
					matchesException = exception.MatchString(w)
				}
			}
			if matchesException {
				// return early
				replacedWords = append(replacedWords, word)
				continue
			}

			// need to find the match indexes to replace across syllable boundaries
			// match against the word, and get the indexes
			// map word-indexes into syllable-letter-indexes
			//     word[i] => syllables[j][k]
			// this is so that the replacement can span across syllables
			// though, there's too many parens in the generated regexp
			// so need to figure out how to get the start and end index
			// w/out using subgroups

			condition := r.ConditionRegexp(e.categories)
			matchIndexes := [][]int{}
			if r.Condition.IsSyllableLevel() {
				for _, syl := range word.Syllables {
					match := condition.FindAllStringIndex(syl, -1)
					if match == nil {
						continue
					}
					matchIndexes = append(matchIndexes, match...)
				}
			} else {
				matchIndexes = condition.FindAllStringIndex(w, -1)
			}
			m, _ := util.ToMap(matchIndexes)
			util.Log(word, m)
			util.Log(r.Source.Regexp(e.categories), r.Replacement)
			util.Log("ConditionRegexp", r.ConditionRegexp(e.categories))
			if matchIndexes == nil || len(matchIndexes) == 0 {
				// return early
				replacedWords = append(replacedWords, word)
				continue
			}
			// TODO: replacement
			// we have the indexes of the matches in either the joined word or syllables

		}
	}
	return replacedWords
}

func (e *Evaluator) sort(words []Word) []Word {
	// letter-based sorting
	if len(e.letters) > 0 {
		sort.Slice(words, func(left, right int) bool {
			// letterize words
			// join into a single string
			l := strings.Join(words[left].Syllables, "")
			r := strings.Join(words[right].Syllables, "")
			// find all (known) letters
			leftLetters := e.letterRegexp.FindAllString(l, -1)
			rightLetters := e.letterRegexp.FindAllString(r, -1)
			// for each letter found, find the index of that letter in the letter directive
			leftIndexes := []int{}
			rightIndexes := []int{}
			for _, letter := range leftLetters {
				leftIndexes = append(leftIndexes, slices.Index(e.letters, letter))
			}
			for _, letter := range rightLetters {
				rightIndexes = append(rightIndexes, slices.Index(e.letters, letter))
			}
			minLen := min(len(leftIndexes), len(rightIndexes))

			for i := 0; i < minLen; i++ {
				if leftIndexes[i] < rightIndexes[i] {
					return true
				}
				if leftIndexes[i] > rightIndexes[i] {
					return false
				}
			}
			if len(leftIndexes) < len(rightIndexes) {
				return true
			}
			if len(leftIndexes) > len(rightIndexes) {
				return false
			}
			return false
		})
	} else {
		sort.Slice(words, func(i, j int) bool {
			a, b := words[i], words[j]
			as, bs := strings.Join(a.Syllables, ""), strings.Join(b.Syllables, "")
			less := as < bs
			if less {
				words[i], words[j] = words[j], words[i]
			}
			return less
		})
	}
	return words
}
