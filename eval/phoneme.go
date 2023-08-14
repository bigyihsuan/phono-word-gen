package eval

type Phoneme struct {
	Value string
}

func NewPhoneme(value string, weights ...float64) CategoryElement {
	p := &Phoneme{Value: value}
	return p
}

func (p *Phoneme) Get(_ map[string]Category) string { return p.Value }
