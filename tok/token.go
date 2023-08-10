package tok

import "fmt"

type Token struct {
	Type   TokenType
	Lexeme string
	Index  int
}

func New(tt TokenType, lexeme string, start int) Token {
	return Token{
		Type:   tt,
		Lexeme: lexeme,
		Index:  start,
	}
}

func (t *Token) String() string {
	return fmt.Sprintf("{%s %q @ %d}", t.Type, t.Lexeme, t.Index)
}
