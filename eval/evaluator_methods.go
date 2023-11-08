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

func (e *Evaluator) evalDirectives(directives []ast.Directive) {
	e.categories = make(map[string]parts.Category)
	e.syllables = []*parts.Syllable{}
	rejections := []parts.Rejection{}
	e.replacements = []parts.Replacement{}
	for _, dir := range directives {
		switch dir := dir.(type) {
		case *ast.CategoryDirective:
			e.categories[dir.Name] = e.evalCategory(dir)
		case *ast.SyllableDirective:
			e.syllables = append(e.syllables, e.evalSyllable(dir))
		case *ast.LettersDirective:
			e.letters = e.evalLetters(dir)
		case *ast.RejectionDirective:
			rejections = append(rejections, e.evalRejection(dir)...)
		case *ast.ReplacementDirective:
			e.replacements = append(e.replacements, e.evalReplacement(dir))
		default:
			util.LogError(fmt.Sprintf("unknown directive: %T (%+v)\n", dir, dir))
		}
	}

	e.processRejections(rejections)
}

func (e *Evaluator) evalCategory(dir *ast.CategoryDirective) parts.Category {
	category := parts.Category{}
	for _, element := range dir.Phonemes {
		e, weight := e.evalCategoryElement(element)
		category.Elements = append(category.Elements, weightedrand.NewChoice(e, weight))
	}
	return category
}

func (e *Evaluator) evalCategoryElement(element ast.CategoryElement) (ele parts.Element, weight int) {
	weight = 1
	switch element := element.(type) {
	case *ast.Phoneme:
		return parts.NewPhoneme(element.Value), 1
	case *ast.CategoryReference:
		return parts.NewReference(element.Name), 1
	case *ast.WeightedElement:
		weight = element.Weight
		ele, _ = e.evalCategoryElement(element.Element)
	}
	return ele, weight
}

func (e *Evaluator) evalSyllable(dir *ast.SyllableDirective) *parts.Syllable {
	elements := []parts.SyllableElement{}
	for _, component := range dir.Components {
		elements = append(elements, e.evalComponent(component))
	}
	return parts.NewSyllable(elements...)
}

func (e *Evaluator) evalComponent(component ast.SyllableComponent) parts.SyllableElement {
	switch component := component.(type) {
	case *ast.Phoneme:
		return parts.NewPhoneme(component.Value)
	case *ast.CategoryReference:
		return parts.NewReference(component.Name)
	case *ast.SyllableGrouping:
		return e.evalGrouping(component)
	case *ast.SyllableOptional:
		return e.evalOptional(component)
	case *ast.SyllableSelection:
		return e.evalSelection(component)
	default:
		fmt.Printf("unknown component: %T (%+v)\n", component, component)
	}
	return nil
}

func (e *Evaluator) evalGrouping(component *ast.SyllableGrouping) parts.SyllableElement {
	components := []parts.SyllableElement{}
	for _, c := range component.Components {
		components = append(components, e.evalComponent(c))
	}
	return parts.NewGrouping(components...)
}

func (e *Evaluator) evalOptional(component *ast.SyllableOptional) parts.SyllableElement {
	components := []parts.SyllableElement{}
	for _, c := range component.Components {
		components = append(components, e.evalComponent(c))
	}
	return parts.NewOptional(components, component.Weight)
}

func (e *Evaluator) evalSelection(component *ast.SyllableSelection) parts.SyllableElement {
	components := []weightedrand.Choice[parts.SyllableElement, int]{}
	for _, c := range component.Components {
		comp, weight := e.evalWeightedComponent(c.(*ast.WeightedSyllableComponent))
		choice := weightedrand.NewChoice(comp, weight)
		components = append(components, choice)
	}
	return parts.NewSelection(components...)
}

