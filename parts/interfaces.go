package parts

import "regexp"

type Element interface {
	Get(categories map[string]Category) (string, error)
	ChoiceCount(categories map[string]Category) int
	Regexp(categories map[string]Category) *regexp.Regexp
}

type SyllableElement interface {
	Element
	syllableElementTag()
}

var _ Element = &Raw{}
var _ Element = &Reference{}
var _ Element = &Grouping{}
var _ Element = &Selection{}
