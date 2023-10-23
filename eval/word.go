package eval

import (
	"phono-word-gen/parts"
)

type Word struct {
	SylTemplates []*parts.Syllable
	Syllables    []string
}

func NewWord(syllables ...*parts.Syllable) Word { return Word{SylTemplates: syllables} }
func (w *Word) GenerateSyllables(categories map[string]parts.Category) error {
	syllables := []string{}
	for _, s := range w.SylTemplates {
		syl, err := s.Get(categories)
		if err != nil {
			return err
		}
		syllables = append(syllables, syl)
	}
	w.Syllables = syllables
	return nil
}

func (w Word) Join() (word string, sylStartIndexes []int) {
	sylStart := 0
	for _, syl := range w.Syllables {
		word += syl
		sylStartIndexes = append(sylStartIndexes, sylStart)
		sylStart += len(syl)
	}
	return
}
