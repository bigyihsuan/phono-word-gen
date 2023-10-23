package parts

import (
	"regexp"
)

type Reference struct {
	Name string
}

func NewReference(name string) *Reference {
	return &Reference{Name: name}
}

func (r *Reference) syllableElementTag() {}
func (r *Reference) Get(categories Categories) (string, error) {
	// look for the existence for the category
	cat, ok := categories[r.Name]
	if ok {
		// if ok, get from the found category
		return cat.Get(categories)
	} else {
		return "", UndefinedCategoryError(r.Name, r.Name)
	}
}
func (r *Reference) ChoiceCount(categories Categories) int {
	cat, ok := categories[r.Name]
	if ok {
		return cat.ChoiceCount(categories)
	} else {
		return 0
	}
}
func (r *Reference) Regexp(categories Categories) *regexp.Regexp {
	cat, ok := categories[r.Name]
	if ok {
		return cat.Regexp(categories)
	} else {
		return regexp.MustCompile("")
	}
}
