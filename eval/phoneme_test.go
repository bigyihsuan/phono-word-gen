package eval

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPhonemeGet(t *testing.T) {
	tests := []struct {
		p        CategoryElement
		expected string
	}{
		{NewPhoneme("a"), "a"},
		{NewPhoneme("i"), "i"},
		{NewPhoneme("j"), "j"},
		{NewPhoneme("k"), "k"},
		{NewPhoneme("l"), "l"},
		{NewPhoneme("p"), "p"},
		{NewPhoneme("r"), "r"},
		{NewPhoneme("t"), "t"},
		{NewPhoneme("w"), "w"},
		{NewPhoneme("ə"), "ə"},
		{NewPhoneme("ā"), "ā"},
		{NewPhoneme("ā"), "ā"},
		{NewPhoneme("āabc"), "āabc"},
	}

	for i, tt := range tests {
		actual := tt.p.Get(make(map[string]Category))
		if !assert.Equal(t, tt.expected, actual, "[%d] get is incorrect, want=%q got=%q", i, tt.expected, actual) {
			continue
		}
	}
}
