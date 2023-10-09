package object

import (
	"fmt"
)

var Builtins = []struct {
	Name    string
	Builtin *Builtin
}{
	{"import", &Builtin{Fn: _import}},
}

func newError(format string, args ...any) *Error {
	return &Error{Message: fmt.Sprintf(format, args...)}
}

func _import(args ...Object) Object {
	if len(args) != 1 {
		return newError("wrong number of arguements. got=%d, want=1", len(args))
	}
	name, ok := args[0].(*String)
	if !ok {
		return newError("argument to import must be string, got %s", args[0].Type())
	}
	return &Import{Name: name.Value}
}

const (
	PKG_INTEGER = "$integer"
	PKG_FLOAT   = "$float"
	PKG_STRING  = "$string"
	PKG_BOOLEAN = "$boolean"
	PKG_ARRAY   = "$array"
	PKG_HASHMAP = "$hashmap"
)

func PackageName(obj Object) (string, bool) {
	switch lv := obj.(type) {
	case *Integer:
		return PKG_INTEGER, true
	case *Float:
		return PKG_FLOAT, true
	case *String:
		return PKG_STRING, true
	case *Boolean:
		return PKG_BOOLEAN, true
	case *Array:
		return PKG_ARRAY, true
	case *HashMap:
		return PKG_HASHMAP, true
	case *Import:
		return lv.Name, true
	}
	return "", false
}
