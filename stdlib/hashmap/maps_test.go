package hashmap

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

func TestLen(t *testing.T) {
	hash := &object.Hash{
		Pairs: map[object.HashKey]object.HashPair{
			(&object.Integer{Value: 1}).HashKey(): {Key: &object.Integer{Value: 1}, Value: &object.Integer{Value: 2}},
		}}
	ap := New()
	mapLen := ap.Member("length")(hash)
	checkInteger(t, mapLen, 1)
}
