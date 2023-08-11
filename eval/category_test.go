package eval

import (
	"math/rand"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCategoryRandomChoice(t *testing.T) {
	random := rand.New(rand.NewSource(1))
	tests := []string{"k", "p", "k", "k", "t", "p", "t"}
	category := NewCategory("C", []CategoryElement{NewPhoneme("p"), NewPhoneme("t"), NewPhoneme("k")})
	for i, expected := range tests {
		actual := category.Get(random)
		if !assert.Equal(t, expected, actual,
			"[%d] phoneme is incorrect, want=%q got=%q", i, expected, actual) {
			continue
		}
	}
}
