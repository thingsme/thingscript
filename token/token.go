package token

import "fmt"

type TokenType string

type Token struct {
	Type     TokenType
	Literal  string
	NoInfix  bool
	Position Position
}

type Position struct {
	Line   int
	Column int
}

func (p Position) String() string {
	return fmt.Sprintf("Ln %d, Col %d", p.Line, p.Column)
}

const (
	ILLEGAL TokenType = "ILLEGAL"
	EOF     TokenType = "EOF"

	// identifier + literal
	IDENT  TokenType = "IDENT"
	INT    TokenType = "INT"
	FLOAT  TokenType = "FLOAT"
	STRING TokenType = "STRING"

	// operator
	ASSIGN   TokenType = "="
	PLUS     TokenType = "+"
	MINUS    TokenType = "-"
	BANG     TokenType = "!"
	ASTERISK TokenType = "*"
	SLASH    TokenType = "/"
	PERCENT  TokenType = "%"
	DOT      TokenType = "."

	IMMEDIATEIF TokenType = "??"
	VARASSIGN   TokenType = ":="
	ADDASSIGN   TokenType = "+="
	SUBASSIGN   TokenType = "-="
	MULASSIGN   TokenType = "*="
	DIVASSIGN   TokenType = "/="
	MODASSIGN   TokenType = "%="

	LT  TokenType = "<"
	LTE TokenType = "<="
	GT  TokenType = ">"
	GTE TokenType = ">="

	EQ     TokenType = "=="
	NOT_EQ TokenType = "!="

	// seperator
	COMMA     TokenType = ","
	COLON     TokenType = ":"
	SEMICOLON TokenType = ";"

	LPAREN   TokenType = "("
	RPAREN   TokenType = ")"
	LBRACE   TokenType = "{"
	RBRACE   TokenType = "}"
	LBRACKET TokenType = "["
	RBRACKET TokenType = "]"

	// reserved keywords
	FUNC   TokenType = "FUNC"
	VAR    TokenType = "VAR"
	TRUE   TokenType = "TRUE"
	FALSE  TokenType = "FALSE"
	IF     TokenType = "IF"
	ELSE   TokenType = "ELSE"
	RETURN TokenType = "RETURN"
	WHILE  TokenType = "WHILE"
	DO     TokenType = "DO"
	BREAK  TokenType = "BREAK"

	// comment
	COMMENT TokenType = "COMMENT"
)

var keywords = map[string]TokenType{
	"func":   FUNC,
	"var":    VAR,
	"true":   TRUE,
	"false":  FALSE,
	"if":     IF,
	"else":   ELSE,
	"return": RETURN,
	"while":  WHILE,
	"do":     DO,
	"break":  BREAK,
	// reserved
	"const":   ILLEGAL,
	"def":     ILLEGAL,
	"let":     ILLEGAL,
	"class":   ILLEGAL,
	"public":  ILLEGAL,
	"private": ILLEGAL,
	"package": ILLEGAL,
	"for":     ILLEGAL,
}

func LookupIdent(ident string) TokenType {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return IDENT
}
