package ast

import (
	"fmt"
	"strings"
)

type ReplacementSource interface {
	fmt.Stringer
	replacementSource()
}

var _ ReplacementSource = &Phoneme{}
var _ ReplacementSource = &Reference{}

type ReplacementDirective struct {
	Source      []ReplacementSource
	Replacement []*Phoneme
	Condition   *ReplacementEnv
	Exception   *ReplacementEnv
}

func (rd *ReplacementDirective) node()      {}
func (rd *ReplacementDirective) directive() {}
func (rd *ReplacementDirective) String() string {
	source := []string{}
	for _, e := range rd.Source {
		source = append(source, e.String())
	}

	replacement := []string{}
	for _, e := range rd.Replacement {
		replacement = append(replacement, e.String())
	}

	condition := rd.Condition.String()

	exception := "()"
	if rd.Exception != nil {
		exception = rd.Exception.String()
	}

	return fmt.Sprintf("(replace (%s) > (%s) / %s // %s)", strings.Join(source, " "), strings.Join(replacement, " "), condition, exception)
}

type ReplacementEnv struct {
	PrefixContext    string
	PrefixComponents []SyllableComponent
	SuffixComponents []SyllableComponent
	SuffixContext    string
}

func (re ReplacementEnv) String() string {
	prefix := []string{}
	for _, e := range re.PrefixComponents {
		prefix = append(prefix, e.String())
	}
	suffix := []string{}
	for _, e := range re.SuffixComponents {
		suffix = append(suffix, e.String())
	}

	p := strings.Join(prefix, " ")
	if re.PrefixContext != "" {
		p = re.PrefixContext + " " + p
	}
	if len(p) > 0 {
		p += " "
	}
	s := strings.Join(suffix, " ")
	if re.SuffixContext != "" {
		s += " " + re.SuffixContext
	}
	if len(s) > 0 && s[0] != ' ' {
		s = " " + s
	}
	return fmt.Sprintf("(%s_%s)", p, s)
}
