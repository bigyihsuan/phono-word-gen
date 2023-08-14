package parts

import (
	"strings"

	"github.com/mroth/weightedrand/v2"
)

type Syllable struct {
	Elements []SyllableElement
}

func (s *Syllable) Get(categories map[string]Category) string {
	elements := []string{}
	for _, e := range s.Elements {
		elements = append(elements, e.Get(categories))
	}
	return strings.Join(elements, "")
}

type SyllableElement interface {
	Element
	syllableElementTag()
}

// A raw string in a syllable
type Raw struct {
	Value string
}

func (r *Raw) syllableElementTag()              {}
func (r *Raw) Get(_ map[string]Category) string { return r.Value }

type Grouping struct {
	Values []SyllableElement
}

func (g *Grouping) syllableElementTag() {}
func (g *Grouping) Get(categories map[string]Category) string {
	// evaluate all elements in the grouping
	values := []string{}
	for _, v := range g.Values {
		values = append(values, v.Get(categories))
	}
	return strings.Join(values, "")
}

type Selection struct {
	Choices *weightedrand.Chooser[SyllableElement, int]
}

func (s *Selection) syllableElementTag() {}
func (s *Selection) Get(catgories map[string]Category) string {
	// pick a random choice in the selection
	return s.Choices.Pick().Get(catgories)
}
