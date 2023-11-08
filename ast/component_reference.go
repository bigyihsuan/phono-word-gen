package ast

type ComponentReference struct {
	Name string
}

func (r *ComponentReference) node()              {}
func (r *ComponentReference) categoryElement()   {}
func (r *ComponentReference) syllableComponent() {}
func (r *ComponentReference) replacementSource() {}
func (r *ComponentReference) String() string     { return "%" + r.Name }
