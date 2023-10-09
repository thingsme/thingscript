package parser

import (
	"fmt"
	"strconv"
	"strings"
	"testing"

	"github.com/thingsme/thingscript/ast"
	"github.com/thingsme/thingscript/lexer"
)

func TestReturnStatements(t *testing.T) {
	tests := []struct {
		input         string
		expectedValue any
	}{
		{`return 5;`, 5},
		{`return true;`, true},
	}
	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)

		program := p.ParseProgram()
		checkParseErrors(t, p)
		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements does not contains 1 statements. got=%d",
				len(program.Statements))
		}
		for _, stmt := range program.Statements {
			returnStmt, ok := stmt.(*ast.ReturnStatement)
			if !ok {
				t.Errorf("stmt not *ast.ReturnStatement. got=%T", stmt)
				continue
			}
			if returnStmt.TokenLiteral() != "return" {
				t.Errorf("returnStmt.TokenLiteral not 'return', got %q",
					returnStmt.TokenLiteral())
			}
			if !testLiteralExpression(t, returnStmt.ReturnValue, tt.expectedValue) {
				return
			}
		}
	}
}

func TestVarStatements(t *testing.T) {
	tests := []struct {
		input              string
		expectedIdentifier string
		expectedValue      any
	}{
		{`var x = 5;`, "x", 5},
		{`x := 5;`, "x", 5},
		{`var z = 0.0`, "z", 0.0},
		{`var f = 3.14`, "f", 3.14},
		{`var h = 0x1a`, "h", 26},
		{`var o = 017`, "o", 15},
		{`var b = 0b11`, "b", 3},
		{`var y = true;`, "y", true},
		{`var y = /*false*/ true; //comment`, "y", true},
		{`y := true;`, "y", true},
		{`var foobar = y;`, "foobar", "y"},
		{`foobar := z;`, "foobar", "z"},
	}
	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)

		program := p.ParseProgram()
		checkParseErrors(t, p)
		if program == nil {
			t.Fatalf("ParseProgram() returned nil")
		}
		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements does not contain 1 statements. got=%d", len(program.Statements))
		}
		stmt := program.Statements[0]
		if !testVarStatement(t, stmt, tt.expectedIdentifier) {
			return
		}
		val := stmt.(*ast.VarStatement).Value
		if !testLiteralExpression(t, val, tt.expectedValue) {
			return
		}
	}
}

func checkParseErrors(t *testing.T, p *Parser) {
	t.Helper()
	errors := p.Errors()
	if len(errors) == 0 {
		return
	}
	t.Errorf("parser has %d errors", len(errors))
	for _, msg := range errors {
		t.Errorf("parser error: %q", msg)
	}
	t.FailNow()
}

func testVarStatement(t *testing.T, s ast.Statement, name string) bool {
	t.Helper()
	if s.TokenLiteral() != "var" {
		t.Errorf("s.TokenLiteral not 'var'. got=%q", s.TokenLiteral())
		return false
	}
	letStmt, ok := s.(*ast.VarStatement)
	if !ok {
		t.Errorf("s not *ast.LetStatement, got=%T", s)
		return false
	}
	if letStmt.Name.Value != name {
		t.Errorf("letStmt.Name.Value not '%s'. got=%s", name, letStmt.Name.Value)
		return false
	}
	if letStmt.Name.TokenLiteral() != name {
		t.Errorf("letStmt.Name.TokenLiteral() not '%s'. got=%s", name, letStmt.Name.TokenLiteral())
		return false
	}
	return true
}

func TestIdentifierExpression(t *testing.T) {
	input := "foobar;"
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParseErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program has not enough statements, got=%d",
			len(program.Statements))
	}
	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement, got=%T",
			program.Statements[0])
	}
	ident, ok := stmt.Expression.(*ast.Identifier)
	if !ok {
		t.Fatalf("exp is not *ast.Identifier, got=%T", stmt.Expression)
	}
	if ident.Value != "foobar" {
		t.Errorf("ident.Value not %s, got=%T", "foobar", ident.Value)
	}
	if ident.TokenLiteral() != "foobar" {
		t.Errorf("ident.TokenLiteral not %s, got=%s", "foobar", ident.TokenLiteral())
	}
}

