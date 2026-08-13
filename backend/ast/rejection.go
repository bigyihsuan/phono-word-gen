package ast

import (
	"fmt"
	"strings"
)

type RejectionDirective struct {
	Elements []RejectionElement
}

func (rd *RejectionDirective) node()      {}
func (rd *RejectionDirective) directive() {}
func (rd *RejectionDirective) String() string {
	elements := []string{}
	for _, e := range rd.Elements {
		elements = append(elements, e.String())
	}
	return fmt.Sprintf("(reject %s)", strings.Join(elements, "|"))
}

type RejectionElement struct {
	PrefixContext string
	Elements      []SyllableComponent
	SuffixContext string
}

func (re RejectionElement) String() string {
	elements := []string{}
	for _, e := range re.Elements {
		elements = append(elements, e.String())
	}
	return fmt.Sprintf("(%s%s%s)", re.PrefixContext, strings.Join(elements, " "), re.SuffixContext)
}
