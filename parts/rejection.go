package parts

import (
	"fmt"
	"regexp"
	"strings"
)

type Rejection struct {
	Prefix   RejectionPrefix
	Elements []SyllableElement
	Suffix   RejectionSuffix
}

func (r Rejection) IsWordLevel() bool {
	return r.Prefix == WORD_START || r.Suffix == WORD_END
}
func (r Rejection) IsSyllableLevel() bool {
	return r.Prefix == SYL_START || r.Suffix == SYL_END
}

func (r Rejection) Regexp(categories map[string]Category) *regexp.Regexp {
	elements := []string{}
	prefix, suffix := "", ""
	switch r.Prefix {
	case WORD_START, SYL_START:
		prefix = "^"
	}
	switch r.Suffix {
	case WORD_END, SYL_END:
		suffix = "$"
	}

	for _, e := range r.Elements {
		elements = append(elements, "("+e.Regexp(categories).String()+")")
	}
	return regexp.MustCompile(fmt.Sprintf("(%s%s%s)", prefix, strings.Join(elements, ""), suffix))
}

//go:generate stringer -type=RejectionPrefix
//go:generate stringer -type=RejectionSuffix
type RejectionPrefix int
type RejectionSuffix int

const (
	NO_PREFIX RejectionPrefix = iota
	WORD_START
	SYL_START
	NOT
)
const (
	NO_SUFFIX RejectionSuffix = iota
	WORD_END
	SYL_END
)
