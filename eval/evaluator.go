package eval

import (
	"errors"
	"fmt"
	"io"
	"math/rand"
	"regexp"
	"slices"
	"strconv"
	"strings"
	"syscall/js"

	"phono-word-gen/ast"
	"phono-word-gen/lex"
	"phono-word-gen/par"
	"phono-word-gen/parts"
	"phono-word-gen/sample"
	"phono-word-gen/util"

	"github.com/mroth/weightedrand/v2"
	"golang.org/x/exp/maps"
	"honnef.co/go/js/dom/v2"
)

type Evaluator struct {
	document dom.Document
	Elements
	Options

	generatedCount, duplicateCount, rejectedCount, replacedCount int

	categories parts.Categories
	components parts.Components
	syllables  []*parts.Syllable

	wordRejections     *regexp.Regexp
	syllableRejections *regexp.Regexp
	generalRejections  *regexp.Regexp

	replacements []parts.Replacement

	letters      []string
	letterRegexp *regexp.Regexp

	*examplePageElements

	errors []error
}

func New() (*Evaluator, error) {
	e := &Evaluator{}
	e.loadDocument()
	e.setEventListeners()
	return e, nil
}

func (e *Evaluator) loadDocument() {
	e.document = dom.GetWindow().Document()
	e.inputTextElement = e.document.QuerySelector("#phonology").(*dom.HTMLTextAreaElement)
	e.outputTextElement = e.document.QuerySelector("#outputText").(*dom.HTMLTextAreaElement)
	e.submitButtonElement = e.document.QuerySelector("#submit").(*dom.HTMLButtonElement)
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
	e.copyButtonElement = e.document.QuerySelector("#copyButton").(*dom.HTMLButtonElement)

	e.generatedAlertElement = e.document.QuerySelector("#generatedAlert").(*dom.HTMLDivElement)
	e.duplicateAlertElement = e.document.QuerySelector("#duplicateAlert").(*dom.HTMLDivElement)
	e.rejectedAlertElement = e.document.QuerySelector("#rejectedAlert").(*dom.HTMLDivElement)
	e.replacedAlertElement = e.document.QuerySelector("#replacedAlert").(*dom.HTMLDivElement)

	e.examplePageElements = nil
	if sampleDropdownElement := e.document.QuerySelector("#samples"); sampleDropdownElement != nil {
		e.examplePageElements = &examplePageElements{
			sampleDropdownElement: sampleDropdownElement.(*dom.HTMLSelectElement),
		}
	}
}

func (e *Evaluator) setEventListeners() {
	e.submitButtonElement.AddEventListener("click", false, e.submitMain)
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
	// set up copy-output button
	e.copyButtonElement.AddEventListener("click", false, func(event dom.Event) {
		js.Global().Get("window").Get("navigator").Get("clipboard").Call("writeText", e.outputTextElement.Value())
	})
}

