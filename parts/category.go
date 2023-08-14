package parts

import (
	"fmt"

	"github.com/mroth/weightedrand/v2"
)

type CategoryChoice = weightedrand.Choice[CategoryElement, int]
type CategoryChooser = *weightedrand.Chooser[CategoryElement, int]

type Category struct {
	Name     string
	Elements CategoryChooser
}

func NewCategory(name string, elements []weightedrand.Choice[CategoryElement, int]) Category {
	c := Category{Name: name}
	chooser, err := weightedrand.NewChooser[CategoryElement, int](elements...)
	if err != nil {
		fmt.Println(err)
	}
	c.Elements = chooser
	return c
}

func (c Category) Get(categories map[string]Category) string {
	// just pick something from the contained elements
	return c.Elements.Pick().Get(categories)
}
