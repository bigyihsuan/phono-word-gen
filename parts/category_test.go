package parts

import (
	"testing"

	wr "github.com/mroth/weightedrand/v2"
	"github.com/stretchr/testify/assert"
)

func TestCategoryGet(t *testing.T) {
	cat, _ := NewCategory([]wr.Choice[Element, int]{wr.NewChoice[Element, int](NewPhoneme("p"), 1)})
	for i := 0; i < 10; i++ {
		expected := "p"
		actual := cat.Get(make(map[string]Category))
		assert.Equal(t, expected, actual, "[%d] incorrect: want=%q got=%q", i, expected, actual)
	}
}

func TestCategoryNestedGet(t *testing.T) {
	c, _ := NewCategory([]wr.Choice[Element, int]{wr.NewChoice[Element, int](NewReference("S"), 1)})
	s, _ := NewCategory([]wr.Choice[Element, int]{wr.NewChoice[Element, int](NewPhoneme("p"), 1)})
	categories := map[string]Category{"C": c, "S": s}
	cat := categories["C"]
	for i := 0; i < 10; i++ {
		expected := "p"
		actual := cat.Get(categories)
		assert.Equal(t, expected, actual, "[%d] incorrect: want=%q got=%q", i, expected, actual)
	}
}