func (e *Evaluator) submitMain(event dom.Event) {
	defer func() {
		if len(e.errors) > 0 {
			e.displayErrors()
		}
		e.clearErrors()
	}()
	// get the values of the various options
	e.getOptions()

	// if this is the example page, load the selected code
	if e.examplePageElements != nil {
		selectedExampleValue := e.sampleDropdownElement.Value()
		if selectedExampleValue == "nothing" {
			util.LogError("no example selected", selectedExampleValue)
			e.addErrors(fmt.Errorf("no example selected"))
			return
		}
		selectedExample, ok := sample.ExampleToFilename[selectedExampleValue]
		if !ok {
			util.LogError("invalid example selection", selectedExampleValue)
			return
		}

		file, err := sample.Examples.Open(selectedExample)
		if err != nil {
			util.LogError("failed to open example", err)
			e.addErrors(fmt.Errorf("failed to open example: %w", err))
			return
		}
		defer file.Close()

		data, err := io.ReadAll(file)
		if err != nil {
			util.LogError("failed to read example", err)
			e.addErrors(fmt.Errorf("failed to read example: %w", err))
			return
		}

		e.inputTextElement.SetValue(string(data))
	}

	// refesh the code input
	directives, err := e.loadCode(e.inputTextElement.Value())
	if err != nil {
		e.addErrors(err)
		return
	}
	e.evalDirectives(directives)
	if ok, err := e.checkCategories(); !ok {
		e.addErrors(err)
		return
	}
	if ok, err := e.checkComponents(); !ok {
		e.addErrors(err)
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
	words := e.generateWords(e.wordCount * 2)
	// convert the words to lists of syllables
	words = e.syllabizeWords(words)
	// if on, remove duplicates
	words = e.removeDuplicates(words)
	if len(words) >= e.wordCount {
		words = words[:e.wordCount]
	}

	// if on, apply rejections
	// TODO: allow contexts in the middle of rejection elements
	words = e.rejectWords(words)

	// TODO: if on, apply replacements
	// words = e.replaceWords(words)

	// if on, force generate to wordCount
	// get number of possible syllables, and abort forced gen if possible < wanted
	count := e.choiceCount(e.categories, e.components)
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
		e.addErrors(fmt.Errorf("not enough choices to force word count: only %d/%d choices available", count, e.wordCount))
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

func (e *Evaluator) loadCode(src string) ([]ast.Directive, error) {
	l := lex.New([]rune(src))
	p := par.New(l)
	directives := p.Directives()
	if len(p.Errors()) > 0 {
		return directives, errors.Join(p.Errors()...)
	}
	return directives, nil
}

func (e *Evaluator) checkCategories() (ok bool, err error) {
	// for each name/cat pair...
	for catName, cat := range e.categories {
		// for each element in the cat's elements...
		for _, element := range cat.Elements {
			// if the current element is a reference...
			reference, ok := element.Item.(*parts.CategoryReference)
			if !ok {
				continue
			}
			// if this reference is defined...
			reffedCat, ok := e.categories[reference.Name]
			if !ok {
				return false, parts.UndefinedCategoryError(catName, reference.Name)
			}
			// does it contain the cat?
			if slices.ContainsFunc(reffedCat.Elements, func(c weightedrand.Choice[parts.Element, int]) bool {
				item, ok := c.Item.(*parts.CategoryReference)
				return ok && item.Name == catName
			}) {
				return false, parts.RecursiveCategoryError(catName, reference.Name)
			}
		}
	}
	return true, nil
}

func (e *Evaluator) checkComponents() (ok bool, err error) {
	// for each name/comp pair...
	for compName, comp := range e.components {
		// for each element in the comp's elements...
		for _, element := range comp.Elements {
			// if the current element is a reference...
			reference, ok := element.(*parts.ComponentReference)
			if !ok {
				continue
			}
			// if this reference is defined...
			reffedComp, ok := e.components[reference.Name]
			if !ok {
				return false, parts.UndefinedComponentError(compName, reference.Name)
			}
			// does it contain the comp?
			if slices.ContainsFunc(reffedComp.Elements, func(c parts.SyllableElement) bool {
				item, ok := c.(*parts.CategoryReference)
				return ok && item.Name == compName
			}) {
				return false, parts.RecursiveComponentError(compName, reference.Name)
			}
		}
	}
	return true, nil
}

func (e *Evaluator) syllabizeWords(words []Word) []Word {
	for i, word := range words {
		err := word.GenerateSyllables(e.categories, e.components)
		if err != nil {
			util.LogError(err.Error())
			e.addErrors(err)
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

func (e *Evaluator) displayErrors() {
	errs := errors.Join(e.errors...)
	util.LogError(errs)
	e.outputTextElement.SetValue(errs.Error())
}

func (e *Evaluator) addErrors(errs ...error) {
	e.errors = append(e.errors, errs...)
}

func (e *Evaluator) clearErrors() {
	e.errors = []error{}
}

func (e *Evaluator) getOptions() {
	e.minSylCount = int(e.minSylCountElement.ValueAsNumber())
	e.maxSylCount = int(e.maxSylCountElement.ValueAsNumber())
	e.wordCount = int(e.wordCountElement.ValueAsNumber())
	e.sentenceCount = int(e.sentenceCountElement.ValueAsNumber())

	// handle minSylCount being larger than maxSylCount
	if e.minSylCount > e.maxSylCount {
		e.maxSylCount = e.minSylCount
		e.maxSylCountElement.SetValue(strconv.Itoa(e.minSylCount))
	}

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

func (e *Evaluator) choiceCount(categories parts.Categories, components parts.Components) int {
	count := len(e.syllables)
	for _, s := range e.syllables {
		count *= s.ChoiceCount(categories, components)
	}
	return count
}
