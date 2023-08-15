package eval

import (
	"phono-word-gen/parts"
)

type Word struct {
	Syllables []*parts.Syllable
}

func NewWord(syllables ...*parts.Syllable) *Word { return &Word{Syllables: syllables} }
func (w *Word) MakeSyllables(categories map[string]parts.Category) []string {
	syllables := []string{}
	for _, s := range w.Syllables {
		syllables = append(syllables, s.Get(categories))
	}
	return syllables
}
