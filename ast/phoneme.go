package ast

// a raw phoneme
type Phoneme struct {
	Value string
}

func (p *Phoneme) node()              {}
func (p *Phoneme) categoryElement()   {}
func (p *Phoneme) syllableComponent() {}
func (p *Phoneme) replacementSource() {}
func (p *Phoneme) String() string     { return p.Value }
