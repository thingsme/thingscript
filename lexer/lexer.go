package lexer

import (
	"unicode"

	"github.com/thingsme/thingscript/token"
)

type Lexer struct {
	input        []rune
	position     int
	readPosition int
	ch           rune
	tabSize      int

	prevToken token.Token
	Position  token.Position
}

func New(input string) *Lexer {
	l := &Lexer{
		input:    []rune(input),
		tabSize:  4,
		Position: token.Position{Line: 1, Column: 0},
	}
	l.readChar()
	return l
}

func (l *Lexer) NextToken() token.Token {
	var tok token.Token
	l.skipWhitespace()
	position := l.Position
	doReadNext := true
	switch l.ch {
	case '=':
		if l.peekChar() == '=' {
			l.readChar()
			tok = token.Token{Type: token.EQ, Literal: "=="}
		} else {
			tok = newToken(token.ASSIGN, l.ch)
		}
	case '!':
		if l.peekChar() == '=' {
			l.readChar()
			tok = token.Token{Type: token.NOT_EQ, Literal: "!="}
		} else {
			tok = newToken(token.BANG, l.ch)
		}
	case ':':
		if l.peekChar() == '=' {
			l.readChar()
			tok = token.Token{Type: token.VARASSIGN, Literal: ":="}
		} else {
			tok = newToken(token.COLON, l.ch)
		}
	case '?':
		if l.peekChar() == '?' {
			l.readChar()
			tok = token.Token{Type: token.IMMEDIATEIF, Literal: "??"}
		} else {
			tok = newToken(token.ILLEGAL, l.ch)
		}
	case ';':
		tok = newToken(token.SEMICOLON, l.ch)
	case ',':
		tok = newToken(token.COMMA, l.ch)
	case '.':
		tok = newToken(token.DOT, l.ch)
	case '+':
		if l.peekChar() == '=' { // +=
			l.readChar()
			tok = token.Token{Type: token.ADDASSIGN, Literal: "+="}
		} else {
			tok = newToken(token.PLUS, l.ch)
		}
	case '-':
		if l.peekChar() == '=' { // -=
			l.readChar()
			tok = token.Token{Type: token.SUBASSIGN, Literal: "-="}
		} else {
			tok = newToken(token.MINUS, l.ch)
		}
	case '*':
		if l.peekChar() == '=' {
			l.readChar()
			tok = token.Token{Type: token.MULASSIGN, Literal: "*="}
		} else {
			tok = newToken(token.ASTERISK, l.ch)
		}
	case '/':
		if l.peekChar() == '/' { // '//'
			l.readChar()
			comment := l.skipLineComment()
			tok = token.Token{Type: token.COMMENT, Literal: comment}
		} else if l.peekChar() == '*' { // '/*'
			l.readChar()
			comment := l.skipBlockComment()
			tok = token.Token{Type: token.COMMENT, Literal: comment}
		} else if l.peekChar() == '=' { // '/='
			l.readChar()
			tok = token.Token{Type: token.DIVASSIGN, Literal: "/="}
		} else {
			tok = newToken(token.SLASH, l.ch)
		}
	case '%':
		if l.peekChar() == '=' { // '%='
			l.readChar()
			tok = token.Token{Type: token.MODASSIGN, Literal: "%="}
		} else {
			tok = newToken(token.PERCENT, l.ch)
		}
	case '<':
		if l.peekChar() == '=' { //  '<='
			l.readChar()
			tok = token.Token{Type: token.LTE, Literal: "<="}
		} else {
			tok = newToken(token.LT, l.ch)
		}
	case '>':
		if l.peekChar() == '=' { //  '>='
			l.readChar()
			tok = token.Token{Type: token.GTE, Literal: ">="}
		} else {
			tok = newToken(token.GT, l.ch)
		}
	case '(':
		tok = newToken(token.LPAREN, l.ch)
	case ')':
		tok = newToken(token.RPAREN, l.ch)
	case '{':
		tok = newToken(token.LBRACE, l.ch)
	case '}':
		tok = newToken(token.RBRACE, l.ch)
	case '[':
		tok = newToken(token.LBRACKET, l.ch)
	case ']':
		tok = newToken(token.RBRACKET, l.ch)
	case '"':
		tok.Literal = l.readString()
		tok.Type = token.STRING
	case 0:
		tok.Literal = ""
		tok.Type = token.EOF
	default:
		if isLetter(l.ch) {
			tok.Literal = l.readIdentifier()
			tok.Type = token.LookupIdent(tok.Literal)
			doReadNext = false
		} else if isDigit(l.ch) {
			tok.Literal, tok.Type = l.readNumber()
			doReadNext = false
		} else {
			tok = newToken(token.ILLEGAL, l.ch)
		}
	}
	tok.Position = position
	if l.prevToken.Position.Line != tok.Position.Line {
		// 'tok' is the first of the current line
		// which can not be evaluated as an infix operator
		tok.NoInfix = true
	}
	if doReadNext {
		l.readChar()
	}
	l.prevToken = tok
	return tok
}

