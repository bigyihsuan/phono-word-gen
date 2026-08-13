package ast

import (
	"fmt"
	"strings"
)

type LettersDirective struct {
	Phonemes []*Phoneme
}

func (l *LettersDirective) node()      {}
func (l *LettersDirective) directive() {}
func (l *LettersDirective) String() string {
	phonemes := []string{}
	for _, p := range l.Phonemes {
		phonemes = append(phonemes, p.String())
	}
	return fmt.Sprintf("(letters [%s])", strings.Join(phonemes, " "))
}
