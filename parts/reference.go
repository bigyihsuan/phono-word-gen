package parts

type Reference struct {
	Name string
}

func NewReference(name string) Element {
	return &Reference{Name: name}
}

func (r *Reference) syllableElementTag() {}
func (r *Reference) Get(categories map[string]Category) (string, error) {
	// look for the existence for the category
	cat, ok := categories[r.Name]
	if ok {
		// if ok, get from the found category
		return cat.Get(categories)
	} else {
		return "", nil
	}
}
