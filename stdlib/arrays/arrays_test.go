package arrays

import (
	"testing"

	"github.com/thingsme/thingscript/object"
)

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
