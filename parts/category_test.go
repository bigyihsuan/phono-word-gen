package parts

import (
	"testing"

	wr "github.com/mroth/weightedrand/v2"
	"github.com/stretchr/testify/assert"
)

func TestCategoryGet(t *testing.T) {
	cat := NewCategory(wr.NewChoice[Element, int](NewPhoneme("p"), 1))
	for i := 0; i < 10; i++ {
		expected := "p"
		actual, _ := cat.Get(make(map[string]Category))
		assert.Equal(t, expected, actual, "[%d] incorrect: want=%q got=%q", i, expected, actual)
	}
}

func TestCategoryNestedGet(t *testing.T) {
	c := NewCategory(wr.NewChoice[Element, int](NewReference("S"), 1))
	s := NewCategory(wr.NewChoice[Element, int](NewPhoneme("p"), 1))
	categories := map[string]Category{"C": c, "S": s}
	cat := categories["C"]
	for i := 0; i < 10; i++ {
		expected := "p"
		actual, _ := cat.Get(categories)
		assert.Equal(t, expected, actual, "[%d] incorrect: want=%q got=%q", i, expected, actual)
	}
}
