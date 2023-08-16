package par

import (
	"phono-word-gen/ast"
	"strconv"
)

func (p *Parser) Phoneme() *ast.Phoneme {
	return &ast.Phoneme{Value: p.curr.Lexeme}
}

func (p *Parser) Reference() *ast.Reference {
	p.getNextToken()
	return &ast.Reference{Name: p.curr.Lexeme}
}

func (p *Parser) Weight() int {
	weight, err := strconv.Atoi(p.curr.Lexeme)
	if err != nil {
		p.errors = append(p.errors, err)
		return -1
	}
	return weight
}
