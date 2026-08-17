package parts

import (
	"math/rand/v2"
	"strings"
)

type Components map[string]Component

type Component struct {
	Elements SyllableElements
}

func NewComponent(elements ...SyllableElement) *Component { return &Component{Elements: elements} }
func (s *Component) Get(categories Categories, components Components, rs *rand.Rand) (string, error) {
	elements := []string{}
	for _, e := range s.Elements {
		ele, err := e.Get(categories, components, rs)
		if err != nil {
			return ele, err
		}
		elements = append(elements, ele)
	}
	return strings.Join(elements, ""), nil
}
func (s *Component) ChoiceCount(categories Categories, components Components) int {
	count := 1
	for _, e := range s.Elements {
		count *= e.ChoiceCount(categories, components)
	}
	return count
}
