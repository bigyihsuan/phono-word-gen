package ast

import "phono-word-gen/parts"

// a raw phoneme
type Phoneme struct {
	Value string
}

func (p *Phoneme) node()              {}
func (p *Phoneme) categoryElement()   {}
func (p *Phoneme) syllableComponent() {}
func (p *Phoneme) replacementSource() {}
func (p *Phoneme) String() string     { return p.Value }

// [Phoneme.PartElement] implements [CategoryElement].
func (p Phoneme) PartElement() (parts.Element, int) {
	return parts.NewPhoneme(p.Value), 1
}
