package parts

import (
	"errors"
	"phono-word-gen/errs"

	wr "github.com/mroth/weightedrand/v2"
)

type CategoryChoice = wr.Choice[Element, int]
type CategoryChooser = *wr.Chooser[Element, int]

type Category struct {
	Elements CategoryChooser
}

func NewCategory(elements ...wr.Choice[Element, int]) (Category, error) {
	c := Category{}
	chooser, err := wr.NewChooser[Element, int](elements...)
	if err != nil {
		return c, errors.Join(errs.CategoryCreationError, err)
	}
	c.Elements = chooser
	return c, nil
}

func (c Category) Get(categories map[string]Category) string {
	// just pick something from the contained elements
	return c.Elements.Pick().Get(categories)
}
