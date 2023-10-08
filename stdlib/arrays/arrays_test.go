package arrays

import (
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
	env := object.NewEnvironment()
	env.RegisterPackages(New())
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

func TestLen(t *testing.T) {
	arr := &object.Array{Elements: []object.Object{
		&object.Integer{Value: 1},
		&object.Integer{Value: 2},
		&object.Integer{Value: 3},
	}}
	ap := New()
	arrLen := ap.Member("length")(arr)
	checkInteger(t, arrLen, 3)
}

func TestPush(t *testing.T) {
	arr := &object.Array{Elements: []object.Object{
		&object.Integer{Value: 1},
		&object.Integer{Value: 2},
		&object.Integer{Value: 3},
	}}

	p := New()
	ret := p.Member("push")(arr, &object.Integer{Value: 4})
	checkIntegerArray(t, ret, []int64{1, 2, 3, 4})
}

func TestInitLast(t *testing.T) {
	arr := &object.Array{Elements: []object.Object{
		&object.Integer{Value: 1},
		&object.Integer{Value: 2},
		&object.Integer{Value: 3},
	}}
	p := New()
	initRet := p.Member("init")(arr)
	lastRet := p.Member("last")(arr)

	checkIntegerArray(t, initRet, []int64{1, 2})
	checkInteger(t, lastRet, 3)
}

func TestHeadTail(t *testing.T) {
	arr := &object.Array{Elements: []object.Object{
		&object.Integer{Value: 1},
		&object.Integer{Value: 2},
		&object.Integer{Value: 3},
	}}
	p := New()
	headRet := p.Member("head")(arr)
	tailRet := p.Member("tail")(arr)

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
		case string:
			errObj, ok := evaluated.(*object.Error)
			if !ok {
				t.Errorf("object is not Error. got=%T (%+v) <= %s", evaluated, evaluated, tt.input)
			}
			if errObj.Message != expected {
				t.Errorf("wrong error message. expected=%q, got=%q",
					expected, errObj.Message)
			}
		}
	}
}
