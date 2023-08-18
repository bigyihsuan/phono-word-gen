package eval

import (
	"errors"
	"math/rand"
	"phono-word-gen/ast"
	"phono-word-gen/lex"
	"phono-word-gen/par"
	"phono-word-gen/parts"
	"phono-word-gen/util"
	"strings"

	"github.com/mroth/weightedrand/v2"
	"golang.org/x/exp/maps"
	"golang.org/x/exp/slices"
	"honnef.co/go/js/dom/v2"
)

type Evaluator struct {
	document dom.Document
	elements elements

	minSylCount, maxSylCount int
	wordCount                int

	forbidDuplicates, forceWordLimit, sortOutput bool

	categories map[string]parts.Category
	syllables  []*parts.Syllable
}

func New() (*Evaluator, error) {
	evaluator := &Evaluator{}
	evaluator.loadDocument()
	evaluator.setEventListeners()
	return evaluator, nil
}

func (evaluator *Evaluator) loadDocument() {
	evaluator.document = dom.GetWindow().Document()
	evaluator.elements.inputTextElement = evaluator.document.GetElementByID("phonology").(*dom.HTMLTextAreaElement)
	evaluator.elements.outputTextElement = evaluator.document.GetElementByID("outputText").(*dom.HTMLTextAreaElement)
	evaluator.elements.submitButton = evaluator.document.GetElementByID("submit").(*dom.HTMLButtonElement)
	evaluator.elements.minSylCountElement = evaluator.document.GetElementByID("minSylCount").(*dom.HTMLInputElement)
	evaluator.elements.maxSylCountElement = evaluator.document.GetElementByID("maxSylCount").(*dom.HTMLInputElement)
	evaluator.elements.wordCountElement = evaluator.document.GetElementByID("wordCount").(*dom.HTMLInputElement)
	evaluator.elements.forbidDuplicatesElement = evaluator.document.GetElementByID("forbidDuplicates").(*dom.HTMLInputElement)
	evaluator.elements.forceWordLimitElement = evaluator.document.GetElementByID("forceWordLimit").(*dom.HTMLInputElement)
	evaluator.elements.sortOutputElement = evaluator.document.GetElementByID("sortOutput").(*dom.HTMLInputElement)
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

	// TODO: if on, remove duplicates
	if e.forbidDuplicates {
		util.Log("forbidding duplicates")
		wordSet := make(map[string][]string)
		for i, word := range wordSyllables {
			w := strings.Join(word, "")
			if _, containsWord := wordSet[w]; !containsWord {
				wordSet[w] = wordSyllables[i]
			}
		}
		wordSyllables = maps.Values(wordSet)
	}

	// TODO: if on, sort
	if e.sortOutput {
		// TODO: letter-based sorting
		wordSyllableSortFunc := func(a, b []string) int {
			x, y := strings.Join(a, ""), strings.Join(b, "")
			if x < y {
				return -1
			} else if x == y {
				return 0
			} else {
				return 1
			}
		}
		slices.SortFunc(wordSyllables, wordSyllableSortFunc)
	}

	// TODO: if on, force generate to wordCount

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
}