func TestAssignStatement(t *testing.T) {
	input := `foo = "bar"`
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParseErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program has not enough statements, got=%d",
			len(program.Statements))
	}
	stmt, ok := program.Statements[0].(*ast.AssignStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.AssignStatement, got=%T",
			program.Statements[0])
	}
	ident := stmt.Name
	if ident.Value != "foo" {
		t.Errorf("ident.Value not %s, got=%T", "foobar", ident.Value)
	}
	value, ok := stmt.Value.(*ast.StringLiteral)
	if !ok {
		t.Errorf("right value not ast.StringLiteral, got=%T", value)
	}
	if value.TokenLiteral() != "bar" {
		t.Errorf("right value not %s, got=%s", "bar", value.TokenLiteral())
	}
}

func TestOperAssignStatement(t *testing.T) {
	input := `foo += 123`
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParseErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program has not enough statements, got=%d",
			len(program.Statements))
	}
	stmt, ok := program.Statements[0].(*ast.OperAssignStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.OperAssignStatement, got=%T",
			program.Statements[0])
	}
	if stmt.Operator != "+" {
		t.Errorf("operator not %q, got=%q", "+", stmt.Operator)
	}
	ident := stmt.Name
	if ident.Value != "foo" {
		t.Errorf("ident.Value not %s, got=%T", "foobar", ident.Value)
	}
	value, ok := stmt.Value.(*ast.IntegerLiteral)
	if !ok {
		t.Errorf("right value not ast.IntegerLiteral, got=%T", value)
	}
	if value.TokenLiteral() != "123" {
		t.Errorf("right value not %s, got=%s", "123", value.TokenLiteral())
	}
}

func TestImmediateIfExpression(t *testing.T) {
	input := "foo ?? bar"
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParseErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program has not enough statements, got=%d",
			len(program.Statements))
	}
	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement, got=%T",
			program.Statements[0])
	}

	exp, ok := stmt.Expression.(*ast.ImmediateIfExpression)
	if !ok {
		t.Fatalf("exp is not *ast.ImmediateIfExpression, got=%T", stmt.Expression)
	}
	if exp.Left.TokenLiteral() != "foo" {
		t.Errorf("exp.Left not %s, got=%s", "foo", exp.Left.TokenLiteral())
	}
	if exp.Right.TokenLiteral() != "bar" {
		t.Errorf("exp.Right not %s, got=%s", "bar", exp.Right.TokenLiteral())
	}
}

func TestIntegerLiteralExpression(t *testing.T) {
	input := "5;"
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParseErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program has not enough statements, got=%d",
			len(program.Statements))
	}
	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement, got=%T",
			program.Statements[0])
	}
	testIntegerLiteral(t, stmt.Expression, 5)
}

func TestBooleanExpression(t *testing.T) {
	input := "true;"
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParseErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program has not enough statements, got=%d",
			len(program.Statements))
	}
	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement, got=%T",
			program.Statements[0])
	}
	testBooleanLiteral(t, stmt.Expression, true)
}

func TestStringLiteralExpression(t *testing.T) {
	input := `"hello world"`
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParseErrors(t, p)
	stmt := program.Statements[0].(*ast.ExpressionStatement)
	literal, ok := stmt.Expression.(*ast.StringLiteral)
	if !ok {
		t.Fatalf("exp not *ast.StringLiteral, got=%T", stmt.Expression)
	}
	if literal.Value != "hello world" {
		t.Errorf("literal.Value not %q. got=%q", "hello world", literal.Value)
	}
}

func TestParsingAccessMembers(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{`[1, 2, 3].field`, `(([1, 2, 3]).(field))`},
		{`[1, 2, 3].call()`, `(([1, 2, 3]).(call()))`},
		{`[1, 2, 3].call(true)`, `(([1, 2, 3]).(call(true)))`},
	}
	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParseErrors(t, p)
		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("statement not ast.ExpressionStatement. got=%T", program.Statements[0])
		}
		access, ok := stmt.Expression.(*ast.AccessExpression)
		if !ok {
			t.Fatalf("expression not ast.AccessExpression, got=%T %s", stmt.Expression, stmt.Expression.String())
		}
		if tt.expected != access.String() {
			t.Errorf("access expr not %q, got=%q", tt.expected, access.String())
		}
	}
}

