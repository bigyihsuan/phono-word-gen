package par

import (
	"phono-word-gen/ast"
	"phono-word-gen/parts"
	"phono-word-gen/tok"
)

func (p *Parser) Syllable() *ast.SyllableDirective {
	if !p.expectPeek(tok.COLON) {
		return nil
	}
	components := p.SyllableComponents()
	if !p.expectPeek(tok.LINE_ENDING) {
		return nil
	}
	return &ast.SyllableDirective{Components: components}
}

func (p *Parser) SyllableComponents() (sc []ast.SyllableComponent) {
	for !p.peekIsAny(tok.LINE_ENDING, tok.RBRACE, tok.RPAREN, tok.RBRACKET, tok.COMMA, tok.STAR) {
		p.getNextToken()
		component := p.SyllableComponent()
		if component == nil {
			return
		}
		sc = append(sc, component)
	}
	return
}

func (p *Parser) SyllableComponent() ast.SyllableComponent {
	switch p.curr.Type {
	case tok.RAW:
		return p.Phoneme()
	case tok.DOLLAR:
		return p.Reference()
	case tok.LBRACE:
		return p.SyllableGrouping()
	case tok.LPAREN:
		return p.SyllableOptional()
	case tok.LBRACKET:
		return p.SyllableSelection()
	default:
		p.errors = append(p.errors, parts.UnexpectedToken(p.curr, p.curr.Type))
		return nil
	}
}

func (p *Parser) SyllableGrouping() ast.SyllableComponent {
	components := p.SyllableComponents()
	if !p.expectPeek(tok.RBRACE) {
		return nil
	}
	return &ast.SyllableGrouping{Components: components}
}

func (p *Parser) SyllableOptional() ast.SyllableComponent {
	components := p.SyllableComponents()
	if !p.expectPeek(tok.RPAREN) {
		return nil
	}
	weight := 50
	if p.peekIs(tok.STAR) {
		p.getNextToken() // rparen
		p.getNextToken() // star
		weight = p.Weight()
	}
	return &ast.SyllableOptional{Components: components, Weight: weight}
}

func (p *Parser) SyllableSelection() ast.SyllableComponent {
	components := p.SelectionElements()
	if components == nil {
		return nil
	}
	if !p.expectPeek(tok.RBRACKET) {
		return nil
	}
	return &ast.SyllableSelection{Components: components}
}
func (p *Parser) SelectionElements() (components []ast.SyllableComponent) {
	for !p.peekIs(tok.RBRACKET) {
		c := p.SelectionElement()
		components = append(components, c)
		if p.peekIs(tok.RBRACKET) {
			break
		}
		if !p.expectPeek(tok.COMMA) {
			return nil
		}
	}
	return
}

func (p *Parser) SelectionElement() *ast.WeightedSyllableComponent {
	components := p.SyllableComponents()
	if p.peekIs(tok.STAR) {
		p.getNextToken() // end of category element
		p.getNextToken() // star
		weight := p.Weight()
		if weight < 0 {
			return nil
		}
		return &ast.WeightedSyllableComponent{Components: components, Weight: weight}
	} else {
		return &ast.WeightedSyllableComponent{Components: components, Weight: 1}
	}
}
