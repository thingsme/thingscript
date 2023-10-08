package object

import "testing"

func TestLen(t *testing.T) {
	arr := &Array{Elements: []Object{
		&Integer{Value: 1},
		&Integer{Value: 2},
		&Integer{Value: 3},
	}}
	arrLen := _length(arr)
	checkInteger(t, arrLen, 3)

	str := &String{Value: "1234"}
	strLen := _length(str)
	checkInteger(t, strLen, 4)
}

func TestPush(t *testing.T) {
	arr := &Array{Elements: []Object{
		&Integer{Value: 1},
		&Integer{Value: 2},
		&Integer{Value: 3},
	}}

	ret := _push(arr, &Integer{Value: 4})
	checkIntegerArray(t, ret, []int64{1, 2, 3, 4})
}

func TestInitLast(t *testing.T) {
	arr := &Array{Elements: []Object{
		&Integer{Value: 1},
		&Integer{Value: 2},
		&Integer{Value: 3},
	}}
	initRet := _init(arr)
	lastRet := _last(arr)

	checkIntegerArray(t, initRet, []int64{1, 2})
	checkInteger(t, lastRet, 3)
}

func TestHeadTail(t *testing.T) {
	arr := &Array{Elements: []Object{
		&Integer{Value: 1},
		&Integer{Value: 2},
		&Integer{Value: 3},
	}}
	headRet := _head(arr)
	tailRet := _tail(arr)

	checkInteger(t, headRet, 1)
	checkIntegerArray(t, tailRet, []int64{2, 3})
}

func checkInteger(t *testing.T, obj Object, expect int64) {
	t.Helper()
	intObj, ok := obj.(*Integer)
	if !ok {
		t.Errorf("obj is not an integer object, got=%T", obj)
	}
	if intObj.Value != expect {
		t.Errorf("integer different, expect %d, got=%d", expect, intObj.Value)
	}
}

func checkIntegerArray(t *testing.T, obj Object, expectArr []int64) {
	t.Helper()
	arrObj, ok := obj.(*Array)
	if !ok {
		t.Errorf("obj is not an array object")
	}
	if len(arrObj.Elements) != len(expectArr) {
		t.Errorf("elements length different, expect %d, got=%d (%+v)", len(expectArr), len(arrObj.Elements), expectArr)
	}
	for i, expect := range expectArr {
		intObj, ok := arrObj.Elements[i].(*Integer)
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
