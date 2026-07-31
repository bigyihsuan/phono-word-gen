package ast

import (
	"fmt"
	"phono-word-gen/parts"
)

type WeightedElement struct {
	Element CategoryElement
	Weight  int
}

func (w *WeightedElement) node()              {}
func (w *WeightedElement) categoryElement()   {}
func (w *WeightedElement) syllableComponent() {}
func (w *WeightedElement) String() string {
	return fmt.Sprintf("%s*%d", w.Element.String(), w.Weight)
}

func (w WeightedElement) PartElement() (parts.Element, int) {
	ele, _ := w.Element.PartElement()
	return ele, w.Weight
}
