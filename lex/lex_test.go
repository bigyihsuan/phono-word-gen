package lex

import (
	"fmt"
	"phono-word-gen/tok"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetNextToken(t *testing.T) {
	tests := []struct {
		tt     tok.TokenType
		lexeme string
	}{
		{tok.RAW, "P"},
		{tok.EQ, "="},
		{tok.RAW, "p"},
		{tok.RAW, "t"},
		{tok.RAW, "k"},
		{tok.LINE_ENDING, "\n"},
		{tok.RAW, "R"},
		{tok.EQ, "="},
		{tok.RAW, "l"},
		{tok.RAW, "r"},
		{tok.RAW, "w"},
		{tok.RAW, "j"},
		{tok.LINE_ENDING, ";"},
		{tok.RAW, "C"},
		{tok.EQ, "="},
		{tok.DOLLAR, "$"},
		{tok.RAW, "P"},
		{tok.DOLLAR, "$"},
		{tok.RAW, "R"},
		{tok.RAW, "ŋ"},
		{tok.LINE_ENDING, "\n"},
		{tok.RAW, "V"},
		{tok.EQ, "="},
		{tok.RAW, "a"},
		{tok.STAR, "*"},
		{tok.NUMBER, "123"},
		{tok.RAW, "i"},
		{tok.STAR, "*"},
		{tok.NUMBER, "0.456"},
		{tok.RAW, "u"},
		{tok.RAW, "ə"},
		{tok.RAW, "ā"},
		{tok.LINE_ENDING, "\n"},
		{tok.SYLLABLE, "syllable"},
		{tok.COLON, ":"},
		{tok.LPAREN, "("},
		{tok.LBRACKET, "["},
		{tok.DOLLAR, "$"},
		{tok.RAW, "C"},
		{tok.STAR, "*"},
		{tok.NUMBER, "0.8"},
		{tok.COMMA, ","},
		{tok.DOLLAR, "$"},
		{tok.RAW, "C"},
		{tok.DOLLAR, "$"},
		{tok.RAW, "R"},
		{tok.RBRACKET, "]"},
		{tok.RPAREN, ")"},
		{tok.DOLLAR, "$"},
		{tok.RAW, "V"},
		{tok.LPAREN, "("},
		{tok.DOLLAR, "$"},
		{tok.RAW, "R"},
		{tok.RPAREN, ")"},
		{tok.LINE_ENDING, "\n"},
		{tok.LETTERS, "letters"},
		{tok.COLON, ":"},
		{tok.RAW, "a"},
		{tok.RAW, "i"},
		{tok.RAW, "j"},
		{tok.RAW, "k"},
		{tok.RAW, "l"},
		{tok.RAW, "p"},
		{tok.RAW, "r"},
		{tok.RAW, "t"},
		{tok.RAW, "w"},
		{tok.LINE_ENDING, "\n"},
		{tok.REJECT, "reject"},
		{tok.COLON, ":"},
		{tok.DOLLAR, "$"},
		{tok.RAW, "V"},
		{tok.DOLLAR, "$"},
		{tok.RAW, "V"},
		{tok.LINE_ENDING, "\n"},
		{tok.REPLACE, "replace"},
		{tok.COLON, ":"},
		{tok.LCURLY, "{"},
		{tok.RAW, "sourceA"},
		{tok.COMMA, ","},
		{tok.RAW, "sourceB"},
		{tok.RCURLY, "}"},
		{tok.ARROW, ">"},
		{tok.RAW, "substitute"},
		{tok.SLASH, "/"},
		{tok.RAW, "condition"},
		{tok.DOUBLESLASH, "//"},
		{tok.RAW, "optionalException"},
		{tok.EOF, ""},
		{tok.EOF, ""},
		{tok.EOF, ""},
		{tok.EOF, ""},
	}

	input := `P = p t k
R = l r w j; C=$P $R ŋ
V = a*123 i*0.456 u ə ā
syllable: ([$C*0.8, $C$R])$V ($R)
letters:  a i j k l p r t w
reject:   $V$V
replace:  {sourceA, sourceB} > substitute / condition // optionalException`

	l := New([]rune(input))

	for i, expected := range tests {
		actual := l.GetNextToken()
		if !assert.Equal(t, expected.tt, actual.Type,
			fmt.Sprintf("[%d] incorrect tokentype: expected=%q got=%q (%s)", i, expected.tt, actual.Type, actual.String())) {
			continue
		}
		assert.Equal(t, expected.lexeme, actual.Lexeme,
			fmt.Sprintf("[%d] incorrect lexeme: expected=%q got=%q (%s)", i, expected.lexeme, actual.Lexeme, actual.String()))
	}
}
