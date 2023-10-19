package lexer

import (
	"testing"

	"github.com/thingsme/thingscript/token"
)

type TokenTest struct {
	expectedType    token.TokenType
	expectedLiteral string
}

func testTokens(t *testing.T, input string, tests []TokenTest) {
	l := New(input)

	for i, tt := range tests {
		tok := l.NextToken()
		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q",
				i, tt.expectedType, tok.Type)
		}
		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - literal wrong, expected=%q, got %q",
				i, tt.expectedLiteral, tok.Literal)
		}
	}
}

func TestVarDef(t *testing.T) {
	input := `
	var five = 5;
	var ten = 10;
	var add = func(x, y) {
		x + y;
	};
	var result = add(five, ten);	
	`
	testTokens(t, input, []TokenTest{
		{token.VAR, "var"}, // 0
		{token.IDENT, "five"},
		{token.ASSIGN, "="},
		{token.INT, "5"},
		{token.SEMICOLON, ";"},
		{token.VAR, "var"}, // 5
		{token.IDENT, "ten"},
		{token.ASSIGN, "="},
		{token.INT, "10"},
		{token.SEMICOLON, ";"},
		{token.VAR, "var"}, // 10
		{token.IDENT, "add"},
		{token.ASSIGN, "="},
		{token.FUNC, "func"},
		{token.LPAREN, "("},
		{token.IDENT, "x"},
		{token.COMMA, ","},
		{token.IDENT, "y"},
		{token.RPAREN, ")"},
		{token.LBRACE, "{"},
		{token.IDENT, "x"},
		{token.PLUS, "+"},
		{token.IDENT, "y"},
		{token.SEMICOLON, ";"},
		{token.RBRACE, "}"},
		{token.SEMICOLON, ";"},
		{token.VAR, "var"},
		{token.IDENT, "result"},
		{token.ASSIGN, "="},
		{token.IDENT, "add"},
		{token.LPAREN, "("},
		{token.IDENT, "five"},
		{token.COMMA, ","},
		{token.IDENT, "ten"},
		{token.RPAREN, ")"},
		{token.SEMICOLON, ";"},
	})
}

func TestUnaryOps(t *testing.T) {
	input := `
		!-/5*;
	`
	testTokens(t, input, []TokenTest{
		// !-/*5;
		{token.BANG, "!"},
		{token.MINUS, "-"},
		{token.SLASH, "/"},
		{token.INT, "5"},
		{token.ASTERISK, "*"},
		{token.SEMICOLON, ";"},
	})
}

func TestLogicalOps(t *testing.T) {
	input := `
		10 == 10;
		10 != 9;
		5 < 10 > 5;
		5 <= 10 >= 5;
		3.14 > 3
		3.14 >= 3
		0xAF < 0xff
		0b10 < 0x1f
		0o7 < 010
	`
	testTokens(t, input, []TokenTest{
		// 10 == 10;
		{token.INT, "10"},
		{token.EQ, "=="},
		{token.INT, "10"},
		{token.SEMICOLON, ";"},
		// 10 != 9;
		{token.INT, "10"},
		{token.NOT_EQ, "!="},
		{token.INT, "9"},
		{token.SEMICOLON, ";"},
		// 5 < 10 > 5;
		{token.INT, "5"},
		{token.LT, "<"},
		{token.INT, "10"},
		{token.GT, ">"},
		{token.INT, "5"},
		{token.SEMICOLON, ";"},
		// 5 <= 10 >= 5;
		{token.INT, "5"},
		{token.LTE, "<="},
		{token.INT, "10"},
		{token.GTE, ">="},
		{token.INT, "5"},
		{token.SEMICOLON, ";"},
		// 3.14 > 3
		{token.FLOAT, "3.14"},
		{token.GT, ">"},
		{token.INT, "3"},
		// 	3.14 >= 3
		{token.FLOAT, "3.14"},
		{token.GTE, ">="},
		{token.INT, "3"},
		// 0xAF < 0xff
		{token.INT, "0xAF"},
		{token.LT, "<"},
		{token.INT, "0xff"},
		// 0b10 < 0x1f
		{token.INT, "0b10"},
		{token.LT, "<"},
		{token.INT, "0x1f"},
		// 0o7 < 010
		{token.INT, "07"},
		{token.LT, "<"},
		{token.INT, "010"},
	})
}

