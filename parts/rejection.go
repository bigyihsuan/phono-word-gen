package parts

import (
	"fmt"
	"regexp"
)

type Rejection struct {
	Prefix   ContextPrefix
	Elements SyllableElements
	Suffix   ContextSuffix
}

func (r Rejection) IsWordLevel() bool {
	return r.Prefix == WORD_START || r.Suffix == WORD_END
}
func (r Rejection) IsSyllableLevel() bool {
	return r.Prefix == SYL_START || r.Suffix == SYL_END
}

func (r Rejection) Regexp(categories Categories, components Components) *regexp.Regexp {
	prefixContext, suffixContext := "", ""
	switch r.Prefix {
	case WORD_START, SYL_START:
		prefixContext = "^"
	}
	switch r.Suffix {
	case WORD_END, SYL_END:
		suffixContext = "$"
	}
	return regexp.MustCompile(fmt.Sprintf("(%s%s%s)", prefixContext, r.Elements.Regexp(categories, components), suffixContext))
}
