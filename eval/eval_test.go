package eval_test

import (
	"bytes"
	gofmt "fmt"
	"testing"

	"github.com/thingsme/thingscript/eval"
	"github.com/thingsme/thingscript/lexer"
	"github.com/thingsme/thingscript/object"
	"github.com/thingsme/thingscript/parser"
	"github.com/thingsme/thingscript/stdlib"
)

func testEval(input string) object.Object {
	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()
	for _, err := range p.Errors() {
		gofmt.Println("Parse Error:", err)
	}
	env := object.NewEnvironment()
	env.RegisterPackages(stdlib.Packages()...)
	return eval.Eval(program, env)
}

func TestEvalIntegerExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"5", 5},
		{"10", 10},
		{"-5", -5},
		{"-10", -10},
		{"5 + 5 + 5 + 5 - 10", 10},
		{"2 * 2 * 2 * 2 * 2", 32},
		{"-50 + 100 + -50", 0},
		{"5 * 2 +10", 20},
		{"5 + 2 * 10", 25},
		{"20 + 2 * -10", 0},
		{"50 / 2 * 2 + 10", 60},
		{"2 * (5 + 10)", 30},
		{"3 * 3 * 3 + 10", 37},
		{"3 * (3 * 3) + 10", 37},
		{"(5 + 10 *2 +15 / 3) * 2 +-10", 50},
		{"13 % 10", 3},
	}
	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testIntegerObject(t, evaluated, tt.expected)
	}
}

func testIntegerObject(t *testing.T, obj object.Object, expected int64) bool {
	t.Helper()
	result, ok := obj.(*object.Integer)
	if !ok {
		t.Errorf("object is not Integer, got=%T (%+v)", obj, obj)
	}
	if result.Value != expected {
		t.Errorf("object has wrong value. got=%d, want=%d", result.Value, expected)
		return false
	}
	return true
}

func checkBoolean(t *testing.T, obj object.Object, expect bool) {
	t.Helper()
	boolObj, ok := obj.(*object.Boolean)
	if !ok {
		t.Errorf("obj is not an integer object, got=%T", obj)
	}
	if boolObj.Value != expect {
		t.Errorf("boolean different, expect %t, got=%t", expect, boolObj.Value)
	}
}

func TestEvalFloatExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected float64
	}{
		{"3.14", 3.14},
		{"10.0", 10},
		{"-5.0", -5},
		{"-10.1", -10.1},
		{"5.0 + 5.0 + 5.0 + 5.0 - 10.0", 10},
		{"2.0 * 2.0 * 2.0 * 2.0 * 2.0", 32},
		{"-50.0 + 100.0 + -50.0", 0},
		{"5.0 * 2.0 + 10.0", 20},
		{"5.0 + 2.0 * 10.0", 25},
		{"20.0 + 2.0 * -10.0", 0},
		{"50.0 / 2.0 * 2.0 + 10.0", 60},
		{"2.0 * (5.0 + 10.0)", 30},
		{"3.0 * 3.0 * 3.0 + 10.0", 37},
		{"3.0 * (3.0 * 3.0) + 10.0", 37},
		{"(5.0 + 10.0 * 2.0 + 15.0 / 3.0) * 2.0 +-10.0", 50},
		{"1 + 2.3", 3.3},
		{"1.2 + 3", 4.2},
	}
	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testFloatObject(t, evaluated, tt.expected)
	}
}

func testFloatObject(t *testing.T, obj object.Object, expected float64) bool {
	result, ok := obj.(*object.Float)
	if !ok {
		t.Errorf("object is not Float, got=%T (%+v)", obj, obj)
	}
	if result.Value != expected {
		t.Errorf("object has wrong value. got=%f, want=%f", result.Value, expected)
		return false
	}
	return true
}

func testNullObject(t *testing.T, obj object.Object) bool {
	if obj != eval.NULL {
		t.Errorf("object is not NULL. got=%T (%+v)", obj, obj)
		return false
	}
	return true
}

func TestEvalBooleanExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"true", true},
		{"false", false},
		{"1 < 2", true},
		{"1.0 < 2", true},
		{"1 <= 2", true},
		{"2 <= 2", true},
		{"1 > 2", false},
		{"1 >= 2", false},
		{"2 >= 2", true},
		{"2 >= 2.0", true},
		{"1 < 1", false},
		{"1 > 1", false},
		{"1 == 1", true},
		{"1 != 1", false},
		{"1 == 2", false},
		{"1 != 2", true},
		{"true == true", true},
		{"false == false", true},
		{"true == false", false},
		{"true != false", true},
		{"false != true", true},
		{"(1 < 2) == true", true},
		{"(1 < 2) == false", false},
		{"(1 > 2) == true", false},
		{"(1 > 2) == false", true},
	}
	for _, tt := range tests {
		evaluated := testEval(tt.input)
		result, ok := evaluated.(*object.Boolean)
		if !ok {
			t.Errorf("object is not Boolean, got=%T (%+v)", evaluated, evaluated)
		}
		if result.Value != tt.expected {
			t.Errorf("object has wrong value. got=%t, want=%t <- %s", result.Value, tt.expected, tt.input)
		}
	}
}

func testBooleanObject(t *testing.T, obj object.Object, expected bool) bool {
	result, ok := obj.(*object.Boolean)
	if !ok {
		t.Errorf("object is not Boolean, got=%T (%+v)", obj, obj)
	}
	if result.Value != expected {
		t.Errorf("object has wrong value. got=%t, want=%t <- ", result.Value, expected)
		return false
	}
	return true
}

func TestStringLiteral(t *testing.T) {
	input := `"Hello World"`
	evaluated := testEval(input)
	str, ok := evaluated.(*object.String)
	if !ok {
		t.Fatalf("object is not String. got=%T (%+v)", evaluated, evaluated)
	}
	if str.Value != "Hello World" {
		t.Errorf("String has wrong value. got=%q", str.Value)
	}
}

func TestStringConcatenation(t *testing.T) {
	input := `"Hello"+ " "+ "World"`
	evaluated := testEval(input)
	str, ok := evaluated.(*object.String)
	if !ok {
		t.Fatalf("object is not String. got=%T (%+v)", evaluated, evaluated)
	}
	if str.Value != "Hello World" {
		t.Errorf("String has wrong value. got=%q", str.Value)
	}
}

func TestBooleanLiterals(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{`var x = true; x`, true},
		{`var x = false; x`, false},
		{`true`, true},
		{`false`, false},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		if obj, ok := evaluated.(*object.Boolean); !ok {
			t.Errorf("expect %t, got=%T(%v)", tt.expected, obj, obj)
		} else if obj.Value != tt.expected {
			t.Errorf("expect %t, got=%T(%v)", tt.expected, obj, obj)
		}
	}
}

func TestArrayLiterals(t *testing.T) {
	input := `[1, 2 + 2, 3 * 3]`
	evaluated := testEval(input)
	result, ok := evaluated.(*object.Array)
	if !ok {
		t.Fatalf("object is not Array. got=%T (%+v)", evaluated, evaluated)
	}
	if len(result.Elements) != 3 {
		t.Fatalf("array has wrong num of elements, got=%d", len(result.Elements))
	}
	testIntegerObject(t, result.Elements[0], 1)
	testIntegerObject(t, result.Elements[1], 4)
	testIntegerObject(t, result.Elements[2], 9)
}

func TestImmediateIfExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected any
	}{
		{`var v = 10; v ?? 20`, 10},
		{`var v = nil; v ?? 20`, 20},
		{`var v = nil; 5 * (v ?? 20)`, 100},
		{`func v(){ return nil }; v() ?? 20`, 20},
		{`var v = func(){ return nil }; var x = func() { return 20 }; v() ?? x()`, 20},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		integer, ok := tt.expected.(int)
		if ok {
			testIntegerObject(t, evaluated, int64(integer))
		} else {
			testNullObject(t, evaluated)
		}
	}
}

func TestWhileExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected any
	}{
		{`var sum = 0; var v = 0; while v < 10 { v += 1; sum += v; }; sum`, 55},
		{`var sum = 0; var v = 0; while v < 20 { v += 1; sum += v; if (v == 10) { break } }; sum`, 55},
		{`var sum = 0; func run(){ var v = 0; while v < 20 { v += 1; sum += v; if (v == 10) { return 10; } } };  run(); sum`, 55},
		{`var sum = 0; func run(){ var v = 0; while v < 20 { v += 1; sum += v; if (v == 10) { return; } } };  run(); sum`, 55},
		{`var sum = 0; func run(){ var v = 0; while v < 20 { v += 1; sum += v; if (v == 10) { return } } };  run(); sum`, 55},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		integer, ok := tt.expected.(int)
		if ok {
			testIntegerObject(t, evaluated, int64(integer))
		} else {
			testNullObject(t, evaluated)
		}
	}
}

func TestDoWhileExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected any
	}{
		{`var sum = 0; var v = 0; do { v += 1; sum += v; } while v < 10 ; sum`, 55},
		{`var sum = 0; var v = 0; do { v += 1; sum += v; if (v == 10) { break } } while v < 20; sum`, 55},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		integer, ok := tt.expected.(int)
		if ok {
			testIntegerObject(t, evaluated, int64(integer))
		} else {
			testNullObject(t, evaluated)
		}
	}
}

func TestArrayIndexExpressions(t *testing.T) {
	tests := []struct {
		input    string
		expected any
	}{
		{`[1,2,3][0]`, 1},
		{`[1,2,3][1]`, 2},
		{`[1,2,3][2]`, 3},
		{`var i = 0; [1][i]`, 1},
		{`[1, 2, 3][1 + 1]`, 3},
		{`var myArray = [1, 2, 3]; myArray[2]`, 3},
		{`var myArray = [1, 2, 3]; myArray[0] + myArray[1] + myArray[2]`, 6},
		{`var myArray = [1, 2, 3]; var i = myArray[0]; myArray[i]`, 2},
		{`[1, 2, 3][3]`, nil},
		{`[1, 2, 3][-1]`, nil},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		integer, ok := tt.expected.(int)
		if ok {
			testIntegerObject(t, evaluated, int64(integer))
		} else {
			testNullObject(t, evaluated)
		}
	}
}

func TestHashLiterals(t *testing.T) {
	input := `var two = "two";
	{
		"one": 10 - 9,
		two: 1 + 1,
		"thr"+"ee": 6 /2,
		4: 4,
		true: 5,
		false: 6
	}
	`
	evaluated := testEval(input)
	result, ok := evaluated.(*object.HashMap)
	if !ok {
		t.Fatalf("Eval didn't return Hash. got=%T (%+v)", evaluated, evaluated)
	}
	expected := map[object.HashKey]int64{
		(&object.String{Value: "one"}).HashKey():   1,
		(&object.String{Value: "two"}).HashKey():   2,
		(&object.String{Value: "three"}).HashKey(): 3,
		(&object.Integer{Value: 4}).HashKey():      4,
		(&object.Boolean{Value: true}).HashKey():   5,
		(&object.Boolean{Value: false}).HashKey():  6,
	}
	if len(result.Pairs) != len(expected) {
		t.Fatalf("Hash has wrong num of pairs. got=%d", len(result.Pairs))
	}
	for expectedKey, expectedValue := range expected {
		pair, ok := result.Pairs[expectedKey]
		if !ok {
			t.Errorf("no pair for given key in Pairs")
		}
		testIntegerObject(t, pair.Value, expectedValue)
	}
}

func TestHashIndexExpressions(t *testing.T) {
	tests := []struct {
		input    string
		expected any
	}{
		{`{"foo": 5}["foo"]`, 5},
		{`{"foo": 5}["bar"]`, nil},
		{`var key = "foo"; {"foo": 5}[key]`, 5},
		{`{}["foo"]`, nil},
		{`{5: 5}[5]`, 5},
		{`{true: 5}[true]`, 5},
		{`{false: 5}[false]`, 5},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		integer, ok := tt.expected.(int)
		if ok {
			testIntegerObject(t, evaluated, int64(integer))
		} else {
			testNullObject(t, evaluated)
		}
	}
}

func TestBangOperator(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{`!true`, false},
		{`!false`, true},
		{`!5`, false},
		{`!!true`, true},
		{`!!false`, false},
		{`!!5`, true},
	}
	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testBooleanObject(t, evaluated, tt.expected)
	}
}

func TestIfElseExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected any
	}{
		{"if (true){ 10 }", 10},
		{"if (false){ 10 }", nil},
		{"if true { 10 }", 10},
		{"if false{ 10 }", nil},
		{"if (1){ 10 }", 10},
		{"if (1 < 2){ 10 }", 10},
		{"if (1 > 2){ 10 }", nil},
		{"if 1 < 2{ 10 }", 10},
		{"if 1 > 2{ 10 }", nil},
		{"if (1 > 2){ 10 } else { 20 }", 20},
		{"if (1 < 2){ 10 } else { 20 }", 10},
		{"if 1 > 2 { 10 } else { 20 }", 20},
		{"if 1 < 2 { 10 } else { 20 }", 10},
		{`if "abc" < "bcd" { 10 } else {20}`, 10},
		{`if "abc" > "bcd" { 10 } else {20}`, 20},
		{`if "abc" != "bcd" { 10 } else {20}`, 10},
		{`if "abc" == "bcd" { 10 } else {20}`, 20},
		{"if nil { 10 } else { 20 }", 20},
	}
	for _, tt := range tests {
		evaluated := testEval(tt.input)
		integer, ok := tt.expected.(int)
		if ok {
			testIntegerObject(t, evaluated, int64(integer))
		} else {
			testNullObject(t, evaluated)
		}
	}
}

func TestReturnStatement(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"return 10;", 10},
		{"return 10; 9;", 10},
		{"return 2 * 5; 9;", 10},
		{"9; return 2 * 5;", 10},
		{`if (10 > 1) { return 10; } return 1; `, 10},
		{`func() { return ( if (10 > 1) { nil } else { 1 } ) }() ?? 10 `, 10},
	}
	for _, tt := range tests {
		evalutated := testEval(tt.input)
		testIntegerObject(t, evalutated, tt.expected)
	}
}

func TestErrorHandling(t *testing.T) {
	tests := []struct {
		input           string
		expectedMessage string
	}{
		{"5 + true", "type mismatch: INTEGER + BOOLEAN"},
		{"5 + true; 5;", "type mismatch: INTEGER + BOOLEAN"},
		{"-true", "unknown operator: -BOOLEAN"},
		{"true + false;", "unknown operator: BOOLEAN + BOOLEAN"},
		{"5; true + false; 5", "unknown operator: BOOLEAN + BOOLEAN"},
		{"if (10 > 1) { true + true;}", "unknown operator: BOOLEAN + BOOLEAN"},
		{`
			if (10 > 1) {
				if (10 > 1) {
					return true + false;
				}
			}
			return 1;
		`, "unknown operator: BOOLEAN + BOOLEAN"},
		{"foobar", "identifier not found: foobar"},
		{"foo = 10", "identifier not found: foo"},
		{`"Hello" - "World"`, "unknown operator: STRING - STRING"},
		{`{"name": "Monkey"}[func(x){x}]`, "unusable as hash key: FUNCTION"},
	}
	for _, tt := range tests {
		evalutated := testEval(tt.input)
		errObj, ok := evalutated.(*object.Error)
		if !ok {
			t.Errorf("no error object returned, got=%T (%+v)", evalutated, evalutated)
		}
		if errObj.Message != tt.expectedMessage {
			t.Errorf("wrong error message. expected=%q, got=%q",
				tt.expectedMessage, errObj.Message)
		}
	}
}

func TestVarStatement(t *testing.T) {
	tests := []struct {
		input    string
		expected any
	}{
		{"var a = 5; a;", 5},
		{"var a = 5 * 5; a;", 25},
		{"var a = 5; var b = a; b;", 5},
		{"var a = 5; var b = a; var c = a + b + 5; c;", 15},
		{"a := 5; a;", 5},
		{"a := 5 * 5; a = a + 1; a;", 26},
		{"a := 5; b := a; b;", 5},
		{"a := 5; b := a; c := a + b + 5; c;", 15},
		{"v := 10; v += 10; v", 20},
		{"v := 10; v -= 10; v", 0},
		{"v := 12; v %= 10; v", 2},
		{"v := 13; v = v % 10; v", 3},
		{"v := 10.0;  func m() { return 10.2 };  v *= m(); v", 102.0},
		{"v := 100.0; v = v / 10.0; func m() { return 10.2}; v *= m(); v", 102.0},
		{"v := 103.0; v /= 10.3; v", 10.0},
	}
	for _, tt := range tests {
		switch expected := tt.expected.(type) {
		case int:
			testIntegerObject(t, testEval(tt.input), int64(expected))
		case float64:
			testFloatObject(t, testEval(tt.input), expected)
		}
	}
}

