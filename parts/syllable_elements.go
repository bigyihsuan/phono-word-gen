package parts

import (
	"regexp"
	"strings"
)

type SyllableElements []SyllableElement

func (se SyllableElements) Regexp(categories Categories) *regexp.Regexp {
	elements := []string{}
	for _, e := range se {
		elements = append(elements, "("+e.Regexp(categories).String()+")")
	}
	return regexp.MustCompile(strings.Join(elements, ""))
}
