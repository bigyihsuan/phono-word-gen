package parts

import (
	"regexp"
)

type CategoryReference struct {
	Name string
}

func NewCategoryReference(name string) *CategoryReference {
	return &CategoryReference{Name: name}
}

func (r *CategoryReference) syllableElementTag() {}
func (r *CategoryReference) Get(categories Categories, components Components) (string, error) {
	// look for the existence for the category
	cat, ok := categories[r.Name]
	if ok {
		// if ok, get from the found category
		return cat.Get(categories, components)
	} else {
		return "", UndefinedCategoryError(r.Name, r.Name)
	}
}
func (r *CategoryReference) ChoiceCount(categories Categories, components Components) int {
	cat, ok := categories[r.Name]
	if ok {
		return cat.ChoiceCount(categories, components)
	} else {
		return 0
	}
}
func (r *CategoryReference) Regexp(categories Categories, components Components) *regexp.Regexp {
	cat, ok := categories[r.Name]
	if ok {
		return cat.Regexp(categories, components)
	} else {
		return regexp.MustCompile("")
	}
}
