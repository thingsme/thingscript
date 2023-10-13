package stdlib

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/thingsme/thingscript/eval"
	"github.com/thingsme/thingscript/lexer"
	"github.com/thingsme/thingscript/object"
	"github.com/thingsme/thingscript/parser"
)

func testEval(input string) object.Object {
	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()
	for _, err := range p.Errors() {
		fmt.Println("Parse Error:", err)
	}
	env := object.NewEnvironment()
	env.RegisterPackages(Packages()...)
	return eval.Eval(program, env)
}

func runTest(t *testing.T, input string, expected any) {
	t.Helper()
	l := lexer.New(input)
	p := parser.New(l)
	for _, err := range p.Errors() {
		t.Errorf("Parse Error: %v", err)
	}
	program := p.ParseProgram()
	env := object.NewEnvironment()
	env.RegisterPackages(Packages()...)
	ret := eval.Eval(program, env)
	if err, ok := ret.(*object.Error); ok {
		if exp, ok := expected.(*object.Error); ok {
			if exp.Message != err.Message {
				t.Errorf("Expect error %q, got=%q <= %s", exp.Message, err.Message, input)
				return
			}
		} else {
			t.Errorf("Result error, %s <= %s", err.Message, input)
		}
		return
	}

	switch v := expected.(type) {
	case int:
		checkInteger(t, ret, int64(v))
	case float64:
		checkFloat(t, ret, v)
	case string:
		checkString(t, ret, v)
	case bool:
		checkBoolean(t, ret, v)
	case []int64:
		checkIntegerArray(t, ret, v)
	default:
		t.Fatalf("unsupproted check type %T", expected)
	}
}

func TestUndefinedMemberError(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{`v := 10; v.undefined(1)`, `function "undefined" not found in "INTEGER"`},
		{`v := 12.3; v.undefined(1)`, `function "undefined" not found in "FLOAT"`},
		{`v := "1234"; v.undefined(1)`, `function "undefined" not found in "STRING"`},
		{`v := true; v.undefined(1)`, `function "undefined" not found in "BOOLEAN"`},
		{`v := [1,2]; v.undefined(1)`, `function "undefined" not found in "ARRAY"`},
		{`v := {"a":1,"b":2.3}; v.undefined(1)`, `function "undefined" not found in "HASHMAP"`},
	}
	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := parser.New(l)
		program := p.ParseProgram()
		env := object.NewEnvironment()
		env.RegisterPackages(Packages()...)
		ret := eval.Eval(program, env)
		if obj, ok := ret.(*object.Error); !ok {
			t.Errorf("Result not error, got=%T (%+v)", ret, ret)
		} else if obj.Message != tt.expected {
			t.Errorf("Result fail; got=%s <= %s", obj.Message, tt.input)
		}
	}
}

func TestType(t *testing.T) {
	tests := []string{
		`v := 10; v.type == "int"`,
		`v := 12.3; v.type == "float"`,
		`v := "1234"; v.type == "string"`,
		`v := true; v.type == "bool"`,
		`v := [1,2]; v.type == "array"`,
		`v := {"a":1,"b":2.3}; v.type == "map"`,
	}
	for _, tt := range tests {
		l := lexer.New(tt)
		p := parser.New(l)
		program := p.ParseProgram()
		env := object.NewEnvironment()
		env.RegisterPackages(Packages()...)
		ret := eval.Eval(program, env)
		if b, ok := ret.(*object.Boolean); !ok {
			t.Errorf("Result not boolean, got=%T (%+v)", ret, ret)
		} else if !b.Value {
			t.Error("Result fail <= ", tt)
		}
	}
}

