package ast

import (
	"fmt"
	"phono-word-gen/parts"
)

// All nodes implement the Node interface.
type Node interface {
	node()
}

// Directive represents a top-level [Node] in the AST hierarchy.
// Only certain Nodes are allowed to be Directives.
type Directive interface {
	Node
	directive()
}

// CategoryElement represents a [Node] that is placed in a category
// (e.g. "C = p t k" where "p t k" are three CategoryElements).
// Only certain Nodes are allowed to be CategoryElements.
type CategoryElement interface {
	Node
	fmt.Stringer
	categoryElement()
	// Element() turns this CategoryElement into a [parts.Element] for the evaluator.
	PartElement() (parts.Element, int)
}

type SyllableComponent interface {
	Node
	fmt.Stringer
	syllableComponent()
}
