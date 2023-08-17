package eval

import (
	"phono-word-gen/parts"
)

type Word struct {
	Syllables []*parts.Syllable
}

func NewWord(syllables ...*parts.Syllable) Word { return Word{Syllables: syllables} }
func (w Word) GenerateSyllables(categories map[string]parts.Category) ([]string, error) {
	syllables := []string{}
	for _, s := range w.Syllables {
		syl, err := s.Get(categories)
		if err != nil {
			return syllables, err
		}
		syllables = append(syllables, syl)
	}
	return syllables, nil
}
