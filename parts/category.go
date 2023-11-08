package parts

import (
	"errors"
	"regexp"
	"strings"

	wr "github.com/mroth/weightedrand/v2"
)

type CategoryChoice = wr.Choice[Element, int]
type CategoryChooser = *wr.Chooser[Element, int]
type Categories map[string]Category

type Category struct {
	Elements []CategoryChoice
}

func NewCategory(elements ...wr.Choice[Element, int]) Category {
	return Category{Elements: elements}
}

func NewCategoryFromPhonemes(phonemes ...string) Category {
	c := Category{}
	for _, s := range phonemes {
		c.Elements = append(c.Elements, wr.NewChoice[Element, int](NewPhoneme(s), 1))
	}
	return c
}

func (c Category) Get(categories Categories, components Components) (string, error) {
	// just pick something from the contained elements
	chooser, err := wr.NewChooser[Element, int](c.Elements...)
	if err != nil {
		return "", errors.Join(CategoryCreationError, err)
	}
	return chooser.Pick().Get(categories, components)
}
func (c Category) ChoiceCount(categories Categories, components Components) int {
	return len(c.Elements)
}
func (c Category) Regexp(categories Categories, components Components) *regexp.Regexp {
	elements := []string{}
	for _, e := range c.Elements {
		elements = append(elements, e.Item.Regexp(categories, components).String())
	}
	return regexp.MustCompile("(" + strings.Join(elements, "|") + ")")
}
