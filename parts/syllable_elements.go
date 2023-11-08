package parts

import (
	"regexp"
	"strings"
)

type SyllableElements []SyllableElement

func (se SyllableElements) Regexp(categories Categories, components Components) *regexp.Regexp {
	elements := []string{}
	for _, e := range se {
		elements = append(elements, "("+e.Regexp(categories, components).String()+")")
	}
	return regexp.MustCompile(strings.Join(elements, ""))
}