func newToken(tokenType token.TokenType, ch rune) token.Token {
	return token.Token{Type: tokenType, Literal: string(ch)}
}

func (l *Lexer) readIdentifier() string {
	position := l.position
	for isLetter(l.ch) {
		l.readChar()
	}
	return string(l.input[position:l.position])
}

func (l *Lexer) readNumber() (string, token.TokenType) {
	if l.ch == '0' && l.peekChar() == 'x' { // base 16 digit: prefix '0x'
		l.readChar()
		l.readChar()
		position := l.position
		for unicode.Is(unicode.Hex_Digit, l.ch) {
			l.readChar()
		}
		return "0x" + string(l.input[position:l.position]), token.INT
	} else if l.ch == '0' && (l.peekChar() == 'o' || isDigit(l.peekChar())) { // base 8 : prefix '0' or '0o'
		l.readChar()
		if l.ch == 'o' {
			l.readChar()
		}
		position := l.position
		for isDigit(l.ch) {
			l.readChar()
		}
		return "0" + string(l.input[position:l.position]), token.INT
	} else if l.ch == '0' && l.peekChar() == 'b' { // base 2 : prefix '0b'
		l.readChar()
		l.readChar()
		position := l.position
		for isDigit(l.ch) {
			l.readChar()
		}
		return "0b" + string(l.input[position:l.position]), token.INT
	} else { // base 10 digit
		position := l.position
		numType := token.INT
		for isDigit(l.ch) || l.ch == '.' {
			if l.ch == '.' {
				if numType == token.FLOAT {
					return "", token.ILLEGAL
				} else {
					numType = token.FLOAT
				}
			}
			l.readChar()
		}
		return string(l.input[position:l.position]), numType
	}
}

func (l *Lexer) readString() string {
	position := l.position + 1
	for {
		l.readChar()
		if l.ch == '"' || l.ch == 0 {
			break
		}
	}
	return string(l.input[position:l.position])
}

func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.ch = 0
	} else {
		l.ch = l.input[l.readPosition]
	}
	l.position = l.readPosition
	l.readPosition += 1
	if l.ch == '\n' {
		l.Position.Line++
		l.Position.Column = 0
	} else {
		if l.ch == '\t' {
			l.Position.Column += l.tabSize
		} else if l.ch == '\r' {
			// ignore
		} else {
			l.Position.Column++
		}
	}
}

func (l *Lexer) peekChar() rune {
	if l.readPosition >= len(l.input) {
		return 0
	} else {
		return l.input[l.readPosition]
	}
}

func (l *Lexer) skipWhitespace() {
	for unicode.IsSpace(l.ch) {
		l.readChar()
	}
}

func (l *Lexer) skipLineComment() string {
	position := l.position + 1
	for l.ch != 0 && l.ch != '\n' {
		l.readChar()
	}
	if position < l.position {
		return string(l.input[position:l.position])
	} else {
		return ""
	}
}

func (l *Lexer) skipBlockComment() string {
	position := l.position + 1
	for l.ch != 0 {
		if l.ch == '*' && l.peekChar() == '/' {
			l.readChar()
			return string(l.input[position : l.position-1])
		}
		l.readChar()
	}
	return ""
}

func isLetter(ch rune) bool {
	return unicode.IsLetter(ch) || ch == '_'
}

func isDigit(ch rune) bool {
	return unicode.IsDigit(ch)
}
