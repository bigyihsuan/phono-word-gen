package par

import (
	"phono-word-gen/ast"
	"phono-word-gen/lex"
	"testing"

	"github.com/stretchr/testify/assert"
)

func checkParseErrors(t *testing.T, p *Parser) {
	for _, err := range p.Errors() {
		t.Error(err)
	}
}

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
		checkParseErrors(t, p)
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
		checkParseErrors(t, p)
		reference, ok := ele.(*ast.Reference)
		if !assert.True(t, ok, "[%d] not a Reference: got=%T (%+v)", i, ele, ele) {
			continue
		}
		assert.Equal(t, tt.expected, reference.Name,
			"[%d] incorrect: want=%q got=%q", i, tt.expected, reference.Name)
	}
}

func TestParseWeightedCategoryElement(t *testing.T) {
	tests := []struct {
		input    string
		expected int
	}{
		{"a*123", 123},
		{"abc*456", 456},
		{"$ā*1", 1},
		{"$noWeight", 1},
		{"noWeight", 1},
	}

	for i, tt := range tests {
		l := lex.New([]rune(tt.input))
		p := New(l)
		ele := p.WeightedCategoryElement()
		checkParseErrors(t, p)
		weighted, ok := ele.(*ast.WeightedElement)
		if !assert.True(t, ok, "[%d] not a WeightedElement: got=%T (%+v)", i, ele, ele) {
			continue
		}
		assert.Equal(t, tt.expected, weighted.Weight,
			"[%d] incorrect: want=%d got=%d", i, tt.expected, weighted.Weight)
	}
}

func TestParseCategory(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"N = m n ñ", "(N = m*1 n*1 ñ*1)"},
		{"C = p*1 t*3", "(C = p*1 t*3)"},
		{"C = $N t*3", "(C = $N*1 t*3)"},
	}
	for i, tt := range tests {
		l := lex.New([]rune(tt.input))
		p := New(l)
		directive := p.Directive()
		checkParseErrors(t, p)
		category, ok := directive.(*ast.Category)
		if !assert.True(t, ok, "[%d] not a Category: got=%T (%+v)", i, directive, directive) {
			continue
		}
		if !assert.NotNil(t, category) {
			continue
		}
		assert.Equal(t, tt.expected, category.String(),
			"[%d] incorrect: want=%q got=%q", tt.expected, category.String())
	}
}
