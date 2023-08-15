package par

import (
	"phono-word-gen/ast"
	"phono-word-gen/lex"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParsePhoneme(t *testing.T) {
	tests := []struct{ input, expected string }{
		{"a", "a"},
		{"abc", "abc"},
		{"ā", "ā"},
		{"ə", "ə"},
	}

	for i, tt := range tests {
		l := lex.New([]rune(tt.input))
		p := New(l)
		ele := p.Phoneme()
		phoneme, ok := ele.(*ast.Phoneme)
		if !assert.True(t, ok, "[%d] not a Phoneme: got=%T (%+v)", i, ele, ele) {
			continue
		}
		assert.Equal(t, tt.expected, phoneme.Value,
			"[%d] incorrect: want=%q got=%q", i, tt.expected, phoneme.Value)
	}
}

func TestParseReference(t *testing.T) {
	tests := []struct{ input, expected string }{
		{"$a", "a"},
		{"$abc", "abc"},
		{"$ā", "ā"},
		{"$ə", "ə"},
	}

	for i, tt := range tests {
		l := lex.New([]rune(tt.input))
		p := New(l)
		ele := p.Reference()
		reference, ok := ele.(*ast.Reference)
		if !assert.True(t, ok, "[%d] not a Reference: got=%T (%+v)", i, ele, ele) {
			continue
		}
		assert.Equal(t, tt.expected, reference.Name,
			"[%d] incorrect: want=%q got=%q", i, tt.expected, reference.Name)
	}
}
