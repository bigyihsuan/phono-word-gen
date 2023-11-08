package parts

import (
	"regexp"
)

type ComponentReference struct {
	Name string
}

func NewComponentReference(name string) *ComponentReference {
	return &ComponentReference{Name: name}
}

func (r *ComponentReference) syllableElementTag() {}
func (r *ComponentReference) Get(categories Categories, components Components) (string, error) {
	// look for the existence for the category
	comp, ok := components[r.Name]
	if ok {
		// if ok, get from the found category
		return comp.Get(categories, components)
	} else {
		return "", UndefinedComponentError(r.Name, r.Name)
	}
}
func (r *ComponentReference) ChoiceCount(categories Categories, components Components) int {
	cat, ok := categories[r.Name]
	if ok {
		return cat.ChoiceCount(categories, components)
	} else {
		return 0
	}
}
func (r *ComponentReference) Regexp(categories Categories, components Components) *regexp.Regexp {
	cat, ok := categories[r.Name]
	if ok {
		return cat.Regexp(categories, components)
	} else {
		return regexp.MustCompile("")
	}
}
