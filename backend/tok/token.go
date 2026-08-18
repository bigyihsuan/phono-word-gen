package tok

import "fmt"

type Token struct {
	Type   TokenType
	Lexeme string
	Pos    Pos
}

func New(tt TokenType, lexeme string, pos Pos) Token {
	return Token{
		Type:   tt,
		Lexeme: lexeme,
		Pos:    pos,
	}
}

func (t *Token) String() string {
	return fmt.Sprintf("{%s %q @ %s}", t.Type, t.Lexeme, t.Pos)
}

// Pos is a representation of a token's location in the source phonology.
type Pos struct {
	Index int // starting index of the Token
	Line  int // the physical line of the Token
	Col   int // the physical start column of the Token
}

func (p Pos) String() string {
	return fmt.Sprintf("%d:%d (idx %d)", p.Line, p.Col, p.Index)
}