func TestFunctionObject(t *testing.T) {
	input := "func(x) { x + 2; };"

	evaluated := testEval(input)
	fn, ok := evaluated.(*object.Function)
	if !ok {
		t.Fatalf("object is not Function, got=%T (%+v)", evaluated, evaluated)
	}
	if len(fn.Parameters) != 1 {
		t.Fatalf("function has wrong parameters. Parameters=%+v", fn.Parameters)
	}
	if fn.Parameters[0].String() != "x" {
		t.Fatalf("parameter is not 'x'. got=%q", fn.Parameters[0])
	}
	expectedBody := "(x + 2)"
	if fn.Body.String() != expectedBody {
		t.Fatalf("body is not %q. got=%q", expectedBody, fn.Body.String())
	}
}

func TestFunctionApplication(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"var identity = func(x) {x;}; identity(5);", 5},         // implicit return value
		{"var identity = func(x) { return x;}; identity(5);", 5}, // return statement
		{"var double = func(x) {x * 2;}; double(5);", 10},
		{"var add = func(x, y) {x + y;}; add(5, 5);", 10},
		{"var add = func(x, y) {x + y;}; add(5 +5, add(5, 5));", 20},
		{"func(x){x;}(5)", 5},
		{"func identity(x) {x;}; identity(5);", 5},         // implicit return value
		{"func identity(x) { return x;}; identity(5);", 5}, // return statement
		{"func double(x) {x * 2;}; double(5);", 10},
		{"func add (x, y) {x + y;}; add(5, 5);", 10},
		{"func add (x, y) {x + y;}; add(5 +5, add(5, 5));", 20},
		{"func(x){x;}(5)", 5},
	}

	for _, tt := range tests {
		testIntegerObject(t, testEval(tt.input), tt.expected)
	}
}

func TestClosure(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{
			input: `
				var newAdder = func(x) {
					func(y) { x + y };
				};
				var addTwo = newAdder(2);
				addTwo(3);
			`,
			expected: 5,
		},
		{
			input: `
				newAdder := func(x) {
					func(y) { x + y };
				};
				var addTwo = newAdder(2);
				addTwo(3);
			`,
			expected: 5,
		},
		{
			input: `
				func newAdder(x) {
					func(y) { x + y };
				};
				var addTwo = newAdder(2);
				addTwo(3);
			`,
			expected: 5,
		},
	}
	for _, tt := range tests {
		testIntegerObject(t, testEval(tt.input), tt.expected)
	}
}

func TestBuiltinFunctionError(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{`1.length()`, "identifier not found: length"},
		{`1.length`, "identifier not found: length"},
		{`"one".length("two")`, "wrong number of arguments. got=1, want=0"},
	}
	for _, tt := range tests {
		evaluated := testEval(tt.input)
		obj, ok := evaluated.(*object.Error)
		if !ok {
			t.Errorf("Eval Error, it should be error, got=%T", evaluated)
		}
		if obj.Message != tt.expected {
			t.Errorf("wrong error message. expected=%q, got=%q",
				tt.expected, obj.Message)
		}
	}
}

