package parts

type CategoryReference struct {
	Name string
}

func NewCategoryReference(name string) CategoryElement {
	return &CategoryReference{Name: name}
}

func (cr CategoryReference) Get(categories map[string]Category) string {
	// look for the existence for the category
	cat, ok := categories[cr.Name]
	if ok {
		// if ok, get from the found category
		return cat.Get(categories)
	} else {
		return ""
	}
}
