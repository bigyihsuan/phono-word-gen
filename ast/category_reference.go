package ast

import "phono-word-gen/parts"

type CategoryReference struct {
	Name string
}

func (r *CategoryReference) node()              {}
func (r *CategoryReference) categoryElement()   {}
func (r *CategoryReference) syllableComponent() {}
func (r *CategoryReference) replacementSource() {}
func (r *CategoryReference) String() string     { return "$" + r.Name }

// [CategoryReference.PartElement] implements [CategoryElement].
func (r CategoryReference) PartElement() (parts.Element, int) {
	return parts.NewCategoryReference(r.Name), 1
}
