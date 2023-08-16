package ast

import (
	"fmt"
	"strings"
)

type CategoryDirective struct {
	Name     string
	Phonemes []CategoryElement
}

func (c *CategoryDirective) node()      {}
func (c *CategoryDirective) directive() {}
func (c *CategoryDirective) String() string {
	phonemes := []string{}
	for _, p := range c.Phonemes {
		phonemes = append(phonemes, p.String())
	}
	return fmt.Sprintf("(%s = %s)", c.Name, strings.Join(phonemes, " "))
}
