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

	"github.com/mroth/weightedrand/v2"
	"golang.org/x/exp/maps"
	"golang.org/x/exp/slices"
	"honnef.co/go/js/dom/v2"
)

type Evaluator struct {
	document dom.Document
	elements

	minSylCount, maxSylCount int
	wordCount                int

	forbidDuplicates, forceWordLimit, sortOutput                 bool
	applyRejections, applyReplacements                           bool
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

func (evaluator *Evaluator) loadDocument() {
	evaluator.document = dom.GetWindow().Document()
	evaluator.inputTextElement = evaluator.document.QuerySelector("#phonology").(*dom.HTMLTextAreaElement)
	evaluator.outputTextElement = evaluator.document.QuerySelector("#outputText").(*dom.HTMLTextAreaElement)
	evaluator.submitButton = evaluator.document.QuerySelector("#submit").(*dom.HTMLButtonElement)
	evaluator.minSylCountElement = evaluator.document.QuerySelector("#minSylCount").(*dom.HTMLInputElement)
	evaluator.maxSylCountElement = evaluator.document.QuerySelector("#maxSylCount").(*dom.HTMLInputElement)
	evaluator.wordCountElement = evaluator.document.QuerySelector("#wordCount").(*dom.HTMLInputElement)
	evaluator.forbidDuplicatesElement = evaluator.document.QuerySelector("#forbidDuplicates").(*dom.HTMLInputElement)
	evaluator.forceWordLimitElement = evaluator.document.QuerySelector("#forceWordLimit").(*dom.HTMLInputElement)
	evaluator.sortOutputElement = evaluator.document.QuerySelector("#sortOutput").(*dom.HTMLInputElement)
	evaluator.applyRejectionsElement = evaluator.document.QuerySelector("#applyRejections").(*dom.HTMLInputElement)
	evaluator.applyReplacementsElement = evaluator.document.QuerySelector("#applyReplacements").(*dom.HTMLInputElement)

	evaluator.generatedAlertElement = evaluator.document.QuerySelector("#generatedAlert").(*dom.HTMLDivElement)
	evaluator.duplicateAlertElement = evaluator.document.QuerySelector("#duplicateAlert").(*dom.HTMLDivElement)
	evaluator.rejectedAlertElement = evaluator.document.QuerySelector("#rejectedAlert").(*dom.HTMLDivElement)
	evaluator.replacedAlertElement = evaluator.document.QuerySelector("#replacedAlert").(*dom.HTMLDivElement)
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
	e.elements.submitButton.AddEventListener("click", false, e.submitMain)
}

func (e *Evaluator) submitMain(event dom.Event) {
	// get the values of the various options
	e.getOptions()

	// refesh the code input
	directives, err := e.loadCode(e.elements.inputTextElement.Value())
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
	words = e.replaceWords(words)

	// if on, force generate to wordCount
	// get number of possible syllables, and abort forced gen if possible < wanted
	count := e.ChoiceCount(e.categories)
	if e.forceWordLimit && count >= e.wordCount {
		for len(words) < e.wordCount {
			words = append(words, e.generateWords(e.wordCount)...)
			words = e.syllabizeWords(words)
			words = e.removeDuplicates(words)
			words = e.rejectWords(words)
			// TODO: apply replacements
		}
		rand.Shuffle(len(words), func(i, j int) {
			words[i], words[j] = words[j], words[i]
		})
		words = words[:e.wordCount]
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

	// display to the output textbox
	e.display(words, syllableSep)
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
			e.elements.outputTextElement.SetValue(err.Error())
			return words
		}
		words[i] = word
	}
	return words
}

func (e *Evaluator) display(words []Word, syllableSep string) {
	wordStrings := []string{}
	text := ""
	for _, word := range words {
		wordStrings = append(wordStrings, strings.Join(word.Syllables, syllableSep))
	}
	text += strings.Join(wordStrings, "\n")
	e.elements.outputTextElement.SetValue(text)

	e.updateAlerts()
}

func (e *Evaluator) displayError(err error) {
	util.LogError(err.Error())
	e.elements.outputTextElement.SetValue(err.Error())
}

func (e *Evaluator) getOptions() {
	e.minSylCount = int(e.minSylCountElement.ValueAsNumber())
	e.maxSylCount = int(e.maxSylCountElement.ValueAsNumber())
	e.wordCount = int(e.wordCountElement.ValueAsNumber())

	e.forbidDuplicates = e.forbidDuplicatesElement.Checked()
	e.forceWordLimit = e.forceWordLimitElement.Checked()
	e.sortOutput = e.sortOutputElement.Checked()
	e.applyRejections = e.applyRejectionsElement.Checked()
	e.applyReplacements = e.applyReplacementsElement.Checked()

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

func (e *Evaluator) ChoiceCount(categories map[string]parts.Category) int {
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

		matchesGeneral := len(e.wordRejections.String()) > 0 && e.generalRejections.MatchString(w)

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
