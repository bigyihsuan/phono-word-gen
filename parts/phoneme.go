package parts

import (
	"regexp"
)

type Phoneme struct {
	Value string
}

func NewPhoneme(value string) *Phoneme                            { return &Phoneme{Value: value} }
func (r *Phoneme) syllableElementTag()                            {}
func (r *Phoneme) Get(_ map[string]Category) (string, error)      { return r.Value, nil }
func (r *Phoneme) ChoiceCount(categories map[string]Category) int { return 1 }
func (r *Phoneme) Regexp(_ map[string]Category) *regexp.Regexp {
	return regexp.MustCompile("(" + regexp.QuoteMeta(r.Value) + ")")
}
