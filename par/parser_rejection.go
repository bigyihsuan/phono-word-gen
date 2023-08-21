package par

import (
	"phono-word-gen/ast"
	"phono-word-gen/tok"
)

func (p *Parser) Rejection() *ast.RejectionDirective {
	r := &ast.RejectionDirective{}
	if !p.expectPeek(tok.COLON) {
		return nil
	}
	for !p.peekIs(tok.LINE_ENDING) {
		r.Elements = append(r.Elements, p.RejectionElement())
		if p.peekIs(tok.LINE_ENDING) {
			p.getNextToken()
			break
		}
		if !p.expectPeek(tok.PIPE) {
			return nil
		}
	}
	return r
}

func (p *Parser) RejectionElement() ast.RejectionElement {
	r := ast.RejectionElement{}
	if p.peekIsAny(tok.CARET, tok.AT, tok.BANG) {
		p.getNextToken()
		r.Prefix = &p.curr
	}
	r.Elements = p.SyllableComponents()
	if p.peekIsAny(tok.BSLASH, tok.AMPERSAND) {
		p.getNextToken()
		r.Suffix = &p.curr
	}
	return r
}