func TestTypeError(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{`v := 10; v.type(1)`, "wrong number of arguments. want=0 got=1"},
		{`v := 12.3; v.type(1)`, "wrong number of arguments. want=0 got=1"},
		{`v := "1234"; v.type(1)`, "wrong number of arguments. want=0 got=1"},
		{`v := true; v.type(1)`, "wrong number of arguments. want=0 got=1"},
		{`v := [1,2,3]; v.type(1)`, "wrong number of arguments. want=0 got=1"},
		{`v := [1,2,3]; v.head(1)`, "wrong number of arguments. want=0 got=1"},
		{`v := [1,2,3]; v.tail(1)`, "wrong number of arguments. want=0 got=1"},
		{`v := [1,2,3]; v.init(1)`, "wrong number of arguments. want=0 got=1"},
		{`v := [1,2,3]; v.last(1)`, "wrong number of arguments. want=0 got=1"},
		{`v := {"a":1,"b":2.3}; v.type(1)`, "wrong number of arguments. want=0 got=1"},
		{`v := {"a":1,"b":2.3}; v.head(1)`, "function \"head\" not found in \"HASHMAP\""},
	}
	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := parser.New(l)
		program := p.ParseProgram()
		env := object.NewEnvironment()
		env.RegisterPackages(Packages()...)
		ret := eval.Eval(program, env)
		if obj, ok := ret.(*object.Error); !ok {
			t.Errorf("Result not error, got=%T (%+v)", ret, ret)
		} else if obj.Message != tt.expected {
			t.Errorf("Result fail got=%q <= %s", obj.Message, tt.input)
		}
	}
}

func TestPrimitives(t *testing.T) {
	tests := []struct {
		input    string
		expected any
	}{
		// int
		{`var a int; a`, 0},
		{`var a int = 1; a`, 1},
		{`a := int(2); a`, 2},
		{`a := 1; b := 2; c := a + b; c `, 3},
		{`a := 1; b := 2; c := a - b; c `, -1},
		{`a := 2; b := 3; c := a * b; c `, 6},
		{`a := 6; b := 2; c := a / b; c `, 3},
		{`a := 5; b := 2; c := a % b; c `, 1},
		{`a := 3; b := 2; c := a < b; c `, false},
		{`a := 2; b := 2; c := a <= b; c `, true},
		{`a := 4; b := 2; c := a > b; c `, true},
		{`a := 4; b := 2; c := a >= b; c `, true},
		{`a := 4; b := 2; c := a == b; c `, false},
		{`a := 4; b := 2; c := a != b; c `, true},
		// int & float
		{`a := 1; b := 2.0; c := a + b; c `, 3.0},
		{`a := 1; b := 2.0; c := a - b; c `, -1.0},
		{`a := 2; b := 3.0; c := a * b; c `, 6.0},
		{`a := 6; b := 2.0; c := a / b; c `, 3.0},
		{`a := 5; b := 2.0; c := a % b; c `, &object.Error{Message: "type mismatch: INTEGER % FLOAT"}},
		{`a := 3; b := 2.0; c := a < b; c `, false},
		{`a := 2; b := 2.0; c := a <= b; c `, true},
		{`a := 4; b := 2.0; c := a > b; c `, true},
		{`a := 4; b := 2.0; c := a >= b; c `, true},
		{`a := 4; b := 2.0; c := a == b; c `, false},
		{`a := 4; b := 2.0; c := a != b; c `, true},
		// float
		{`var a float; a`, 0.0},
		{`var a float = 1.2; a`, 1.2},
		{`a := 2.345; a`, 2.345},
		{`a := 2.345; a = a + 1; a`, 3.345},
		{`a := 2.345; a = a * 10; a`, 23.45},
		{`a := 2.468; a = a / 2; a`, 1.234},
		{`a := 2.345; a = a - 2; a`, 0.345},
		{`a := float(2); a`, 2.0},
		{`a := 3.1; b := 2.1; c := a < b; c `, false},
		{`a := 2.1; b := 2.1; c := a <= b; c `, true},
		{`a := 4.2; b := 2.1; c := a > b; c `, true},
		{`a := 4.0; b := 2.0; c := a >= b; c `, true},
		{`a := 4.0; b := 2.0; c := a == b; c `, false},
		{`a := 4.0; b := 2.0; c := a != b; c `, true},
		// string
		{`var a string; a`, ""},
		{`var a string = "hello"; a`, "hello"},
		{`a := "world"; a`, "world"},
		{`a := "hello"; b := "world"; c := a + b; c`, "helloworld"},
		{`a := "a"; b := "b"; c := a < b; c`, true},
		{`a := "a"; b := "b"; c := a <= b; c`, true},
		{`a := "a"; b := "b"; c := a > b; c`, false},
		{`a := "a"; b := "b"; c := a >= b; c`, false},
		{`a := "a"; b := "b"; c := a == b; c`, false},
		{`a := "a"; b := "b"; c := a != b; c`, true},
		// bool
		{`var a bool; a`, false},
		{`var a bool = 1 < 2; a`, true},
		{`a := 2 > 1; a`, true},
		{`a := true; var b bool; a == b`, false},
		{`a := true; var b bool; a != b`, true},
		// array
		{`var a = [1,2,3]; a`, []int64{1, 2, 3}},
		{`var a = [1,2,3]; var b array; b = a; b`, []int64{1, 2, 3}},
		// TODO; lvalue and rvalue
		// {`var a = [1,2,3]; a[0] = a[0]*10; a`, []int64{10, 2, 3}},
	}
	for _, tt := range tests {
		runTest(t, tt.input, tt.expected)
	}
}