func TestBuiltinFunction(t *testing.T) {
	tests := []struct {
		input    string
		expected any
	}{
		{`"".length()`, 0},
		{`"".length`, 0},
		{`"four".length()`, 4},
		{`"four".length`, 4},
		{`"hello world".length()`, 11},
		{`("hello" + " " + "world").length`, 11},
		{`[1, 2, 3].length()`, 3},
		{`[1, 2, 3].length`, 3},
		{`[1, 2, 3].head()`, 1},
		{`[1,2,3].tail().tail().length()`, 1},
		{`[1,2,3].tail.tail.length`, 1},
		{`[1,2,3].tail().tail()[0]`, 3},
		{`[1,2,3].tail().tail[0]`, 3},
		{`[1, 2, 3].last()`, 3},
		{`[1,2,3].init().length()`, 2},
		{`[1,2,3].init()[0]`, 1},
		{`[1,2,3].init()[1]`, 2},
		{`sum := 0; [1,2,3].foreach(func(idx,elm){ sum += elm}); sum`, 6},
		{`sum := 0; func iter(idx, elm){ sum += elm}; [1,2,3].foreach(iter); sum`, 6},
		{`sum := 0; iter := func(idx, elm){ sum += elm}; [1,2,3].foreach(iter); sum`, 6},
		{`sum := 0.0; [1.1,2.2,3.3].foreach(func(idx,elm){ sum += elm}); sum`, 6.6},
		{`sum := ""; ["1","2","3"].foreach(func(idx,elm){ sum += elm}); sum`, "123"},
		{`sum := ""; func cat(idx, elm){ sum+=elm}; ["1","2","3"].foreach(cat); sum`, "123"},
		{`sum := ""; cat := func(idx, elm){ sum+=elm}; ["1","2","3"].foreach(cat); sum`, "123"},
		{`ret := true; [true, true, false].foreach(func(idx,elm){ ret = elm }); ret`, false},
		{`ret := true; func iter(idx,elm){ ret = elm }; [true, true, false].foreach(iter); ret`, false},
		{`ret := true; var iter = func(idx,elm){ ret = elm }; [true, true, false].foreach(iter); ret`, false},
		{`ret := true; iter := func(idx,elm){ ret = elm }; [true, true, false].foreach(iter); ret`, false},
		{`func arr(){return [1,2,3]}; arr().head()`, 1},
		{`func arr(){return [1,2,3]}; arr().last()`, 3},
		{`var b = [1,2,3].push(4); b[3]`, 4},
		{`h := {1:"a", 2:"b"}; h.length()`, 2},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		if evaluated != nil && evaluated.Type() == object.ERROR_OBJ {
			t.Error("Eval Error:", evaluated, "<=", tt.input)
		}
		switch expected := tt.expected.(type) {
		case int:
			testIntegerObject(t, evaluated, int64(expected))
		case bool:
			checkBoolean(t, evaluated, expected)
		case string:
			obj, ok := evaluated.(*object.String)
			if !ok {
				t.Error("Eval Error: not string", evaluated, "<=", tt.input)
			}
			if obj.Value != expected {
				t.Errorf("wrong string. expected=%q, got=%q", expected, obj.Value)
			}
		case float64:
			testFloatObject(t, evaluated, expected)
		default:
			t.Errorf("wrong test type. expected=%q, got=%q", expected, evaluated)
		}
	}
}

func TestAccessOperation(t *testing.T) {
	tests := []struct {
		input    string
		expected any
	}{
		{`"hello".length()`, 5},
		{`("hello"+", world").length()`, 12},
		{`[1,2,3].length()`, 3},
		{`[].length()`, 0},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		if evaluated.Type() == object.ERROR_OBJ {
			t.Errorf("Error: %s <= %s", evaluated.Inspect(), tt.input)
			continue
		}
		switch expected := tt.expected.(type) {
		case int:
			testIntegerObject(t, evaluated, int64(expected))
		}
	}
}

func TestImports(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			input: `
				var out = import("fmt")
				out.println(1, 2, "x")
			`,
			expected: "1 2 x\n",
		},
		{
			input: `
				out := import("fmt")
				out.println(1, 2, "x")
			`,
			expected: "1 2 x\n",
		},
	}
	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := parser.New(l)
		program := p.ParseProgram()
		env := object.NewEnvironment()
		out := &bytes.Buffer{}
		env.RegisterPackages(stdlib.FmtPackage(stdlib.WithWriter(out)))
		ret := eval.Eval(program, env)
		if ret != nil && ret.Type() == object.ERROR_OBJ {
			t.Errorf("result is error; %s", ret.Inspect())
		}
		if out.String() != tt.expected {
			t.Errorf("result is not %q, got=%q", tt.expected, out.String())
		}
	}
}
