package lex

import (
	"phono-word-gen/tok"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetNextToken(t *testing.T) {
	tests := []struct {
		tt     tok.TokenType
		lexeme string
	}{
		{tok.COMMENT, "# this is a comment"},
		{tok.LINE_ENDING, "\n"},
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
		{tok.NUMBER, "456"},
		{tok.RAW, "u"},
		{tok.RAW, "ə"},
		{tok.RAW, "ā"},
		{tok.LINE_ENDING, "\n"},
		{tok.RAW, "C"},
		{tok.EQ, "="},
		{tok.RAW, "p͡ɸ"},
		{tok.RAW, "t͡s"},
		{tok.RAW, "k͡x"},
		{tok.RAW, "k͡xʷ"},
		{tok.LINE_ENDING, "\n"},
		{tok.RAW, "N"},
		{tok.EQ, "="},
		{tok.NUMBER, "1"},
		{tok.NUMBER, "2"},
		{tok.NUMBER, "3"},
		{tok.NUMBER, "4"},
		{tok.NUMBER, "5"},
		{tok.LINE_ENDING, "\n"},
		{tok.RAW, "T"},
		{tok.EQ, "="},
		{tok.RAW, "˥"},
		{tok.RAW, "˦"},
		{tok.RAW, "˧"},
		{tok.RAW, "˨"},
		{tok.RAW, "˩"},
		{tok.RAW, "˩˥˧"},
		{tok.RAW, "˧˥˩"},
		{tok.RAW, "˧˥˧"},
		{tok.RAW, "˩˧˩"},
		{tok.LINE_ENDING, "\n"},
		{tok.COMPONENT, "component"},
		{tok.COLON, ":"},
		{tok.RAW, "cv"},
		{tok.EQ, "="},
		{tok.DOLLAR, "$"},
		{tok.RAW, "C"},
		{tok.DOLLAR, "$"},
		{tok.RAW, "V"},
		{tok.LINE_ENDING, "\n"},
		{tok.SYLLABLE, "syllable"},
		{tok.COLON, ":"},
		{tok.LPAREN, "("},
		{tok.LBRACKET, "["},
		{tok.DOLLAR, "$"},
		{tok.RAW, "C"},
		{tok.STAR, "*"},
		{tok.NUMBER, "8"},
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
		{tok.COMMENT, "#comment"},
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
		{tok.BANG, "!"},
		{tok.DOLLAR, "$"},
		{tok.RAW, "V"},
		{tok.DOLLAR, "$"},
		{tok.RAW, "V"},
		{tok.PIPE, "|"},
		{tok.DOLLAR, "$"},
		{tok.RAW, "C"},
		{tok.DOLLAR, "$"},
		{tok.RAW, "C"},
		{tok.LINE_ENDING, "\n"},
		{tok.REPLACE, "replace"},
		{tok.COLON, ":"},
		{tok.LBRACE, "{"},
		{tok.RAW, "sourceA"},
		{tok.COMMA, ","},
		{tok.RAW, "sourceB"},
		{tok.RBRACE, "}"},
		{tok.ARROW, ">"},
		{tok.RAW, "substitute"},
		{tok.SLASH, "/"},
		{tok.CARET, "^"},
		{tok.UNDERSCORE, "_"},
		{tok.RAW, "condition"},
		{tok.BSLASH, "\\"},
		{tok.DOUBLESLASH, "//"},
		{tok.AT, "@"},
		{tok.RAW, "optionalException"},
		{tok.UNDERSCORE, "_"},
		{tok.AMPERSAND, "&"},
		{tok.LINE_ENDING, "\n"},
		{tok.EOF, ""},
		{tok.EOF, ""},
		{tok.EOF, ""},
		{tok.EOF, ""},
	}

	input := `# this is a comment
P = p t k
R = l r w j; C=$P $R ŋ
V = a*123 i*456 u ə ā
C = p͡ɸ t͡s k͡x k͡xʷ
N = 1 2 3 4 5
T = ˥ ˦ ˧ ˨ ˩ ˩˥˧ ˧˥˩ ˧˥˧ ˩˧˩
component: cv = $C $V
syllable: ([$C*8, $C$R])$V ($R)
#comment
letters:  a i j k l p r t w
reject:   !$V$V|$C$C
replace:  {sourceA, sourceB} > substitute / ^ _ condition\ // @optionalException _ &`

	l := New([]rune(input))

	for i, expected := range tests {
		actual := l.GetNextToken()
		if !assert.Equal(t, expected.tt, actual.Type,
			"[%d] incorrect tokentype: expected=%q got=%q (%s)", i, expected.tt, actual.Type, actual.String()) {
			t.Fatal()
		}
		assert.Equal(t, expected.lexeme, actual.Lexeme,
			"[%d] incorrect lexeme: expected=%q got=%q (%s)", i, expected.lexeme, actual.Lexeme, actual.String())
	}
}