func TestParsingArrayLiterals(t *testing.T) {
	input := "[1, 2 * 2, 3 + 3]"
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParseErrors(t, p)

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("statement not ast.ExpressionStatement. got=%T", program.Statements[0])
	}
	array, ok := stmt.Expression.(*ast.ArrayLiteral)
	if !ok {
		t.Fatalf("expression not ast.ArrayLiteral. got=%T", stmt.Expression)
	}
	if len(array.Elements) != 3 {
		t.Fatalf("len(array.Elements) not 3, got=%d", len(array.Elements))
	}
	testIntegerLiteral(t, array.Elements[0], 1)
	testInfixExpression(t, array.Elements[1], 2, "*", 2)
	testInfixExpression(t, array.Elements[2], 3, "+", 3)
}

func TestParsingHashLiteral(t *testing.T) {
	input := `{"one": 1, "two": 2, "three": 3}`
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParseErrors(t, p)

	stmt := program.Statements[0].(*ast.ExpressionStatement)
	hash, ok := stmt.Expression.(*ast.HashMapLiteral)
	if !ok {
		t.Fatalf("exp is not ast.HashLiteral. got=%T", stmt.Expression)
	}
	if len(hash.Pairs) != 3 {
		t.Errorf("hash.Pairs has wrong length. got=%d", len(hash.Pairs))
	}
	expected := map[string]int64{
		"one":   1,
		"two":   2,
		"three": 3,
	}
	for key, value := range hash.Pairs {
		literal, ok := key.(*ast.StringLiteral)
		if !ok {
			t.Errorf("key is not ast.StringLiteral. got=%T", key)
		}
		expectedValue := expected[literal.Value]
		testIntegerLiteral(t, value, expectedValue)
	}
}

func TestParsingEmptyHashLiteral(t *testing.T) {
	input := `{}`
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParseErrors(t, p)

	stmt := program.Statements[0].(*ast.ExpressionStatement)
	hash, ok := stmt.Expression.(*ast.HashMapLiteral)
	if !ok {
		t.Fatalf("exp is not ast.HashLiteral. got=%T", stmt.Expression)
	}
	if len(hash.Pairs) != 0 {
		t.Errorf("hash.Pairs has wrong length. got=%d", len(hash.Pairs))
	}
}

func TestParsingIndexExpressions(t *testing.T) {
	input := "myArray[1 + 1]"
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParseErrors(t, p)

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("statement not ast.ExpressionStatement. got=%T", program.Statements[0])
	}
	indexExp, ok := stmt.Expression.(*ast.IndexExpression)
	if !ok {
		t.Fatalf("expression not ast.IndexExpression. got=%T", stmt.Expression)
	}
	if !testIdentifier(t, indexExp.Left, "myArray") {
		return
	}
	if !testInfixExpression(t, indexExp.Index, 1, "+", 1) {
		return
	}
}

func TestParsingPrefixExpressions(t *testing.T) {
	prefixTests := []struct {
		input        string
		operator     string
		integerValue any
	}{
		{`!5`, "!", 5},
		{`-15`, "-", 15},
		{`!true;`, "!", true},
		{`!false;`, "!", false},
	}

	for _, tt := range prefixTests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParseErrors(t, p)

		if len(program.Statements) != 1 {
			t.Fatalf("program has not contain %d statements, got=%d",
				1, len(program.Statements))
		}
		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("program.Statements[0] is not ast.ExpressionStatement, got=%T",
				program.Statements[0])
		}
		exp, ok := stmt.Expression.(*ast.PrefixExpression)
		if !ok {
			t.Fatalf("stmt is not ast.PrefixExpression, got=%T", stmt.Expression)
		}
		if exp.Operator != tt.operator {
			t.Fatalf("exp.Operator is not '%s', got=%s", tt.operator, exp.Operator)
		}
		if !testLiteralExpression(t, exp.Right, tt.integerValue) {
			return
		}
	}
}

