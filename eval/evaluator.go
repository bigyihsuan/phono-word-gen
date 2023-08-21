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
	generatedCount, duplicateCount, rejectedCount, replacedCount int

	categories   map[string]parts.Category
	syllables    []*parts.Syllable
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
	wordSyllables := e.syllabizeWords(words)

	// if on, remove duplicates
	words, wordSyllables = e.removeDuplicates(words, wordSyllables)

	// if on, force generate to wordCount
	// get number of possible syllables, and abort forced gen if possible < wanted
	count := e.ChoiceCount(e.categories)
	util.Log(e.forceWordLimit, e.wordCount, len(wordSyllables), count)
	if e.forceWordLimit && count >= e.wordCount {
		for e.wordCount > len(wordSyllables) {
			needed := e.wordCount - len(wordSyllables)
			words = append(words, e.generateWords(needed)...)
			wordSyllables = append(wordSyllables, e.syllabizeWords(words)...)
			words, wordSyllables = e.removeDuplicates(words, wordSyllables)
		}
	} else if e.forceWordLimit && count < e.wordCount {
		e.displayError(fmt.Errorf("not enough choices to force word count: only %d/%d choices available", count, e.wordCount))
		return
	}

	// if on, sort output
	if e.sortOutput {
		// letter-based sorting
		if len(e.letters) > 0 {
			sort.Slice(wordSyllables, func(left, right int) bool {
				// letterize words
				l := strings.Join(wordSyllables[left], "")
				leftLetters := e.letterRegexp.FindAllString(l, -1)
				leftIndexes := []int{}
				for _, letter := range leftLetters {
					leftIndexes = append(leftIndexes, slices.Index(e.letters, letter))
				}
				r := strings.Join(wordSyllables[right], "")
				rightLetters := e.letterRegexp.FindAllString(r, -1)
				rightIndexes := []int{}
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
			sort.Slice(wordSyllables, func(i, j int) bool {
				a, b := wordSyllables[i], wordSyllables[j]
				as, bs := strings.Join(a, ""), strings.Join(b, "")
				return as < bs
			})
		}
	}

	syllableSep := ""
	// TODO: if on, display with syllable separators

	// display to the output textbox
	e.display(wordSyllables, syllableSep)
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

func (e *Evaluator) syllabizeWords(words []Word) (wordSyllables [][]string) {
	for _, word := range words {
		wordSyls, err := word.GenerateSyllables(e.categories)
		if err != nil {
			util.LogError(err.Error())
			e.elements.outputTextElement.SetValue(err.Error())
			return
		}
		wordSyllables = append(wordSyllables, wordSyls)
	}
	return
}

func (e *Evaluator) display(wordSyllables [][]string, syllableSep string) {
	wordStrings := []string{}
	text := ""
	for _, word := range wordSyllables {
		wordStrings = append(wordStrings, strings.Join(word, syllableSep))
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
	e.minSylCount = int(e.elements.minSylCountElement.ValueAsNumber())
	e.maxSylCount = int(e.elements.maxSylCountElement.ValueAsNumber())
	e.wordCount = int(e.elements.wordCountElement.ValueAsNumber())

	e.forbidDuplicates = e.elements.forbidDuplicatesElement.Checked()
	e.forceWordLimit = e.elements.forceWordLimitElement.Checked()
	e.sortOutput = e.elements.sortOutputElement.Checked()

	e.generatedCount = 0
	e.duplicateCount = 0
	e.rejectedCount = 0
	e.replacedCount = 0
}

func (e *Evaluator) updateAlerts() {
	e.generatedAlertElement.SetInnerHTML(fmt.Sprintf("generated %d words", e.generatedCount))
	e.duplicateAlertElement.SetInnerHTML(fmt.Sprintf("removed %d duplicates", e.duplicateCount))
	e.rejectedAlertElement.SetInnerHTML(fmt.Sprintf("rejected %d words", e.rejectedCount))
	e.replacedAlertElement.SetInnerHTML(fmt.Sprintf("rejected %d words", e.replacedCount))
}

func (e *Evaluator) removeDuplicates(words []Word, wordSyllables [][]string) (ws []Word, syls [][]string) {
	if e.forbidDuplicates {
		type entry struct {
			word Word
			syls []string
		}
		oldLen := len(wordSyllables)
		wordSet := make(map[string]entry)
		for i, word := range wordSyllables {
			w := strings.Join(word, "")
			if _, containsWord := wordSet[w]; !containsWord {
				wordSet[w] = entry{words[i], wordSyllables[i]}
			}
		}
		values := maps.Values(wordSet)
		ws = []Word{}
		syls = [][]string{}
		for _, v := range values {
			ws = append(ws, v.word)
			syls = append(syls, v.syls)
		}
		e.duplicateCount = oldLen - len(syls)
		return ws, syls
	} else {
		return words, wordSyllables
	}
}

func (e *Evaluator) ChoiceCount(categories map[string]parts.Category) int {
	count := len(e.syllables)
	for _, s := range e.syllables {
		count *= s.ChoiceCount(categories)
	}
	return count
}

type elements struct {
	inputTextElement        *dom.HTMLTextAreaElement
	outputTextElement       *dom.HTMLTextAreaElement
	submitButton            *dom.HTMLButtonElement
	minSylCountElement      *dom.HTMLInputElement
	maxSylCountElement      *dom.HTMLInputElement
	wordCountElement        *dom.HTMLInputElement
	forbidDuplicatesElement *dom.HTMLInputElement
	forceWordLimitElement   *dom.HTMLInputElement
	sortOutputElement       *dom.HTMLInputElement
	generatedAlertElement   *dom.HTMLDivElement
	duplicateAlertElement   *dom.HTMLDivElement
	rejectedAlertElement    *dom.HTMLDivElement
	replacedAlertElement    *dom.HTMLDivElement
}
