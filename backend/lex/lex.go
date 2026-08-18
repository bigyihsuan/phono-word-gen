package lex

import (
	"phono-word-gen/tok"
	"strconv"
	"strings"
	"unicode"
)

const symbols = "\n;*#$=[](){}|,/_:>^@!\\&%"

type Lexer struct {
	src     []rune
	ch      rune // current character
	currIdx int
	peekIdx int
	// for Pos
	currLine int // 1-indexed line
	currCol  int // 1-indexed col on the line
}

func New(src []rune) *Lexer {
	l := &Lexer{src: append(src, '\n'), currIdx: 0, peekIdx: 0, currLine: 0, currCol: 0}
	// init currIdx and peekIdx for peekRune
	l.nextRune()
	l.advanceLine()
	return l
}

func (l *Lexer) GetNextToken() tok.Token {
	var token tok.Token
	l.skipSpace()
	switch l.ch {
	case '#': // comment
		lexeme := string(l.ch)
		for l.peekRune() != '\n' && l.peekRune() != ';' {
			l.nextRune()
			lexeme += string(l.ch)
		}
		token = tok.New(tok.COMMENT, lexeme, l.currPos())
	case '(':
		token = tok.New(tok.LPAREN, string(l.ch), l.currPos())
	case ')':
		token = tok.New(tok.RPAREN, string(l.ch), l.currPos())
	case '[':
		token = tok.New(tok.LBRACKET, string(l.ch), l.currPos())
	case ']':
		token = tok.New(tok.RBRACKET, string(l.ch), l.currPos())
	case '{':
		token = tok.New(tok.LBRACE, string(l.ch), l.currPos())
	case '}':
		token = tok.New(tok.RBRACE, string(l.ch), l.currPos())
	case ',':
		token = tok.New(tok.COMMA, string(l.ch), l.currPos())
	case '*':
		token = tok.New(tok.STAR, string(l.ch), l.currPos())
	case ':':
		token = tok.New(tok.COLON, string(l.ch), l.currPos())
	case '=':
		token = tok.New(tok.EQ, string(l.ch), l.currPos())
	case '$':
		token = tok.New(tok.DOLLAR, string(l.ch), l.currPos())
	case '%':
		token = tok.New(tok.PERCENT, string(l.ch), l.currPos())
	case '>':
		token = tok.New(tok.ARROW, string(l.ch), l.currPos())
	case '^':
		token = tok.New(tok.CARET, string(l.ch), l.currPos())
	case '\\':
		token = tok.New(tok.BSLASH, string(l.ch), l.currPos())
	case '@':
		token = tok.New(tok.AT, string(l.ch), l.currPos())
	case '&':
		token = tok.New(tok.AMPERSAND, string(l.ch), l.currPos())
	case '|':
		token = tok.New(tok.PIPE, string(l.ch), l.currPos())
	case '_':
		token = tok.New(tok.UNDERSCORE, string(l.ch), l.currPos())
	case '!':
		token = tok.New(tok.BANG, string(l.ch), l.currPos())
	case '/':
		lexeme := string(l.ch)
		if l.peekRune() == '/' {
			l.nextRune()
			lexeme += string(l.ch)
			token = tok.New(tok.DOUBLESLASH, lexeme, l.currPos())
		} else {
			token = tok.New(tok.SLASH, lexeme, l.currPos())
		}
	case ';', '\n':
		token = tok.New(tok.LINE_ENDING, string(l.ch), l.currPos())
		if l.ch == '\n' { // newlines advance the line counter, semicolons do not
			l.advanceLine()
		}
	case 0:
		token = tok.New(tok.EOF, "", l.currPos())
	default:
		token.Pos.Index = l.currIdx
		token.Lexeme = l.raw()
		if tok.IsKeyword(token.Lexeme) {
			token.Type = tok.Keyword(token.Lexeme)
		} else if _, err := strconv.ParseInt(token.Lexeme, 0, 64); err == nil {
			token.Type = tok.NUMBER
		} else {
			token.Type = tok.RAW
		}
		return token
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
	l.currCol++
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

func (l *Lexer) keywordOrRaw() string {
	startPosition := l.currIdx
	for !isRawEnder(l.ch) {
		l.nextRune()
		s := string(l.src[startPosition:l.currIdx])
		tt := tok.Keyword(string(l.src[startPosition:l.currIdx]))
		if tt != tok.RAW {
			return s
		}
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

func (l *Lexer) raw() string {
	startPosition := l.currIdx
	for !isRawEnder(l.ch) {
		l.nextRune()
	}
	return string(l.src[startPosition:l.currIdx])
}

func (l *Lexer) advanceLine() {
	l.currLine++
	l.currCol = 1
}

func (l Lexer) currPos() tok.Pos {
	return tok.Pos{
		Index: l.currIdx,
		Line:  l.currLine,
		Col:   l.currCol,
	}
}

func isSpace(r rune) bool  { return r == ' ' }
func isSymbol(r rune) bool { return strings.ContainsRune(symbols, r) }
func isLetter(r rune) bool { return unicode.IsLetter(r) }
func isDigit(r rune) bool  { return unicode.IsDigit(r) }

func isRawEnder(r rune) bool { return isSpace(r) || isSymbol(r) }
