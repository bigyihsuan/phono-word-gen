package par

import (
	"fmt"
	"phono-word-gen/ast"
	"phono-word-gen/errs"
	"phono-word-gen/lex"
	"phono-word-gen/tok"
	"strconv"
)

type Parser struct {
	l          *lex.Lexer
	curr, peek tok.Token
	errors     []error
}

func New(l *lex.Lexer) *Parser {
	p := &Parser{l: l, errors: []error{}}
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

func (p *Parser) currIs(tt tok.TokenType) bool { return p.curr.Type == tt }
func (p *Parser) peekIs(tt tok.TokenType) bool { return p.peek.Type == tt }
func (p *Parser) expect(tt tok.TokenType) bool {
	if p.peekIs(tt) {
		p.getNextToken()
		return true
	} else {
		p.errors = append(p.errors, errs.ParserUnexpectedToken(p.curr.Type, tt))
		return false
	}
}

func (p *Parser) Directive() ast.Directive {
	// skip duplicate line endings
	for p.currIs(tok.LINE_ENDING) {
		p.getNextToken()
	}
	switch p.curr.Type {
	case tok.RAW:
		return p.Category()
	default:
		p.errors = append(p.errors, errs.ParserUnknownDirective(p.curr.Type))
		return nil
	}
}

func (p *Parser) Category() *ast.Category {
	c := &ast.Category{Name: p.curr.Lexeme, Phonemes: []ast.CategoryElement{}}
	p.getNextToken() // name
	p.getNextToken() // eq
	for !p.currIs(tok.LINE_ENDING) && !p.currIs(tok.EOF) {
		element := p.WeightedCategoryElement()
		if element != nil {
			c.Phonemes = append(c.Phonemes, element)
		}
		p.getNextToken()
	}
	p.getNextToken() // line ending
	return c
}

func (p *Parser) WeightedCategoryElement() ast.CategoryElement {
	ele := p.CategoryElement()
	if p.peekIs(tok.STAR) {
		p.getNextToken() // star
		p.getNextToken() // weight
		weight, err := strconv.Atoi(p.curr.Lexeme)
		if err != nil {
			p.errors = append(p.errors, err)
			return nil
		}
		return &ast.WeightedElement{Element: ele, Weight: weight}
	} else {
		return &ast.WeightedElement{Element: ele, Weight: 1}
	}
}

func (p *Parser) CategoryElement() ast.CategoryElement {
	fmt.Println("CategoryElement", p.curr.String())
	switch p.curr.Type {
	case tok.RAW:
		return p.Phoneme()
	case tok.DOLLAR:
		return p.Reference()
	default:
		p.errors = append(p.errors, errs.ParserUnexpectedToken(p.peek.Type))
		return nil
	}
}

func (p *Parser) Phoneme() ast.CategoryElement {
	return &ast.Phoneme{Value: p.curr.Lexeme}
}

func (p *Parser) Reference() ast.CategoryElement {
	p.getNextToken()
	return &ast.Reference{Name: p.curr.Lexeme}
}
