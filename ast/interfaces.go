package ast

import "fmt"

type Node interface {
	node()
}

type Directive interface {
	Node
	directive()
}

type CategoryElement interface {
	Node
	fmt.Stringer
	categoryElement()
}

type SyllableComponent interface {
	Node
	fmt.Stringer
	syllableComponent()
}
