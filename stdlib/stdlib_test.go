package stdlib

import (
	"fmt"
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
		`v := 10; v.type == "integer"`,
		`v := 12.3; v.type == "float"`,
		`v := "1234"; v.type == "string"`,
		`v := true; v.type == "boolean"`,
		`v := [1,2]; v.type == "array"`,
		`v := {"a":1,"b":2.3}; v.type == "hashmap"`,
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
		{`v := 10; v.type(1)`, "wrong number of arguments. got=1, want=0"},
		{`v := 12.3; v.type(1)`, "wrong number of arguments. got=1, want=0"},
		{`v := "1234"; v.type(1)`, "wrong number of arguments. got=1, want=0"},
		{`v := true; v.type(1)`, "wrong number of arguments. got=1, want=0"},
		{`v := [1,2,3]; v.type(1)`, "wrong number of arguments. got=1, want=0"},
		{`v := [1,2,3]; v.head(1)`, "wrong number of arguments. got=1, want=0"},
		{`v := [1,2,3]; v.tail(1)`, "wrong number of arguments. got=1, want=0"},
		{`v := [1,2,3]; v.init(1)`, "wrong number of arguments. got=1, want=0"},
		{`v := [1,2,3]; v.last(1)`, "wrong number of arguments. got=1, want=0"},
		{`v := {"a":1,"b":2.3}; v.type(1)`, "wrong number of arguments. got=1, want=0"},
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
		{`v := "1234"; v.length(1)`, "wrong number of arguments. got=1, want=0"},
		{`v := [1,2]; v.length(1)`, "wrong number of arguments. got=1, want=0"},
		{`v := {"a":1,"b":2.3}; v.length(1)`, "wrong number of arguments. got=1, want=0"},
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
