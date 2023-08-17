package parts

type Phoneme struct {
	Value string
}

func NewPhoneme(value string) Element {
	p := &Phoneme{Value: value}
	return p
}

func (p *Phoneme) Get(_ map[string]Category) (string, error) { return p.Value, nil }
