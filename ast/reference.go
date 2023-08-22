package ast

type Reference struct {
	Name string
}

func (r *Reference) node()              {}
func (r *Reference) categoryElement()   {}
func (r *Reference) syllableComponent() {}
func (r *Reference) replacementSource() {}
func (r *Reference) String() string     { return "$" + r.Name }
