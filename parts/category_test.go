package parts

import (
	"testing"

	"github.com/mroth/weightedrand/v2"
	"github.com/stretchr/testify/assert"
)

func TestCategoryGet(t *testing.T) {
	cat := NewCategory("C", []weightedrand.Choice[CategoryElement, int]{weightedrand.NewChoice[CategoryElement, int](NewPhoneme("p"), 1)})
	for i := 0; i < 10; i++ {
		expected := "p"
		actual := cat.Get(make(map[string]Category))
		assert.Equal(t, expected, actual, "[%d] incorrect: want=%q got=%q", i, expected, actual)
	}
}

func TestCategoryNestedGet(t *testing.T) {
	categories := map[string]Category{
		"C": NewCategory("C", []weightedrand.Choice[CategoryElement, int]{weightedrand.NewChoice[CategoryElement, int](NewCategoryReference("S"), 1)}),
		"S": NewCategory("S", []weightedrand.Choice[CategoryElement, int]{weightedrand.NewChoice[CategoryElement, int](NewPhoneme("p"), 1)}),
	}
	cat := categories["C"]
	for i := 0; i < 10; i++ {
		expected := "p"
		actual := cat.Get(categories)
		assert.Equal(t, expected, actual, "[%d] incorrect: want=%q got=%q", i, expected, actual)
	}
}
