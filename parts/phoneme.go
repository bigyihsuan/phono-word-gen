package parts

import (
	"regexp"
)

type Phoneme struct {
	Value string
}

func NewPhoneme(value string) Element {
	p := &Phoneme{Value: value}
	return p
}

func (p *Phoneme) Get(_ map[string]Category) (string, error) { return p.Value, nil }
func (p *Phoneme) ChoiceCount(_ map[string]Category) int     { return 1 }
func (p *Phoneme) Regexp(categories map[string]Category) *regexp.Regexp {
	return regexp.MustCompile(regexp.QuoteMeta(p.Value))
}
