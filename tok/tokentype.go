package tok

//go:generate stringer -type=TokenType

type TokenType int

const (
	_ TokenType = iota
	EOF
	ILLEGAL
	COMMENT
	// common
	LPAREN
	RPAREN
	LBRACKET
	RBRACKET
	LBRACE
	RBRACE
	COMMA
	STAR
	COLON
	LINE_ENDING
	NUMBER
	// categories
	EQ
	DOLLAR
	RAW
	// context sigils
	CARET
	BSLASH
	AT
	AMPERSAND
	BANG
	// rejections
	PIPE
	// replacements
	ARROW
	SLASH
	DOUBLESLASH
	// keywords
	SYLLABLE
	LETTERS
	REJECT
	REPLACE
)

var keywords = map[string]TokenType{
	"syllable": SYLLABLE,
	"letters":  LETTERS,
	"reject":   REJECT,
	"replace":  REPLACE,
}

func IsKeywordOrRaw(lexeme string) TokenType {
	if tt, ok := keywords[lexeme]; ok {
		return tt
	} else {
		return RAW
	}
}
