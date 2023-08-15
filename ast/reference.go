package ast

type Reference struct {
	Name string
}

func (r *Reference) nodeTag()            {}
func (r *Reference) categoryElementTag() {}
