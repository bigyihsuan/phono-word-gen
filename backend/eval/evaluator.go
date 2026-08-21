package eval

import (
	"errors"
	"fmt"
	"math/rand/v2"
	"regexp"
	"slices"
	"strings"
	"syscall/js"
	"time"

	"phono-word-gen/ast"
	"phono-word-gen/lex"
	"phono-word-gen/par"
	"phono-word-gen/parts"
	"phono-word-gen/util"

	"github.com/mroth/weightedrand/v3"
	"golang.org/x/exp/maps"
)

type Evaluator struct {
	Options

	GeneratedCount, DuplicateCount, RejectedCount, ReplacedCount int

	categories parts.Categories
	components parts.Components
	syllables  []*parts.Syllable

	wordRejections     *regexp.Regexp
	syllableRejections *regexp.Regexp
	generalRejections  *regexp.Regexp

	replacements []parts.Replacement

	letters      []string
	letterRegexp *regexp.Regexp

	Errors []error

	rs *rand.Rand // global random source
}

func Generate(this js.Value, inputs []js.Value) any {
	input := inputs[0]
	opts := OptionsFromJsValue(input)
	e := New(opts, nil)
	words, sep, sentences := e.Run()

	wordText := ""
	if len(words) > 0 {
		wordText = e.reifyWords(words, sep)
	}
	sentenceText := ""
	sentenceText = strings.Join(sentences, " ")
	errorText := ""
	if len(e.Errors) > 0 {
		errorText = e.reifyErrors(e.Errors)
	}

	return js.ValueOf(map[string]any{
		"words":          wordText,
		"sentences":      sentenceText,
		"errors":         errorText,
		"generatedCount": e.GeneratedCount,
		"duplicateCount": e.DuplicateCount,
		"rejectedCount":  e.RejectedCount,
		"replacedCount":  e.ReplacedCount,
	})
}

func New(opts Options, rs *rand.Rand) *Evaluator {
	e := new(Evaluator)
	e.Options = opts
	if rs == nil {
		// default random to current unix nano timestamp x2
		e.rs = rand.New(rand.NewPCG(uint64(time.Now().UnixNano()), uint64(time.Now().UnixNano())))
	} else {
		// use provided random
		e.rs = rs
	}
	e.resetCounters()
	return e
}

func (e *Evaluator) Run() (words []Word, syllableSep string, sentences []string) {
	// refesh the code input
	directives, err := e.LoadCode(e.Phonology)
	if err != nil {
		e.AddErrors(fmt.Errorf("parse error: %w", err))
		return
	}
	if !e.hasSyllableDirective(directives) {
		e.AddErrors(errors.New("at least one syllable directive is required"))
		return
	}

	e.evalDirectives(directives)
	if ok, err := e.checkCategories(); !ok {
		e.AddErrors(fmt.Errorf("category check error: %w", err))
		return
	}
	if ok, err := e.checkComponents(); !ok {
		e.AddErrors(fmt.Errorf("component check error: %w", err))
		return
	}

	if e.GenerateSentences {
		sentences = e.createSentences()
		return
	}
	// generate N words
	words = e.generateWords(e.WordCount * 2)
	// convert the words to lists of syllables
	words = e.syllabizeWords(words)
	// if on, remove duplicates
	words = e.removeDuplicates(words)
	if len(words) >= e.WordCount {
		words = words[:e.WordCount]
	}

	// if on, apply rejections
	// TODO: allow contexts in the middle of rejection elements
	words = e.rejectWords(words)

	// TODO: if on, apply replacements
	// words = e.replaceWords(words)

	// if on, force generate to WordCount
	// get number of possible syllables, and abort forced gen if possible < wanted
	count := e.choiceCount(e.categories, e.components)
	if e.ForceWordLimit && count >= e.WordCount {
		for len(words) < e.WordCount {
			words = e.generateWords(e.WordCount * 2)
			words = e.syllabizeWords(words)
			words = e.removeDuplicates(words)
			words = e.rejectWords(words)
			// TODO: apply replacements
			// words = e.replaceWords(words)
		}
		rand.Shuffle(len(words), func(i, j int) {
			words[i], words[j] = words[j], words[i]
		})
		if len(words) >= e.WordCount {
			words = words[:e.WordCount]
		}
	} else if e.ForceWordLimit && count < e.WordCount {
		e.AddErrors(fmt.Errorf("not enough choices to force word count: only %d/%d choices available", count, e.WordCount))
	}

	// if on, sort output
	if e.SortOutput {
		words = e.sort(words)
	}

	syllableSep = ""
	// TODO: if on, display with syllable separators
	if e.MarkSyllables {
		syllableSep = "."
	}

	return words, syllableSep, sentences
}

func (e *Evaluator) hasSyllableDirective(directives []ast.Directive) bool {
	hasSylDir := false
	for _, dir := range directives {
		switch dir.(type) {
		case *ast.SyllableDirective:
			hasSylDir = true
		default:
			continue
		}
	}
	return hasSylDir
}

func (e *Evaluator) LoadCode(src string) ([]ast.Directive, error) {
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
		err := word.GenerateSyllables(e.categories, e.components, e.rs)
		if err != nil {
			util.LogError(err.Error())
			e.AddErrors(err)
			return words
		}
		words[i] = word
	}
	return words
}

func (e *Evaluator) AddErrors(errs ...error) {
	e.Errors = append(e.Errors, errs...)
}

func (e *Evaluator) ClearErrors() {
	e.Errors = []error{}
}
func (e *Evaluator) removeDuplicates(words []Word) (ws []Word) {
	if !e.ForbidDuplicates {
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
	e.DuplicateCount = oldLen - len(ws)
	return ws
}

func (e *Evaluator) choiceCount(categories parts.Categories, components parts.Components) int {
	count := len(e.syllables)
	for _, s := range e.syllables {
		count *= s.ChoiceCount(categories, components)
	}
	return count
}

func (e *Evaluator) resetCounters() {
	e.GeneratedCount = 0
	e.DuplicateCount = 0
	e.RejectedCount = 0
	e.ReplacedCount = 0
}

func (e Evaluator) reifyWords(wds []Word, sylSep string) string {
	var b strings.Builder
	for _, word := range wds {
		b.WriteString(strings.Join(word.Syllables, sylSep) + "\n")
	}
	return strings.TrimSpace(b.String())
}

func (e Evaluator) reifyErrors(es []error) string {
	errs := errors.Join(es...)
	return fmt.Errorf("ERROR: execution failed due to:\n%w", errs).Error()
}
