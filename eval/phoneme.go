package eval

type Phoneme struct {
	Value string
}

func NewPhoneme(value string, weights ...float64) CategoryElement {
	p := &Phoneme{Value: value}
	return p
}

func (p *Phoneme) ResolveCategories(categories map[string]Category) []CategoryElement {
	return []CategoryElement{p}
}
func (p *Phoneme) Get() string { return p.Value }
