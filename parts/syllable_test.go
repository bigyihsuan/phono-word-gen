package parts

import (
	"testing"

	"github.com/mroth/weightedrand/v2"
	"github.com/stretchr/testify/assert"
)

var categories = map[string]Category{}

func TestSyllableRaw(t *testing.T) {
	raw := Raw{Value: "a"}
	assert.Equal(t, "a", raw.Get(categories))
}

func TestSyllableGrouping(t *testing.T) {
	grouping := Grouping{
		Elements: []SyllableElement{
			&Raw{Value: "a"},
			&Raw{Value: "b"},
			&Raw{Value: "c"},
		},
	}
	assert.Equal(t, "abc", grouping.Get(categories))
}

func TestSyllableOptional(t *testing.T) {
	optional, _ := NewOptional(
		[]SyllableElement{
			&Raw{Value: "a"},
			&Raw{Value: "b"},
			&Raw{Value: "c"},
		},
	)
	for i := 0; i < 10; i++ {
		actual := optional.Get(categories)
		assert.True(t, actual == "abc" || actual == "")
	}
}

func TestSyllableSelection(t *testing.T) {
	chooser, _ := weightedrand.NewChooser[SyllableElement, int](
		weightedrand.NewChoice[SyllableElement, int](&Raw{"a"}, 1),
		weightedrand.NewChoice[SyllableElement, int](&Raw{"b"}, 1),
	)
	selection := Selection{Choices: chooser}
	for i := 0; i < 10; i++ {
		actual := selection.Get(categories)
		assert.True(t, actual == "a" || actual == "b", "invalid output: want=%q/%q got=%q", "a", "b", actual)
	}
}

func TestSyllableGet(t *testing.T) {
	chooser, _ := weightedrand.NewChooser[SyllableElement, int](
		weightedrand.NewChoice[SyllableElement, int](&Raw{"n"}, 1),
		weightedrand.NewChoice[SyllableElement, int](&Raw{"d"}, 1),
	)
	syllable := Syllable{
		Elements: []SyllableElement{
			&Raw{Value: "b"},
			&Raw{Value: "a"},
			&Selection{chooser},
		},
	}
	for i := 0; i < 10; i++ {
		actual := syllable.Get(categories)
		assert.True(t, actual == "ban" || actual == "bad", "invalid output: want=%q/%q got=%q", "ban", "bad", actual)
	}
}
