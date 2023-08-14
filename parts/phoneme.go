package parts

type Phoneme struct {
	Value string
}

func NewPhoneme(value string, weights ...float64) Element {
	p := &Phoneme{Value: value}
	return p
}

func (p *Phoneme) Get(_ map[string]Category) string { return p.Value }