func (e *Evaluator) evalWeightedComponent(component *ast.WeightedSyllableComponent) (parts.SyllableElement, int) {
	components := []parts.SyllableElement{}
	for _, c := range component.Components {
		components = append(components, e.evalComponent(c))
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

func (e *Evaluator) evalRejection(dir *ast.RejectionDirective) []parts.Rejection {
	r := []parts.Rejection{}
	for _, reject := range dir.Elements {
		r = append(r, e.evalRejectionElement(reject))
	}
	return r
}

func (e *Evaluator) evalRejectionElement(dir ast.RejectionElement) parts.Rejection {
	r := parts.Rejection{}
	switch {
	case dir.PrefixContext == "^":
		r.Prefix = parts.WORD_START
	case dir.PrefixContext == "@":
		r.Prefix = parts.SYL_START
	case dir.PrefixContext == "!":
		r.Prefix = parts.NOT
	default:
		r.Prefix = parts.NO_PREFIX
	}

	for _, element := range dir.Elements {
		r.Elements = append(r.Elements, e.evalComponent(element))
	}

	switch {
	case dir.SuffixContext == "\\":
		r.Suffix = parts.WORD_END
	case dir.SuffixContext == "&":
		r.Suffix = parts.SYL_END
	default:
		r.Suffix = parts.NO_SUFFIX
	}
	return r
}

func (e *Evaluator) processRejections(rejections []parts.Rejection) {
	wordOnlyRejections := []parts.Rejection{}
	syllableOnlyRejections := []parts.Rejection{}
	generalRejections := []parts.Rejection{}
	bothRejections := []parts.Rejection{}
	// sort out rejections
	for _, r := range rejections {
		switch {
		case r.IsWordLevel() && r.IsSyllableLevel():
			bothRejections = append(bothRejections, r)
		case !r.IsWordLevel() && !r.IsSyllableLevel():
			generalRejections = append(generalRejections, r)
		case r.IsWordLevel():
			wordOnlyRejections = append(wordOnlyRejections, r)
		case r.IsSyllableLevel():
			syllableOnlyRejections = append(syllableOnlyRejections, r)
		}
	}

	e.wordRejections = e.mergeRejections(wordOnlyRejections)
	e.syllableRejections = e.mergeRejections(syllableOnlyRejections)
	e.generalRejections = e.mergeRejections(generalRejections)
}

func (e *Evaluator) mergeRejections(rejections []parts.Rejection) *regexp.Regexp {
	w := []string{}
	for _, r := range rejections {
		reg := r.Regexp(e.categories).String()
		w = append(w, reg)
	}
	return regexp.MustCompile(strings.Join(w, "|"))
}

func (e *Evaluator) evalReplacement(dir *ast.ReplacementDirective) parts.Replacement {
	r := parts.Replacement{}

	for _, s := range dir.Source {
		r.Source = append(r.Source, e.evalComponent(s.(ast.SyllableComponent)))
	}

	r.Replacement = ""
	for _, p := range dir.Replacement {
		r.Replacement += p.Value
	}

	r.Condition = *e.evalReplacementEnv(dir.Condition)
	if dir.Exception != nil {
		r.Exception = e.evalReplacementEnv(dir.Exception)
	}

	return r
}

func (e *Evaluator) evalReplacementEnv(env *ast.ReplacementEnv) *parts.ReplacementEnv {
	r := &parts.ReplacementEnv{}
	switch {
	case env.PrefixContext == "^":
		r.Prefix = parts.WORD_START
	case env.PrefixContext == "@":
		r.Prefix = parts.SYL_START
	case env.PrefixContext == "!":
		r.Prefix = parts.NOT
	default:
		r.Prefix = parts.NO_PREFIX
	}

	for _, element := range env.PrefixComponents {
		r.PrefixComponents = append(r.PrefixComponents, e.evalComponent(element))
	}
	for _, element := range env.SuffixComponents {
		r.SuffixComponents = append(r.SuffixComponents, e.evalComponent(element))
	}

	switch {
	case env.SuffixContext == "\\":
		r.Suffix = parts.WORD_END
	case env.SuffixContext == "&":
		r.Suffix = parts.SYL_END
	default:
		r.Suffix = parts.NO_SUFFIX
	}
	return r
}
