package parts

import (
	"testing"

	wr "github.com/mroth/weightedrand/v2"
	"github.com/stretchr/testify/assert"
)

var categories = map[string]Category{}

func TestSyllablePhoneme(t *testing.T) {
	raw := NewPhoneme("a")
	actual, err := raw.Get(categories)
	assert.Nil(t, err)
	assert.Equal(t, "a", actual)
}

func TestSyllableGrouping(t *testing.T) {
	grouping := NewGrouping(NewPhoneme("a"), NewPhoneme("b"), NewPhoneme("c"))
	actual, err := grouping.Get(categories)
	assert.Nil(t, err)
	assert.Equal(t, "abc", actual)
}

func TestSyllableOptional(t *testing.T) {
	optional := NewOptional([]SyllableElement{NewPhoneme("a"), NewPhoneme("b"), NewPhoneme("c")})
	for i := 0; i < 10; i++ {
		actual, err := optional.Get(categories)
		assert.Nil(t, err)
		assert.True(t, actual == "abc" || actual == "")
	}
}

func TestSyllableSelection(t *testing.T) {
	selection := NewSelection(
		wr.NewChoice[SyllableElement, int](NewPhoneme("a"), 1),
		wr.NewChoice[SyllableElement, int](NewPhoneme("b"), 1),
	)
	for i := 0; i < 10; i++ {
		actual, err := selection.Get(categories)
		assert.Nil(t, err)
		assert.True(t, actual == "a" || actual == "b", "invalid output: want=%q/%q got=%q", "a", "b", actual)
	}
}

func TestSyllableGet(t *testing.T) {
	syllable := Syllable{
		Elements: []SyllableElement{
			NewPhoneme("b"),
			NewPhoneme("a"),
			NewSelection(
				wr.NewChoice[SyllableElement, int](NewPhoneme("n"), 1),
				wr.NewChoice[SyllableElement, int](NewPhoneme("d"), 1),
			),
		},
	}
	for i := 0; i < 10; i++ {
		actual, err := syllable.Get(categories)
		assert.Nil(t, err)
		assert.True(t, actual == "ban" || actual == "bad", "invalid output: want=%q/%q got=%q", "ban", "bad", actual)
	}
}
