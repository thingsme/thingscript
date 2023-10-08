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