func TestParsingInfixExpression(t *testing.T) {
	infixTests := []struct {
		input      string
		leftValue  any
		operator   string
		rightValue any
	}{
		{`5 + 5`, 5, "+", 5},
		{`5 - 5`, 5, "-", 5},
		{`5 * 5`, 5, "*", 5},
		{`5 / 5`, 5, "/", 5},
		{`5 > 5`, 5, ">", 5},
		{`5 < 5`, 5, "<", 5},
		{`5 >= 5`, 5, ">=", 5},
		{`5 <= 5`, 5, "<=", 5},
		{`5 == 5`, 5, "==", 5},
		{`5 != 5`, 5, "!=", 5},
		{`true == true`, true, "==", true},
		{`true != false`, true, "!=", false},
		{`false == false`, false, "==", false},
	}

	for _, tt := range infixTests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParseErrors(t, p)

		if len(program.Statements) != 1 {
			t.Fatalf("program has not contain %d statements, got=%d",
				1, len(program.Statements))
		}
		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("program.Statements[0] is not ast.ExpressionStatement, got=%T",
				program.Statements[0])
		}

		if !testInfixExpression(t, stmt.Expression, tt.leftValue, tt.operator, tt.rightValue) {
			return
		}
	}
}

func testFloatLiteral(t *testing.T, fl ast.Expression, value float64) bool {
	float, ok := fl.(*ast.FloatLiteral)
	if !ok {
		t.Fatalf("fl is not *ast.FloatLiteral, got=%T", fl)
		return false
	}
	if float.Value != value {
		t.Errorf("float.Value not %f, got=%T", value, float.Value)
		return false
	}
	evaled, err := strconv.ParseFloat(float.TokenLiteral(), 64)
	if err != nil {
		t.Errorf("float wrong format; %s", err)
		return false
	}
	if evaled != value {
		t.Errorf("literal.TokenLiteral not %f, got=%s", value, float.TokenLiteral())
		return false
	}
	return true
}
func testIntegerLiteral(t *testing.T, il ast.Expression, value int64) bool {
	integ, ok := il.(*ast.IntegerLiteral)
	if !ok {
		t.Fatalf("il is not *ast.IntegerLiteral, got=%T", il)
		return false
	}
	if integ.Value != value {
		t.Errorf("integ.Value not %d, got=%T", value, integ.Value)
		return false
	}
	literal := integ.TokenLiteral()
	strValue := fmt.Sprintf("%d", value)
	if strings.HasPrefix(literal, "0x") {
		strValue = fmt.Sprintf("0x%x", value)
	} else if strings.HasPrefix(literal, "0b") {
		strValue = fmt.Sprintf("0b%b", value)
	} else if strings.HasPrefix(literal, "0") {
		strValue = fmt.Sprintf("0%o", value)
	}
	if integ.TokenLiteral() != strValue {
		t.Errorf("literal.TokenLiteral not %d, got=%s", value, integ.TokenLiteral())
		return false
	}
	return true
}

func testBooleanLiteral(t *testing.T, exp ast.Expression, value bool) bool {
	bo, ok := exp.(*ast.Boolean)
	if !ok {
		t.Fatalf("il is not *ast.Boolean, got=%T", bo)
		return false
	}
	if bo.Value != value {
		t.Errorf("bo.Value not %t, got=%T", value, bo.Value)
		return false
	}
	if bo.TokenLiteral() != fmt.Sprintf("%t", value) {
		t.Errorf("literal.TokenLiteral not %t, got=%s", value, bo.TokenLiteral())
		return false
	}
	return true
}

func testIdentifier(t *testing.T, exp ast.Expression, value string) bool {
	t.Helper()
	ident, ok := exp.(*ast.Identifier)
	if !ok {
		t.Errorf("exp not *ast.Identifier, got=%T", exp)
		return false
	}
	if ident.Value != value {
		t.Errorf("ident.Value not %s, got=%s", value, ident.Value)
		return false
	}
	if ident.TokenLiteral() != value {
		t.Errorf("ident.TokenLiteral not %s. got=%s", value, ident.TokenLiteral())
		return false
	}
	return true
}

