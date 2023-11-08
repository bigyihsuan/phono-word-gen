package par

import (
	"errors"
	"phono-word-gen/ast"
	"phono-word-gen/tok"
)

func (p *Parser) Replacement() *ast.ReplacementDirective {
	r := &ast.ReplacementDirective{}
	if !p.expectPeek(tok.COLON) {
		return nil
	}
	if !p.peekIs(tok.ARROW) {
		r.Source = p.ReplacementSource()
	}
	if !p.expectPeek(tok.ARROW) {
		return nil
	}
	if !p.peekIs(tok.SLASH) {
		r.Replacement = p.ReplacementReplacement()
	}
	if !p.expectPeek(tok.SLASH) {
		return nil
	}
	r.Condition = p.ReplacementEnv()
	if r.Condition == nil {
		p.errors = append(p.errors, errors.New("replacement requires condition env"))
		return nil
	}
	if p.peekIs(tok.DOUBLESLASH) {
		p.getNextToken()
		r.Exception = p.ReplacementEnv()
	}
	return r
}

func (p *Parser) ReplacementSource() (source []ast.ReplacementSource) {
	for p.peekIsAny(tok.DOLLAR, tok.RAW) {
		p.getNextToken()
		switch p.curr.Type {
		case tok.RAW:
			source = append(source, p.Phoneme())
		case tok.DOLLAR:
			source = append(source, p.CategoryReference())
		}
	}
	return
}

func (p *Parser) ReplacementReplacement() (replacement []*ast.Phoneme) {
	for p.peekIs(tok.RAW) {
		p.getNextToken()
		replacement = append(replacement, p.Phoneme())
	}
	return
}

func (p *Parser) ReplacementEnv() *ast.ReplacementEnv {
	r := &ast.ReplacementEnv{}
	if p.peekIsAny(tok.CARET, tok.AT, tok.BANG) {
		p.getNextToken()
		r.PrefixContext = p.curr.Lexeme
	}
	r.PrefixComponents = p.SyllableComponents()
	if !p.expectPeek(tok.UNDERSCORE) {
		return nil
	}
	// p.getNextToken()
	if p.peekIs(tok.DOUBLESLASH) {
		return r
	}
	r.SuffixComponents = p.SyllableComponents()
	if p.peekIsAny(tok.BSLASH, tok.AMPERSAND) {
		p.getNextToken()
		r.SuffixContext = p.curr.Lexeme
	}
	return r
}
