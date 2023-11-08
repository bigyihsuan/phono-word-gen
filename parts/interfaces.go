package parts

import "regexp"

type Element interface {
	Get(categories Categories, components Components) (string, error)
	ChoiceCount(categories Categories, components Components) int
	Regexp(categories Categories, components Components) *regexp.Regexp
}

type SyllableElement interface {
	Element
	syllableElementTag()
}

var _ Element = &Phoneme{}
var _ Element = &CategoryReference{}
var _ Element = &ComponentReference{}
var _ Element = &Grouping{}
var _ Element = &Selection{}
var _ Element = &Optional{}
