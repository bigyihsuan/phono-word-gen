package lex

import (
	"phono-word-gen/tok"
	"unicode"
)

type Lexer struct {
	src     []rune
	ch      rune
	currIdx int
	peekIdx int
}

func New(src []rune) *Lexer {
	l := &Lexer{src: append(src, '\n'), currIdx: 0, peekIdx: 0}
	l.nextRune()
	return l
}

func (l *Lexer) GetNextToken() tok.Token {
	var token tok.Token
	l.skipSpace()
	switch l.ch {
	case '(':
		token = tok.New(tok.LPAREN, string(l.ch), l.currIdx)
	case ')':
		token = tok.New(tok.RPAREN, string(l.ch), l.currIdx)
	case '[':
		token = tok.New(tok.LBRACKET, string(l.ch), l.currIdx)
	case ']':
		token = tok.New(tok.RBRACKET, string(l.ch), l.currIdx)
	case '{':
		token = tok.New(tok.LBRACE, string(l.ch), l.currIdx)
	case '}':
		token = tok.New(tok.RBRACE, string(l.ch), l.currIdx)
	case ',':
		token = tok.New(tok.COMMA, string(l.ch), l.currIdx)
	case '*':
		token = tok.New(tok.STAR, string(l.ch), l.currIdx)
	case ':':
		token = tok.New(tok.COLON, string(l.ch), l.currIdx)
	case '=':
		token = tok.New(tok.EQ, string(l.ch), l.currIdx)
	case '$':
		token = tok.New(tok.DOLLAR, string(l.ch), l.currIdx)
	case '>':
		token = tok.New(tok.ARROW, string(l.ch), l.currIdx)
	case '^':
		token = tok.New(tok.CARET, string(l.ch), l.currIdx)
	case '\\':
		token = tok.New(tok.BSLASH, string(l.ch), l.currIdx)
	case '@':
		token = tok.New(tok.AT, string(l.ch), l.currIdx)
	case '&':
		token = tok.New(tok.AMPERSAND, string(l.ch), l.currIdx)
	case '!':
		token = tok.New(tok.BANG, string(l.ch), l.currIdx)
	case '/':
		lexeme := string(l.ch)
		if l.peekRune() == '/' {
			l.nextRune()
			lexeme += string(l.ch)
			token = tok.New(tok.DOUBLESLASH, lexeme, l.currIdx)
		} else {
			token = tok.New(tok.SLASH, lexeme, l.currIdx)
		}
	case ';':
		token = tok.New(tok.LINE_ENDING, string(l.ch), l.currIdx)
	case '\n':
		token = tok.New(tok.LINE_ENDING, string(l.ch), l.currIdx)
	case 0:
		token = tok.New(tok.EOF, "", l.currIdx)
	default:
		token.Index = l.currIdx
		if unicode.IsLetter(l.ch) {
			token.Lexeme = l.raw()
			token.Type = tok.IsKeywordOrRaw(token.Lexeme)
			return token
		} else if unicode.IsDigit(l.ch) {
			token.Lexeme = l.number()
			token.Type = tok.NUMBER
			return token
		} else {
			token = tok.New(tok.ILLEGAL, string(l.ch), l.currIdx)
		}
	}
	l.nextRune()
	return token
}

func (l *Lexer) nextRune() {
	if l.peekIdx >= len(l.src) {
		l.ch = 0
	} else {
		l.ch = l.src[l.peekIdx]
	}
	l.currIdx = l.peekIdx
	l.peekIdx++
}
func (l *Lexer) peekRune() rune {
	if l.peekIdx >= len(l.src) {
		return 0
	} else {
		return l.src[l.peekIdx]
	}
}

func (l *Lexer) skipSpace() {
	for isSpace(l.ch) {
		l.nextRune()
	}
}

func (l *Lexer) raw() string {
	startPosition := l.currIdx
	for isLetter(l.ch) {
		l.nextRune()
	}
	return string(l.src[startPosition:l.currIdx])
}

func (l *Lexer) number() string {
	startPosition := l.currIdx
	for unicode.IsDigit(l.ch) {
		l.nextRune()
	}
	return string(l.src[startPosition:l.currIdx])
}

func isSpace(r rune) bool  { return r == ' ' }
func isLetter(r rune) bool { return unicode.IsLetter(r) }
func isDigit(r rune) bool  { return unicode.IsDigit(r) }
