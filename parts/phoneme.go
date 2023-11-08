package parts

import (
	"regexp"
)

type Phoneme struct {
	Value string
}

func NewPhoneme(value string) *Phoneme                            { return &Phoneme{Value: value} }
func (r *Phoneme) syllableElementTag()                            {}
func (r *Phoneme) Get(_ Categories, _ Components) (string, error) { return r.Value, nil }
func (r *Phoneme) ChoiceCount(_ Categories, _ Components) int     { return 1 }
func (r *Phoneme) Regexp(_ Categories, _ Components) *regexp.Regexp {
	return regexp.MustCompile("(" + regexp.QuoteMeta(r.Value) + ")")
}
