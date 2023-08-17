package par

import (
	"phono-word-gen/ast"
	"phono-word-gen/lex"
	"phono-word-gen/parts"
	"phono-word-gen/tok"
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
func (p *Parser) currIsAny(tts ...tok.TokenType) bool {
	for _, tt := range tts {
		if tt == p.curr.Type {
			return true
		}
	}
	return false
}
func (p *Parser) peekIs(tt tok.TokenType) bool { return p.peek.Type == tt }
func (p *Parser) peekIsAny(tts ...tok.TokenType) bool {
	for _, tt := range tts {
		if tt == p.peek.Type {
			return true
		}
	}
	return false
}
func (p *Parser) expectPeek(tt tok.TokenType) bool {
	if p.peekIs(tt) {
		p.getNextToken()
		return true
	} else {
		p.errors = append(p.errors, parts.UnexpectedToken(p.peek, p.peek.Type, tt))
		return false
	}
}
func (p *Parser) expectCurr(tt tok.TokenType) bool {
	if p.currIs(tt) {
		p.getNextToken()
		return true
	} else {
		p.errors = append(p.errors, parts.UnexpectedToken(p.curr, p.curr.Type, tt))
		return false
	}
}

func (p *Parser) Directives() (directives []ast.Directive) {
	for p.peek.Type != tok.EOF {
		dir := p.Directive()
		if dir != nil {
			directives = append(directives, dir)
		} else {
			return
		}
		p.getNextToken()
	}
	return
}

func (p *Parser) Directive() ast.Directive {
	// skip duplicate line endings
	for p.currIs(tok.LINE_ENDING) {
		p.getNextToken()
	}
	switch p.curr.Type {
	case tok.RAW:
		return p.Category()
	case tok.SYLLABLE:
		return p.Syllable()
	default:
		p.errors = append(p.errors, parts.UnknownDirective(p.curr.Type))
		return nil
	}
}
