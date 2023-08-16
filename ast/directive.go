package ast

type Directive interface {
	Node
	directiveTag()
}
