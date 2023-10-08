package strings

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
	str := &object.String{Value: "1234"}
	sp := New()
	strLen := sp.Member("length")(str)
	checkInteger(t, strLen, 4)
}