func testLiteralExpression(t *testing.T, exp ast.Expression, expected any) bool {
	t.Helper()
	switch v := expected.(type) {
	case int:
		return testIntegerLiteral(t, exp, int64(v))
	case int64:
		return testIntegerLiteral(t, exp, v)
	case float32:
		return testFloatLiteral(t, exp, float64(v))
	case float64:
		return testFloatLiteral(t, exp, v)
	case string:
		return testIdentifier(t, exp, v)
	case bool:
		return testBooleanLiteral(t, exp, v)
	}
	t.Errorf("type of exp not handled. got=%T", exp)
	return false
}

func testInfixExpression(t *testing.T, exp ast.Expression, left any, operator string, right any) bool {
	t.Helper()
	opExp, ok := exp.(*ast.InfixExpression)
	if !ok {
		t.Errorf("exp is not ast.InfixExpression, got=%T(%s)", exp, exp)
		return false
	}
	if !testLiteralExpression(t, opExp.Left, left) {
		return false
	}
	if opExp.Operator != operator {
		t.Errorf("exp.Operator is not %s, got=%q", operator, opExp.Operator)
		return false
	}
	if !testLiteralExpression(t, opExp.Right, right) {
		return false
	}
	return true
}

func TestOperatorPrecedenceParsing(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{`-a *b`, `((-a) * b)`},
		{`!-a`, `(!(-a))`},
		{`a + b +c`, `((a + b) + c)`},
		{`a + b -c`, `((a + b) - c)`},
		{`a * b * c`, `((a * b) * c)`},
		{`a * b / c`, `((a * b) / c)`},
		{`a + b / c`, `(a + (b / c))`},
		{`a + b * c + d /e -f`, `(((a + (b * c)) + (d / e)) - f)`},
		{`3 + 4; -5*5`, `(3 + 4)((-5) * 5)`},
		{`5 > 4 == 3 < 4`, `((5 > 4) == (3 < 4))`},
		{`3 + 4 * 5 == 3 * 1 + 4 *5`, `((3 + (4 * 5)) == ((3 * 1) + (4 * 5)))`},
		{`true`, `true`},
		{`false`, `false`},
		{`3 > 5 == false`, `((3 > 5) == false)`},
		{`3 < 5 == true`, `((3 < 5) == true)`},
		{`1 + (2 + 3) + 4`, `((1 + (2 + 3)) + 4)`},
		{`(5 + 5) * 2`, `((5 + 5) * 2)`},
		{`2 / ( 5 + 5)`, `(2 / (5 + 5))`},
		{`-(5+5)`, `(-(5 + 5))`},
		{`!(true == true)`, `(!(true == true))`},
		{`a + add(b *c) + d`, `((a + add((b * c))) + d)`},
		{`add(a, b, 1, 2*3, 4 + 5, add(6, 7 * 8))`, `add(a, b, 1, (2 * 3), (4 + 5), add(6, (7 * 8)))`},
		{`add(a +b + c * d /f + g)`, `add((((a + b) + ((c * d) / f)) + g))`},
		{`a * [1, 2, 3, 4][b*c] *d`, `((a * ([1, 2, 3, 4][(b * c)])) * d)`},
		{`add(a * b[2], b[1], 2 * [1,2][1])`, `add((a * (b[2])), (b[1]), (2 * ([1, 2][1])))`},
		{`"hello".len()`, `((hello).(len()))`},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParseErrors(t, p)

		actual := program.String()
		if actual != tt.expected {
			t.Errorf("expected=%q, got=%q", tt.expected, actual)
		}
	}
}