func TestIfElse(t *testing.T) {
	input := `
		if (5 < 10) {
			return true;
		} else {
			return false;
		}
	`

	testTokens(t, input, []TokenTest{
		// if (5 < 10) {
		// 	return true;
		{token.IF, "if"},
		{token.LPAREN, "("},
		{token.INT, "5"},
		{token.LT, "<"},
		{token.INT, "10"},
		{token.RPAREN, ")"},
		{token.LBRACE, "{"},
		{token.RETURN, "return"},
		{token.TRUE, "true"},
		{token.SEMICOLON, ";"},
		// } else {
		{token.RBRACE, "}"},
		{token.ELSE, "else"},
		{token.LBRACE, "{"},
		// 	return false;
		{token.RETURN, "return"},
		{token.FALSE, "false"},
		{token.SEMICOLON, ";"},
		// }
		{token.RBRACE, "}"},
	})
}

func TestOpAssign(t *testing.T) {
	input := `
		v += 10
		v -= 10
		v %= 10
		v = v % 10
		v *= 10.2
		v /= 10.3
	`
	testTokens(t, input, []TokenTest{
		// v += 10
		{token.IDENT, "v"},
		{token.ADDASSIGN, "+="},
		{token.INT, "10"},
		// v -= 10
		{token.IDENT, "v"},
		{token.SUBASSIGN, "-="},
		{token.INT, "10"},
		// v %= 10
		{token.IDENT, "v"},
		{token.MODASSIGN, "%="},
		{token.INT, "10"},
		// v = v % 10
		{token.IDENT, "v"},
		{token.ASSIGN, "="},
		{token.IDENT, "v"},
		{token.PERCENT, "%"},
		{token.INT, "10"},
		// v *= 10.2
		{token.IDENT, "v"},
		{token.MULASSIGN, "*="},
		{token.FLOAT, "10.2"},
		// v /= 10.3
		{token.IDENT, "v"},
		{token.DIVASSIGN, "/="},
		{token.FLOAT, "10.3"},
	})
}

func TestWhile(t *testing.T) {
	input := `
		while i < 10 { i -= 1 }
		do { i += 1 } while (i < 10)
	`
	testTokens(t, input, []TokenTest{
		// 	while i < 10 { i -= 1 }
		{token.WHILE, "while"},
		{token.IDENT, "i"},
		{token.LT, "<"},
		{token.INT, "10"},
		{token.LBRACE, "{"},
		{token.IDENT, "i"},
		{token.SUBASSIGN, "-="},
		{token.INT, "1"},
		{token.RBRACE, "}"},
		// do { i += 1 } while (i < 10)
		{token.DO, "do"},
		{token.LBRACE, "{"},
		{token.IDENT, "i"},
		{token.ADDASSIGN, "+="},
		{token.INT, "1"},
		{token.RBRACE, "}"},
		{token.WHILE, "while"},
		{token.LPAREN, "("},
		{token.IDENT, "i"},
		{token.LT, "<"},
		{token.INT, "10"},
		{token.RPAREN, ")"},
	})
}

func TestComments(t *testing.T) {
	input := `
		/*
			block comment
		*/ x = y // line comment
	`
	testTokens(t, input, []TokenTest{
		// comment
		{token.COMMENT, "\n\t\t\tblock comment\n\t\t"},
		{token.IDENT, "x"},
		{token.ASSIGN, "="},
		{token.IDENT, "y"},
		{token.COMMENT, " line comment"},
	})
}

func TestStringLiteral(t *testing.T) {
	input := `
		"foobar"
		"foo bar"
	`
	testTokens(t, input, []TokenTest{
		//"foobar"
		{token.STRING, "foobar"},
		//"foo bar"
		{token.STRING, "foo bar"},
	})
}

func TestArrayLiteral(t *testing.T) {
	input := `
		[1, 2]	
	`
	testTokens(t, input, []TokenTest{
		// [1, 2]
		{token.LBRACKET, "["}, // tests[75]
		{token.INT, "1"},
		{token.COMMA, ","},
		{token.INT, "2"},
		{token.RBRACKET, "]"},
	})
}

func TestMapLiteral(t *testing.T) {
	input := `
		{"foo": "bar"}
	`
	testTokens(t, input, []TokenTest{
		// {"foot": "bar"}
		{token.LBRACE, "{"},
		{token.STRING, "foo"},
		{token.COLON, ":"},
		{token.STRING, "bar"},
		{token.RBRACE, "}"},
	})
}

func TestVarAssign(t *testing.T) {
	input := `
		twenty := 20
		notnil := twenty ?? 30
	`
	tests := []TokenTest{
		// twenty := 20
		{token.IDENT, "twenty"},
		{token.VARASSIGN, ":="},
		{token.INT, "20"},
		// notnil := twenty ?? 30
		{token.IDENT, "notnil"},
		{token.VARASSIGN, ":="},
		{token.IDENT, "twenty"},
		{token.IMMEDIATEIF, "??"},
		{token.INT, "30"},
		// EOF
		{token.EOF, ""},
	}
	testTokens(t, input, tests)
}
