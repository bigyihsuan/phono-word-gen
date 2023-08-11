package eval

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPhonemeResolveCategories(t *testing.T) {
	tests := []struct {
		p        *Phoneme
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
		actual := tt.p.ResolveCategories(make(map[string]Category))
		if !assert.Len(t, actual, 1,
			"[%d] len incorrect, want=1 got=%d", i, len(actual)) {
			continue
		}
		if !assert.IsType(t, tt.p, actual[0],
			"[%d] type incorrect, want=%T got=%T (%+v)", i, tt.p, actual[0], actual[0]) {
			continue
		}
		if !assert.Equal(t, tt.expected, actual[0].Value,
			"[%d] value incorrect, want=%q got=%q", i, tt.expected, actual[0].Value) {
			continue
		}
	}
}

func TestPhonemeGet(t *testing.T) {
	tests := []struct {
		p        *Phoneme
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
		actual := tt.p.Get(nil)
		if !assert.Equal(t, tt.expected, actual, "[%d] get is incorrect, want=%q got=%q", i, tt.expected, actual) {
			continue
		}
	}
}