func TestIfExpression(t *testing.T) {
	input := `if (x < y) { x }`
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParseErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain %d statements. got=%d",
			1, len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement, got=%T",
			program.Statements[0])
	}
	exp, ok := stmt.Expression.(*ast.IfExpression)
	if !ok {
		t.Fatalf("stmt.Expression is not ast.IfStatement, got=%T", stmt.Expression)
	}
	if !testInfixExpression(t, exp.Condition, "x", "<", "y") {
		return
	}
	if len(exp.Consequence.Statements) != 1 {
		t.Errorf("consequence is not 1 statements. got=%d", len(exp.Consequence.Statements))
	}
	consequence, ok := exp.Consequence.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("Statements[0] is not ast.ExpressionStatement, got=%T",
			exp.Consequence.Statements[0])
	}
	if !testIdentifier(t, consequence.Expression, "x") {
		return
	}
	if exp.Alternative != nil {
		t.Errorf("exp.Alternative.Statements was not nil, got=%+v", exp.Alternative)
	}
}

func TestIfElseExpression(t *testing.T) {
	input := `if (x < y) { x } else { y }`
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParseErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain %d statements. got=%d",
			1, len(program.Statements))
	}
	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement, got=%T",
			program.Statements[0])
	}
	exp, ok := stmt.Expression.(*ast.IfExpression)
	if !ok {
		t.Fatalf("stmt.Expression is not ast.IfStatement, got=%T", stmt.Expression)
	}
	if !testInfixExpression(t, exp.Condition, "x", "<", "y") {
		return
	}
	if len(exp.Consequence.Statements) != 1 {
		t.Errorf("consequence is not 1 statements. got=%d", len(exp.Consequence.Statements))
	}
	consequence, ok := exp.Consequence.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("Statements[0] is not ast.ExpressionStatement, got=%T",
			exp.Consequence.Statements[0])
	}
	if !testIdentifier(t, consequence.Expression, "x") {
		return
	}
	if exp.Alternative == nil {
		t.Errorf("exp.Alternative.Statements was nil")
	}
	if len(exp.Alternative.Statements) != 1 {
		t.Errorf("exp.Alternative.Statements is not 1 statement, got=%d", len(exp.Alternative.Statements))
	}
	alternative, ok := exp.Alternative.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("Alternative.Statements[0] is not ast.ExpressionStatement, got=%T",
			exp.Alternative.Statements[0])
	}
	if !testIdentifier(t, alternative.Expression, "y") {
		return
	}
}

func TestWhileExpression(t *testing.T) {
	input := `while (x < 5) { break }`
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParseErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain %d statements. got=%d",
			1, len(program.Statements))
	}
	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement, got=%T",
			program.Statements[0])
	}
	exp, ok := stmt.Expression.(*ast.WhileExpression)
	if !ok {
		t.Fatalf("stmt.Expression is not ast.WhileStatement, got=%T", stmt.Expression)
	}
	if exp.Condition.String() != "(x < 5)" {
		t.Errorf("condition is not %s. got=%s", "x < 5", exp.Condition.String())
		return
	}
	if len(exp.Block.Statements) != 1 {
		t.Errorf("block is not 1 statements. got=%d", len(exp.Block.Statements))
	}
	breakStatement, ok := exp.Block.Statements[0].(*ast.BreakStatement)
	if !ok {
		t.Fatalf("Statements[0] is not ast.BreakStatement, got=%T",
			exp.Block.Statements[0])
	}
	if breakStatement.TokenLiteral() != "break" {
		t.Errorf("block is not break statement. got=%s", breakStatement.String())
	}
}

func TestDoWhileExpression(t *testing.T) {
	input := `do { break } while (x < 5)`
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParseErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain %d statements. got=%d",
			1, len(program.Statements))
	}
	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement, got=%T",
			program.Statements[0])
	}
	exp, ok := stmt.Expression.(*ast.DoWhileExpression)
	if !ok {
		t.Fatalf("stmt.Expression is not ast.DoWhileStatement, got=%T", stmt.Expression)
	}
	if exp.Condition.String() != "(x < 5)" {
		t.Errorf("condition is not %s. got=%s", "x < 5", exp.Condition.String())
		return
	}
	if len(exp.Block.Statements) != 1 {
		t.Errorf("block is not 1 statements. got=%d", len(exp.Block.Statements))
	}
	breakStatement, ok := exp.Block.Statements[0].(*ast.BreakStatement)
	if !ok {
		t.Fatalf("Statements[0] is not ast.BreakStatement, got=%T",
			exp.Block.Statements[0])
	}
	if breakStatement.TokenLiteral() != "break" {
		t.Errorf("block is not break statement. got=%s", breakStatement.String())
	}
}

