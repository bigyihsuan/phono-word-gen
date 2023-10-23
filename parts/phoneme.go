package parts

import (
	"regexp"
)

type Phoneme struct {
	Value string
}

func NewPhoneme(value string) *Phoneme                   { return &Phoneme{Value: value} }
func (r *Phoneme) syllableElementTag()                   {}
func (r *Phoneme) Get(_ Categories) (string, error)      { return r.Value, nil }
func (r *Phoneme) ChoiceCount(categories Categories) int { return 1 }
func (r *Phoneme) Regexp(_ Categories) *regexp.Regexp {
	return regexp.MustCompile("(" + regexp.QuoteMeta(r.Value) + ")")
}
