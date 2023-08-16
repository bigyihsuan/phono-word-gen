package ast

type Reference struct {
	Name string
}

func (r *Reference) nodeTag()            {}
func (r *Reference) categoryElementTag() {}
func (r *Reference) String() string      { return "$" + r.Name }