func TestFunctionLiteralParsing(t *testing.T) {
	input := `func(x, y){ x + y }`
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParseErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain %d statements. got=%d", 1, len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])
	}
	fn, ok := stmt.Expression.(*ast.FunctionLiteral)
	if !ok {
		t.Fatalf("stmt.Expression is not ast.FunctionLiteral. got=%T", stmt.Expression)
	}
	if len(fn.Parameters) != 2 {
		t.Fatalf("function literal paramsters wrong. want 2, got=%d", len(fn.Parameters))
	}
	testLiteralExpression(t, fn.Parameters[0], "x")
	testLiteralExpression(t, fn.Parameters[1], "y")

	if len(fn.Body.Statements) != 1 {
		t.Fatalf("function.Body.Statements has not 1 statements. got=%d", len(fn.Body.Statements))
	}
	bodyStmt, ok := fn.Body.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("function body stmt is not ast.ExpressionStatement. got=%T", fn.Body.Statements[0])
	}
	testInfixExpression(t, bodyStmt.Expression, "x", "+", "y")
}

func TestFunction(t *testing.T) {
	input := `func myFunction(){  }`
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParseErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain %d statements. got=%d", 1, len(program.Statements))
	}

	fn, ok := program.Statements[0].(*ast.FunctionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.FunctionStatement. got=%T", program.Statements[0])
	}
	if fn.Name.Value != "myFunction" {
		t.Fatalf("function literal name wrong. want 'myFunction', got=%q", fn.Name.Value)
	}
}

func TestFunctionLiteralwithName(t *testing.T) {
	input := `var myFunction = func(){  }`
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParseErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain %d statements. got=%d", 1, len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.VarStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])
	}
	fn, ok := stmt.Value.(*ast.FunctionLiteral)
	if !ok {
		t.Fatalf("stmt.Expression is not ast.FunctionLiteral. got=%T", stmt.Value)
	}
	if fn.Name != "myFunction" {
		t.Fatalf("function literal name wrong. want 'myFunction', got=%q", fn.Name)
	}
}

func TestFunctionVoidListeralParsing(t *testing.T) {
	input := `func() { return x + y; }`
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParseErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain %d statements. got=%d", 1, len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])
	}
	fn, ok := stmt.Expression.(*ast.FunctionLiteral)
	if !ok {
		t.Fatalf("stmt.Expression is not ast.FunctionLiteral. got=%T", stmt.Expression)
	}
	if len(fn.Parameters) != 0 {
		t.Fatalf("function literal paramsters wrong. want 0, got=%d", len(fn.Parameters))
	}

	if len(fn.Body.Statements) != 1 {
		t.Fatalf("function.Body.Statements has not 1 statements. got=%d", len(fn.Body.Statements))
	}
	returnStmt, ok := fn.Body.Statements[0].(*ast.ReturnStatement)
	if !ok {
		t.Fatalf("function body stmt is not ast.ReturnStatement. got=%T", fn.Body.Statements[0])
	}
	testInfixExpression(t, returnStmt.ReturnValue, "x", "+", "y")
}

func TestCallExpressionParsing(t *testing.T) {
	input := "add(1, 2*3, 4 + 5);"
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParseErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain %d statements. got=%d", 1, len(program.Statements))
	}
	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])
	}
	exp, ok := stmt.Expression.(*ast.CallExpression)
	if !ok {
		t.Fatalf("stmt.Expression is not ast.CallExpression. got=%T", stmt.Expression)
	}
	if !testIdentifier(t, exp.Function, "add") {
		return
	}
	if len(exp.Arguments) != 3 {
		t.Fatalf("wrong length of arguments. got=%d", len(exp.Arguments))
	}
	testLiteralExpression(t, exp.Arguments[0], 1)
	testInfixExpression(t, exp.Arguments[1], 2, "*", 3)
	testInfixExpression(t, exp.Arguments[2], 4, "+", 5)
}
