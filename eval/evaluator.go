package eval

import (
	"errors"
	"math/rand"
	"phono-word-gen/ast"
	"phono-word-gen/lex"
	"phono-word-gen/par"
	"phono-word-gen/parts"
	"strconv"
	"strings"

	"github.com/mroth/weightedrand/v2"
	"golang.org/x/exp/slices"
	"honnef.co/go/js/dom/v2"
)

type Evaluator struct {
	document dom.Document

	inputTextElement   *dom.HTMLTextAreaElement
	outputTextElement  *dom.HTMLTextAreaElement
	submitButton       *dom.HTMLButtonElement
	minSylCountElement *dom.HTMLInputElement
	maxSylCountElement *dom.HTMLInputElement
	wordCountElement   *dom.HTMLInputElement

	minSylCount, maxSylCount int
	wordCount                int

	categories map[string]parts.Category
	syllables  []*parts.Syllable
}

func New() (*Evaluator, error) {
	evaluator := &Evaluator{}

	evaluator.loadDocument()
	directives, err := evaluator.loadCode(evaluator.inputTextElement.Value())
	if err != nil {
		return evaluator, err
	}
	evaluator.evalDirectives(directives)
	if ok, err := evaluator.checkCategories(); !ok {
		return evaluator, err
	}
	evaluator.setEventListeners()
	return evaluator, nil
}

func (evaluator *Evaluator) loadDocument() {
	evaluator.document = dom.GetWindow().Document()
	evaluator.inputTextElement = evaluator.document.GetElementByID("phonology").(*dom.HTMLTextAreaElement)
	evaluator.outputTextElement = evaluator.document.GetElementByID("outputText").(*dom.HTMLTextAreaElement)
	evaluator.submitButton = evaluator.document.GetElementByID("submit").(*dom.HTMLButtonElement)
	evaluator.minSylCountElement = evaluator.document.GetElementByID("minSylCount").(*dom.HTMLInputElement)
	evaluator.maxSylCountElement = evaluator.document.GetElementByID("maxSylCount").(*dom.HTMLInputElement)
	evaluator.wordCountElement = evaluator.document.GetElementByID("wordCount").(*dom.HTMLInputElement)
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

func (evaluator *Evaluator) setEventListeners() {
	evaluator.submitButton.AddEventListener("click", false, func(event dom.Event) {
		text := ""
		evaluator.minSylCount, _ = strconv.Atoi(evaluator.minSylCountElement.Value())
		evaluator.maxSylCount, _ = strconv.Atoi(evaluator.maxSylCountElement.Value())
		evaluator.wordCount, _ = strconv.Atoi(evaluator.wordCountElement.Value())

		words := [][]string{}

		directives, err := evaluator.loadCode(evaluator.inputTextElement.Value())
		if err != nil {
			logError(err.Error())
			evaluator.outputTextElement.SetValue(err.Error())
			return
		}
		evaluator.evalDirectives(directives)
		if ok, err := evaluator.checkCategories(); !ok {
			logError(err.Error())
			evaluator.outputTextElement.SetValue(err.Error())
			return
		}

		syllable := evaluator.syllables[rand.Intn(len(evaluator.syllables))]

		for i := 0; i < evaluator.wordCount; i++ {
			syllables := []string{}

			syllableCount := min(evaluator.minSylCount+powerLaw(evaluator.maxSylCount, 50), evaluator.maxSylCount)
			for i := 0; i < syllableCount; i++ {
				syl, err := syllable.Get(evaluator.categories)
				if err != nil {
					text += err.Error()
				}
				syllables = append(syllables, syl)
			}
			// syllables = append([]string{fmt.Sprintf("%d ", syllableCount)}, syllables...)
			words = append(words, syllables)
		}
		wordStrings := []string{}
		for _, word := range words {
			wordStrings = append(wordStrings, strings.Join(word, ""))
		}
		text += strings.Join(wordStrings, "\n")
		evaluator.outputTextElement.SetValue(text)
	})
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

func log(o ...any) {
	dom.GetWindow().Console().Call("log", o...)
}
func logError(o ...any) {
	dom.GetWindow().Console().Call("error", o...)
}
