package par

import (
	"phono-word-gen/ast"
	"phono-word-gen/tok"
)

func (p *Parser) Letters() *ast.LettersDirective {
	l := &ast.LettersDirective{}
	if !p.expectPeek(tok.COLON) {
		return nil
	}
	for !p.peekIs(tok.LINE_ENDING) {
		p.getNextToken()
		l.Phonemes = append(l.Phonemes, p.Phoneme())
	}
	return l
}
