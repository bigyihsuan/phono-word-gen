package syllable

import (
	"fmt"
	"regexp"
)

// Represents a component that can be evaluated
// (category, optional, selection, raw phoneme, etc)
type Evaluable interface {
	fmt.Stringer
	// Turn this into a list of phonemes
	Evaluate() []string
	// Generate all phoneme combinations possible from this component
	EvaluateAll() [][]string
}

// Matches a string to a component
type Matchable interface {
	// Returns whether the given string is an output of this component
	Matches(other string) bool
	Regexp() *regexp.Regexp
}

// Get a random choice for this component
type RandomlyChoosable interface {
	// Get a random list of phonemes
	RandomChoice() []string
}

type Component interface {
	componentTag()
}
