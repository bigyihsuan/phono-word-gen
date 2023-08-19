package eval

import (
	"fmt"
	"phono-word-gen/ast"
	"phono-word-gen/parts"
	"phono-word-gen/util"
	"regexp"
	"strings"

	"github.com/mroth/weightedrand/v2"
	"golang.org/x/exp/slices"
)

func (evaluator *Evaluator) evalDirectives(directives []ast.Directive) {
	evaluator.categories = make(map[string]parts.Category)
	evaluator.syllables = []*parts.Syllable{}

	for _, dir := range directives {
		switch dir := dir.(type) {
		case *ast.CategoryDirective:
			evaluator.categories[dir.Name] = evaluator.evalCategory(dir)
		case *ast.SyllableDirective:
			evaluator.syllables = append(evaluator.syllables, evaluator.evalSyllable(dir))
		case *ast.LettersDirective:
			evaluator.letters = evaluator.evalLetters(dir)
		default:
			util.LogError(fmt.Sprintf("unknown directive: %T (%+v)\n", dir, dir))
		}
	}
}

func (evaluator *Evaluator) evalCategory(dir *ast.CategoryDirective) parts.Category {
	category := parts.Category{}
	for _, element := range dir.Phonemes {
		e, weight := evalCategoryElement(element)
		category.Elements = append(category.Elements, weightedrand.NewChoice(e, weight))
	}
	return category
}

func evalCategoryElement(element ast.CategoryElement) (e parts.Element, weight int) {
	weight = 1
	switch element := element.(type) {
	case *ast.Phoneme:
		return parts.NewPhoneme(element.Value), 1
	case *ast.Reference:
		return parts.NewReference(element.Name), 1
	case *ast.WeightedElement:
		weight = element.Weight
		e, _ = evalCategoryElement(element.Element)
	}
	return e, weight
}

func (evaluator *Evaluator) evalSyllable(dir *ast.SyllableDirective) *parts.Syllable {
	elements := []parts.SyllableElement{}
	for _, component := range dir.Components {
		elements = append(elements, evaluator.evalComponent(component))
	}
	return parts.NewSyllable(elements...)
}

func (evaluator *Evaluator) evalComponent(component ast.SyllableComponent) parts.SyllableElement {
	switch component := component.(type) {
	case *ast.Phoneme:
		return parts.NewRaw(component.Value)
	case *ast.Reference:
		return parts.NewReference(component.Name).(parts.SyllableElement)
	case *ast.SyllableGrouping:
		return evaluator.evalGrouping(component)
	case *ast.SyllableOptional:
		return evaluator.evalOptional(component)
	case *ast.SyllableSelection:
		return evaluator.evalSelection(component)
	default:
		fmt.Printf("unknown component: %T (%+v)\n", component, component)
	}
	return nil
}

func (evaluator *Evaluator) evalGrouping(component *ast.SyllableGrouping) parts.SyllableElement {
	components := []parts.SyllableElement{}
	for _, c := range component.Components {
		components = append(components, evaluator.evalComponent(c))
	}
	return parts.NewGrouping(components...)
}

func (evaluator *Evaluator) evalOptional(component *ast.SyllableOptional) parts.SyllableElement {
	components := []parts.SyllableElement{}
	for _, c := range component.Components {
		components = append(components, evaluator.evalComponent(c))
	}
	return parts.NewOptional(components, component.Weight)
}

func (evaluator *Evaluator) evalSelection(component *ast.SyllableSelection) parts.SyllableElement {
	components := []weightedrand.Choice[parts.SyllableElement, int]{}
	for _, c := range component.Components {
		comp, weight := evaluator.evalWeightedComponent(c.(*ast.WeightedSyllableComponent))
		choice := weightedrand.NewChoice(comp, weight)
		components = append(components, choice)
	}
	return parts.NewSelection(components...)
}

func (evaluator *Evaluator) evalWeightedComponent(component *ast.WeightedSyllableComponent) (parts.SyllableElement, int) {
	components := []parts.SyllableElement{}
	for _, c := range component.Components {
		components = append(components, evaluator.evalComponent(c))
	}
	return parts.NewGrouping(components...), component.Weight
}

func (e *Evaluator) evalLetters(dir *ast.LettersDirective) (letters []string) {
	for _, phoneme := range dir.Phonemes {
		letters = append(letters, phoneme.Value)
	}
	// make letter regexp
	ls := []string{}
	for _, l := range letters {
		ls = append(ls, "("+l+")")
	}
	slices.SortStableFunc(ls, func(a, b string) int { return len([]rune(b)) - len([]rune(a)) })
	e.letterRegexp = regexp.MustCompile(strings.Join(ls, "|"))
	return
}