func TestLength(t *testing.T) {
	tests := []string{
		`v := "1234"; v.length == 4`,
		`v := [1,2]; v.length == 2`,
		`v := {"a":1,"b":2.3}; v.length == 2`,
	}
	for _, tt := range tests {
		l := lexer.New(tt)
		p := parser.New(l)
		program := p.ParseProgram()
		env := object.NewEnvironment()
		env.RegisterPackages(Packages()...)
		ret := eval.Eval(program, env)
		if b, ok := ret.(*object.Boolean); !ok {
			t.Errorf("Result not boolean, got=%T (%+v)", ret, ret)
		} else if !b.Value {
			t.Error("Result fail <= ", tt)
		}
	}
}

func TestLengthError(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{`v := "1234"; v.length(1)`, "wrong number of arguments. want=0 got=1"},
		{`v := [1,2]; v.length(1)`, "wrong number of arguments. want=0 got=1"},
		{`v := {"a":1,"b":2.3}; v.length(1)`, "wrong number of arguments. want=0 got=1"},
	}
	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := parser.New(l)
		program := p.ParseProgram()
		env := object.NewEnvironment()
		env.RegisterPackages(Packages()...)
		ret := eval.Eval(program, env)
		if obj, ok := ret.(*object.Error); !ok {
			t.Errorf("Result not error, got=%T (%+v)", ret, ret)
		} else if obj.Message != tt.expected {
			t.Error("Result fail <= ", tt)
		}
	}
}

func TestPush(t *testing.T) {
	arr := &object.Array{Elements: []object.Object{
		&object.Integer{Value: 1},
		&object.Integer{Value: 2},
		&object.Integer{Value: 3},
	}}

	ret := Arrays("push")(arr, &object.Integer{Value: 4})
	checkIntegerArray(t, ret, []int64{1, 2, 3, 4})
}

func TestInitLast(t *testing.T) {
	arr := &object.Array{Elements: []object.Object{
		&object.Integer{Value: 1},
		&object.Integer{Value: 2},
		&object.Integer{Value: 3},
	}}
	initRet := Arrays("init")(arr)
	lastRet := Arrays("last")(arr)

	checkIntegerArray(t, initRet, []int64{1, 2})
	checkInteger(t, lastRet, 3)
}

func TestHeadTail(t *testing.T) {
	arr := &object.Array{Elements: []object.Object{
		&object.Integer{Value: 1},
		&object.Integer{Value: 2},
		&object.Integer{Value: 3},
	}}
	headRet := Arrays("head")(arr)
	tailRet := Arrays("tail")(arr)

	checkInteger(t, headRet, 1)
	checkIntegerArray(t, tailRet, []int64{2, 3})
}

