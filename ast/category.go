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
	Weight  *Weight
}

type Weight struct{ Value float64 }
