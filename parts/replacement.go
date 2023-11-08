package parts

import (
	"fmt"
	"regexp"
	"strings"
)

type Replacement struct {
	Source      SyllableElements
	Replacement string
	Condition   ReplacementEnv
	Exception   *ReplacementEnv
}

func (r Replacement) ConditionRegexp(categories Categories, components Components) *regexp.Regexp {
	// PREFIX SOURCE SUFFIX
	prefix, suffix := r.Condition.Regexp(categories, components)
	return regexp.MustCompile(
		fmt.Sprintf("(%s(%s)%s)",
			prefix,
			r.Source.Regexp(categories, components),
			suffix,
		),
	)
}
func (r Replacement) ExceptionRegexp(categories Categories, components Components) *regexp.Regexp {
	// PREFIX SOURCE SUFFIX
	if r.Exception == nil {
		return nil
	}
	prefix, suffix := r.Exception.Regexp(categories, components)
	return regexp.MustCompile(
		fmt.Sprintf("(%s(%s)%s)",
			prefix,
			r.Source.Regexp(categories, components),
			suffix,
		),
	)
}

type ReplacementEnv struct {
	Prefix           ContextPrefix
	PrefixComponents SyllableElements
	SuffixComponents SyllableElements
	Suffix           ContextSuffix
}

func (r ReplacementEnv) IsWordLevel() bool {
	return r.Prefix == WORD_START || r.Suffix == WORD_END
}
func (r ReplacementEnv) IsSyllableLevel() bool {
	return r.Prefix == SYL_START || r.Suffix == SYL_END
}

func (r ReplacementEnv) Regexp(categories Categories, components Components) (prefix, suffix *regexp.Regexp) {
	prefixContext, suffixContext := "", ""
	switch r.Prefix {
	case WORD_START, SYL_START:
		prefixContext = "^"
	}
	switch r.Suffix {
	case WORD_END, SYL_END:
		suffixContext = "$"
	}

	prefixElements := []string{}
	for _, e := range r.PrefixComponents {
		prefixElements = append(prefixElements, e.Regexp(categories, components).String())
	}
	suffixElements := []string{}
	for _, e := range r.SuffixComponents {
		suffixElements = append(suffixElements, e.Regexp(categories, components).String())
	}
	prefixStr := prefixContext
	if len(prefixElements) > 0 {
		prefixStr = fmt.Sprintf("(%s%s)", prefixContext, strings.Join(prefixElements, ""))
	}
	suffixStr := suffixContext
	if len(prefixElements) > 0 {
		suffixStr = fmt.Sprintf("(%s%s)", strings.Join(suffixElements, ""), suffixContext)
	}
	return regexp.MustCompile(prefixStr), regexp.MustCompile(suffixStr)
}
