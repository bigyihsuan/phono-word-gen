package par

import (
	"phono-word-gen/ast"
	"phono-word-gen/parts"
	"phono-word-gen/tok"
)

func (p *Parser) Category() *ast.CategoryDirective {
	c := &ast.CategoryDirective{Name: p.curr.Lexeme, Phonemes: []ast.CategoryElement{}}
	if !p.expectPeek(tok.EQ) {
		return nil
	}
	for !p.peekIs(tok.LINE_ENDING) && !p.peekIs(tok.EOF) {
		p.getNextToken()
		element := p.WeightedCategoryElement()
		if element != nil {
			c.Phonemes = append(c.Phonemes, element)
		} else {
			return nil
		}
	}
	return c
}

func (p *Parser) WeightedCategoryElement() *ast.WeightedElement {
	ele := p.CategoryElement()
	if p.peekIs(tok.STAR) {
		p.getNextToken() // end of category element
		p.getNextToken() // star
		weight := p.Weight()
		if weight < 0 {
			return nil
		}
		return &ast.WeightedElement{Element: ele, Weight: weight}
	} else {
		return &ast.WeightedElement{Element: ele, Weight: 1}
	}
}

func (p *Parser) CategoryElement() ast.CategoryElement {
	switch p.curr.Type {
	case tok.RAW:
		return p.Phoneme()
	case tok.DOLLAR:
		return p.Reference()
	default:
		p.errors = append(p.errors, parts.UnexpectedToken(p.peek, p.peek.Type))
		return nil
	}
}
