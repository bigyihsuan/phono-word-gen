package eval

import (
	"fmt"

	"github.com/mroth/weightedrand/v2"
)

type Category struct {
	Name     string
	Elements *weightedrand.Chooser[CategoryElement, int]
}

func NewCategory(name string, elements []weightedrand.Choice[CategoryElement, int]) *Category {
	c := &Category{Name: name}
	chooser, err := weightedrand.NewChooser[CategoryElement, int](elements...)
	if err != nil {
		fmt.Println(err)
	}
	c.Elements = chooser
	return c
}

func (c *Category) Get() string {
	return c.Elements.Pick().Get()
}

type CategoryElement interface {
	ResolveCategories(categories map[string]Category) []CategoryElement
	Get() string
}
