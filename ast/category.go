package ast

import (
	"fmt"
	"strings"
)

type Category struct {
	Name     string
	Phonemes []CategoryElement
}

func (c *Category) nodeTag()      {}
func (c *Category) directiveTag() {}
func (c *Category) String() string {
	phonemes := []string{}
	for _, p := range c.Phonemes {
		phonemes = append(phonemes, p.String())
	}
	return fmt.Sprintf("(%s = %s)", c.Name, strings.Join(phonemes, " "))
}

type CategoryElement interface {
	fmt.Stringer
	categoryElementTag()
}

type WeightedElement struct {
	Element CategoryElement
	Weight  int
}

func (w *WeightedElement) categoryElementTag() {}
func (w *WeightedElement) String() string {
	return fmt.Sprintf("%s*%d", w.Element.String(), w.Weight)
}
