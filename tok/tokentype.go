package tok

//go:generate stringer -type=TokenType

type TokenType int

const (
	_ TokenType = iota
	EOF
	ILLEGAL
	// common
	LPAREN
	RPAREN
	LBRACKET
	RBRACKET
	LCURLY
	RCURLY
	COMMA
	STAR
	COLON
	LINE_ENDING
	NUMBER
	// categories
	EQ
	DOLLAR
	RAW
	// rejections
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
