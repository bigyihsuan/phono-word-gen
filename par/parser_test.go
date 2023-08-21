package par

import (
	"phono-word-gen/ast"
	"phono-word-gen/lex"
	"testing"

	"github.com/stretchr/testify/assert"
)

func checkParseErrors(t *testing.T, p *Parser, i ...int) {
	if len(p.errors) == 0 {
		return
	}
	if len(i) > 0 {
		t.Errorf("[%d] parser has %d errors", i, len(p.errors))
	} else {
		t.Errorf("parser has %d errors", len(p.errors))
	}
	for _, err := range p.errors {
		t.Errorf("parser error: %s", err)
	}
	t.FailNow()
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
		phoneme := p.Phoneme()
		checkParseErrors(t, p, i)
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
		reference := p.Reference()
		checkParseErrors(t, p, i)
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
		weighted := p.WeightedCategoryElement()
		checkParseErrors(t, p, i)
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
		checkParseErrors(t, p, i)
		category, ok := directive.(*ast.CategoryDirective)
		if !assert.True(t, ok, "[%d] not a CategoryDirective: got=%T (%+v)", i, directive, directive) {
			continue
		}
		if !assert.NotNil(t, category, "[%d] category was nil") {
			continue
		}
		assert.Equal(t, tt.expected, category.String(),
			"[%d] incorrect: want=%q got=%q", i, tt.expected, category.String())
	}
}

func TestParseSyllableDirective(t *testing.T) {
	tests := []struct{ input, expected string }{
		{"syllable: [$C*9, $Cr$R, $Cl$L, s[c,t]]", "(syllable [($C * 9), ($Cr $R * 1), ($Cl $L * 1), (s [(c * 1), (t * 1)] * 1)])"},
		{"syllable: $C$V", "(syllable $C $V)"},
		{"syllable: {$C$V}", "(syllable {$C $V})"},
		{"syllable: ($C)*123$V", "(syllable (($C) * 123) $V)"},
		{"syllable: ($C)$V", "(syllable (($C) * 50) $V)"},
		{"syllable: [$C,$V]", "(syllable [($C * 1), ($V * 1)])"},
	}
	for i, tt := range tests {
		l := lex.New([]rune(tt.input))
		p := New(l)
		directive := p.Directive()
		checkParseErrors(t, p, i)
		syllable, ok := directive.(*ast.SyllableDirective)
		if !assert.True(t, ok, "[%d] not a SyllableDirective: got=%T (%+v)", i, directive, directive) {
			continue
		}
		if !assert.NotNil(t, syllable, "[%d] syllable was nil", i) {
			continue
		}
		assert.Equal(t, tt.expected, syllable.String(),
			"[%d] incorrect: want=%q got=%q", tt.expected, syllable.String())
	}
}

func TestParseLettersDirective(t *testing.T) {
	input := `letters: ñ ng p t k a i u`
	expected := "(letters [ñ ng p t k a i u])"
	l := lex.New([]rune(input))
	p := New(l)
	directive := p.Directive()
	checkParseErrors(t, p)
	letters, ok := directive.(*ast.LettersDirective)
	if !assert.True(t, ok, "not a LettersDirective: got=%T (%+v)", directive, directive) {
		return
	}
	if !assert.NotNil(t, letters, "letters was nil") {
		return
	}
	assert.Equal(t, expected, letters.String(),
		"incorrect: want=%q got=%q", expected, letters.String())
}

func TestParseRejectionDirective(t *testing.T) {
	input := `reject: ^$Vweak$Vweak | $R$C n | $Vstrong$Vstrong`
	expected := `(reject (^$Vweak $Vweak)|($R $C n)|($Vstrong $Vstrong))`
	l := lex.New([]rune(input))
	p := New(l)
	directive := p.Directive()
	checkParseErrors(t, p)
	reject, ok := directive.(*ast.RejectionDirective)
	if !assert.True(t, ok, "not a RejectionDirective: got=%T (%+v)", directive, directive) {
		return
	}
	if !assert.NotNil(t, reject, "reject was nil") {
		return
	}
	assert.Equal(t, expected, reject.String(),
		"incorrect: want=%q got=%q", expected, reject.String())
}
