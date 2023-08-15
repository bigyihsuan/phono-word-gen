package par

import (
	"phono-word-gen/ast"
	"phono-word-gen/errs"
	"phono-word-gen/lex"
	"phono-word-gen/tok"
)

type Parser struct {
	l          *lex.Lexer
	curr, peek tok.Token
	errors     []error
}

func New(l *lex.Lexer) *Parser {
	p := &Parser{l: l}
	// fill curr and peek
	p.getNextToken()
	p.getNextToken()
	return p
}

func (p *Parser) Errors() []error {
	return p.errors
}

func (p *Parser) getNextToken() {
	p.curr = p.peek
	p.peek = p.l.GetNextToken()
}

func (p *Parser) peekIs(tt tok.TokenType) bool { return p.peek.Type != tt }
func (p *Parser) expect(tt tok.TokenType) bool {
	if !p.peekIs(tt) {
		p.errors = append(p.errors, errs.ParserUnexpectedTokenError(p.curr.Type, tt))
		return false
	} else {
		p.getNextToken()
		return true
	}
}

func (p *Parser) Phoneme() ast.CategoryElement {
	if !p.peekIs(tok.RAW) {
		return nil
	}
	return &ast.Phoneme{Value: p.curr.Lexeme}
}

func (p *Parser) Reference() ast.CategoryElement {
	if !p.expect(tok.DOLLAR) {
		return nil
	}
	if !p.peekIs(tok.RAW) {
		return nil
	}
	return &ast.Reference{Name: p.curr.Lexeme}
}
