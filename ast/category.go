package ast

type Category struct {
	Name     string
	Phonemes []CategoryElement
}

type CategoryElement interface {
	categoryElementTag()
}

type WeightedElement struct {
	Element CategoryElement
	Weight  int
}

func (w *WeightedElement) categoryElementTag() {}
