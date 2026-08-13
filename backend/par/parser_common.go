package par

import (
	"phono-word-gen/ast"
	"strconv"
)

func (p *Parser) Phoneme() *ast.Phoneme {
	return &ast.Phoneme{Value: p.curr.Lexeme}
}

func (p *Parser) CategoryReference() *ast.CategoryReference {
	p.getNextToken()
	return &ast.CategoryReference{Name: p.curr.Lexeme}
}

func (p *Parser) ComponentReference() *ast.ComponentReference {
	p.getNextToken()
	return &ast.ComponentReference{Name: p.curr.Lexeme}
}

func (p *Parser) Weight() int {
	weight, err := strconv.Atoi(p.curr.Lexeme)
	if err != nil {
		p.errors = append(p.errors, err)
		return -1
	}
	return weight
}