func TestFunctions(t *testing.T) {
	tests := []struct {
		input    string
		expected any
	}{
		{`[1, 2, 3].length()`, 3},
		{`[1, 2, 3].length`, 3},
		{`[1, 2, 3].head()`, 1},
		{`[1, 2, 3].head`, 1},
		{`[1, 2, 3].tail().length`, 2},
		{`[1, 2, 3].tail[0]`, 2},
		{`[1, 2, 3].tail[1]`, 3},
		{`[1,2,3].tail().tail().length()`, 1},
		{`[1,2,3].tail.tail.length`, 1},
		{`[1,2,3].tail().tail()[0]`, 3},
		{`[1,2,3].tail().tail[0]`, 3},
		{`[1, 2, 3].last()`, 3},
		{`[1, 2, 3].last`, 3},
		{`[1,2,3].init().length()`, 2},
		{`[1,2,3].init.length`, 2},
		{`[1,2,3].init()[0]`, 1},
		{`[1,2,3].init[0]`, 1},
		{`[1,2,3].init()[1]`, 2},
		{`[1,2,3].init[1]`, 2},
		{`sum := 0; [1,2,3].foreach(func(idx,elm){ sum += elm}); sum`, 6},
		{`sum := ""; ["1","2","3"].foreach(func(idx,elm){ sum += elm}); sum`, "123"},
		{`ret := true; [true, true, false].foreach(func(idx,elm){ ret = elm }); ret`, false},
		{`func arr(){return [1,2,3]}; arr().head()`, 1},
		{`func arr(){return [1,2,3]}; arr().head`, 1},
		{`func arr(){return [1,2,3]}; arr().last()`, 3},
		{`func arr(){return [1,2,3]}; arr().last`, 3},
		{`var b = [1,2,3].push(4); b[3]`, 4},
	}
	for _, tt := range tests {
		evaluated := testEval(tt.input)
		switch expected := tt.expected.(type) {
		case int:
			checkInteger(t, evaluated, int64(expected))
		case bool:
			checkBoolean(t, evaluated, expected)
		case string:
			switch obj := evaluated.(type) {
			case *object.String:
				if obj.Value != expected {
					t.Errorf("wrong string. expected=%q, got=%q",
						expected, obj.Value)
				}
			case *object.Error:
				if obj.Message != expected {
					t.Errorf("wrong error message. expected=%q, got=%q",
						expected, obj.Message)
				}
			default:
				t.Errorf("object is not Error. got=%T (%+v) <= %s", evaluated, evaluated, tt.input)
			}
		}
	}
}

func checkInteger(t *testing.T, obj object.Object, expect int64) {
	t.Helper()
	intObj, ok := obj.(*object.Integer)
	if !ok {
		t.Errorf("obj is not an integer object, got=%T", obj)
	}
	if intObj.Value != expect {
		t.Errorf("integer different, expect %d, got=%d", expect, intObj.Value)
	}
}

func checkFloat(t *testing.T, obj object.Object, expect float64) {
	t.Helper()
	floatObj, ok := obj.(*object.Float)
	if !ok {
		t.Errorf("obj is not an float object, got=%T", obj)
	}
	str := strconv.FormatFloat(expect, 'E', 7, 64)
	val := strconv.FormatFloat(floatObj.Value, 'E', 7, 64)
	if str != val {
		t.Errorf("float different, expect %f, got=%f", expect, floatObj.Value)
	}
	// if floatObj.Value != expect {
	// 	t.Errorf("float different, expect %f, got=%f", expect, floatObj.Value)
	// }
}

func checkString(t *testing.T, obj object.Object, expect string) {
	t.Helper()
	stringObj, ok := obj.(*object.String)
	if !ok {
		t.Errorf("obj is not an string object, got=%T", obj)
	}
	if stringObj.Value != expect {
		t.Errorf("string different, expect %s, got=%s", expect, stringObj.Value)
	}
}

func checkBoolean(t *testing.T, obj object.Object, expect bool) {
	t.Helper()
	boolObj, ok := obj.(*object.Boolean)
	if !ok {
		t.Errorf("obj is not an boolean object, got=%T", obj)
	}
	if boolObj.Value != expect {
		t.Errorf("boolean different, expect %t, got=%t", expect, boolObj.Value)
	}
}

func checkIntegerArray(t *testing.T, obj object.Object, expectArr []int64) {
	t.Helper()
	arrObj, ok := obj.(*object.Array)
	if !ok {
		t.Errorf("obj is not an array object")
	}
	if len(arrObj.Elements) != len(expectArr) {
		t.Errorf("elements length different, expect %d, got=%d (%+v)", len(expectArr), len(arrObj.Elements), expectArr)
	}
	for i, expect := range expectArr {
		intObj, ok := arrObj.Elements[i].(*object.Integer)
		if !ok {
			t.Errorf("element[%d] is not an integer, got=%T", i, arrObj.Elements[i])
			return
		}
		if intObj.Value != expect {
			t.Errorf("element[%d] is not %d, got=%d", i, expect, intObj.Value)
			return
		}
	}
}
