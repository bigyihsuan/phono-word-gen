package ast

// a raw phoneme
type Phoneme struct {
	Value string
}

func (p *Phoneme) nodeTag()            {}
func (p *Phoneme) categoryElementTag() {}
