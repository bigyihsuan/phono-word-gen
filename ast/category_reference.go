package ast

type CategoryReference struct {
	Name string
}

func (r *CategoryReference) node()              {}
func (r *CategoryReference) categoryElement()   {}
func (r *CategoryReference) syllableComponent() {}
func (r *CategoryReference) replacementSource() {}
func (r *CategoryReference) String() string     { return "$" + r.Name }
