package ast_test

import (
	"testing"

	"github.com/thingsme/thingscript/lexer"
	"github.com/thingsme/thingscript/parser"
)

func TestAst(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{`return`, "return;"},
		{`1 +-2`, "(1 + (-2))"},
		{`var myVar = anotherVar`, "var myVar = anotherVar;"},
		{`myVar = 10`, "myVar = 10;"},
		{`myVar += 10`, "myVar += 10;"},
		{`func myFn(a, b){ return a + b}`, "func <myFn>(a, b) {return (a + b);}"},
		{`var myFn = func(){return true}`, "var myFn = func<myFn>() { return true; };"},
		{`if a < b { break }`, "if (a < b) { break; }"},
		{`a = b ?? 10`, "a = (b ?? 10);"},
		{`while(true) { a += 1}`, "while ( true ) { a += 1; }"},
		{`do{ a += 1 } while(true)`, "do { a += 1; } while (true);"},
		{`a = [1, 2, 3]`, "a = [1, 2, 3];"},
		{`h = {a: "alpha", b: "beta"}`, "h = {a:alpha, b:beta};"},
		{`call(ab, cd)`, "call(ab, cd)"},
		{`arr[12]`, "(arr[12])"},
		{`"hello-world".length`, "((hello-world).(length))"},
		{`[1, 2, 3].field`, `(([1, 2, 3]).(field))`},
		{`[1, 2, 3].call()`, `(([1, 2, 3]).(call()))`},
		{`[1, 2, 3].call(true)`, `(([1, 2, 3]).(call(true)))`},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := parser.New(l)
		ret := p.ParseProgram()
		for _, err := range p.Errors() {
			t.Error("Parse Error:", err)
		}
		if ret.String() != tt.expected {
			t.Errorf("AST wrong. got=%q", ret.String())
		}
	}
}
