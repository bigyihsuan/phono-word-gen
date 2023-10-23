package parts

import "regexp"

type Element interface {
	Get(categories Categories) (string, error)
	ChoiceCount(categories Categories) int
	Regexp(categories Categories) *regexp.Regexp
}

type SyllableElement interface {
	Element
	syllableElementTag()
}

var _ Element = &Phoneme{}
var _ Element = &Reference{}
var _ Element = &Grouping{}
var _ Element = &Selection{}
