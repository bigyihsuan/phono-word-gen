package parts

import (
	"testing"

	wr "github.com/mroth/weightedrand/v2"
	"github.com/stretchr/testify/assert"
)

var categories = map[string]Category{}

func TestSyllableRaw(t *testing.T) {
	raw := NewRaw("a")
	actual, err := raw.Get(categories)
	assert.Nil(t, err)
	assert.Equal(t, "a", actual)
}

func TestSyllableGrouping(t *testing.T) {
	grouping := NewGrouping(NewRaw("a"), NewRaw("b"), NewRaw("c"))
	actual, err := grouping.Get(categories)
	assert.Nil(t, err)
	assert.Equal(t, "abc", actual)
}

func TestSyllableOptional(t *testing.T) {
	optional := NewOptional([]SyllableElement{NewRaw("a"), NewRaw("b"), NewRaw("c")})
	for i := 0; i < 10; i++ {
		actual, err := optional.Get(categories)
		assert.Nil(t, err)
		assert.True(t, actual == "abc" || actual == "")
	}
}

func TestSyllableSelection(t *testing.T) {
	selection := NewSelection(
		wr.NewChoice[SyllableElement, int](NewRaw("a"), 1),
		wr.NewChoice[SyllableElement, int](NewRaw("b"), 1),
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
			NewRaw("b"),
			NewRaw("a"),
			NewSelection(
				wr.NewChoice[SyllableElement, int](NewRaw("n"), 1),
				wr.NewChoice[SyllableElement, int](NewRaw("d"), 1),
			),
		},
	}
	for i := 0; i < 10; i++ {
		actual, err := syllable.Get(categories)
		assert.Nil(t, err)
		assert.True(t, actual == "ban" || actual == "bad", "invalid output: want=%q/%q got=%q", "ban", "bad", actual)
	}
}
