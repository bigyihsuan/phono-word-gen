package ast

import (
	"fmt"
	"strings"
)

type ComponentDirective struct {
	Name       string
	Components []SyllableComponent
}

func (c *ComponentDirective) node()      {}
func (c *ComponentDirective) directive() {}
func (c *ComponentDirective) String() string {
	phonemes := []string{}
	for _, p := range c.Components {
		phonemes = append(phonemes, p.String())
	}
	return fmt.Sprintf("(component %s = %s)", c.Name, strings.Join(phonemes, " "))
}
