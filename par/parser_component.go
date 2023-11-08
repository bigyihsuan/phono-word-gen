package par

import (
	"phono-word-gen/ast"
	"phono-word-gen/tok"
)

func (p *Parser) Component() *ast.ComponentDirective {
	if !p.expectPeek(tok.COLON) {
		return nil
	}
	if !p.peekIs(tok.RAW) {
		return nil
	}
	p.getNextToken()
	name := p.curr.Lexeme
	if !p.expectPeek(tok.EQ) {
		return nil
	}
	components := p.SyllableComponents()
	if !p.expectPeek(tok.LINE_ENDING) {
		return nil
	}
	return &ast.ComponentDirective{Name: name, Components: components}
}
