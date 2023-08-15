package parts

import (
	"errors"
	"phono-word-gen/errs"

	wr "github.com/mroth/weightedrand/v2"
)

type CategoryChoice = wr.Choice[Element, int]
type CategoryChooser = *wr.Chooser[Element, int]

type Category struct {
	Elements []CategoryChoice
}

func NewCategory(elements ...wr.Choice[Element, int]) Category {
	return Category{Elements: elements}
}

func (c Category) Get(categories map[string]Category) (string, error) {
	// just pick something from the contained elements
	chooser, err := wr.NewChooser[Element, int](c.Elements...)
	if err != nil {
		return "", errors.Join(errs.CategoryCreationError, err)
	}
	return chooser.Pick().Get(categories)
}
