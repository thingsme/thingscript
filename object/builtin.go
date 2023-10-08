package object

import (
	"fmt"
)

var Builtins = []struct {
	Name    string
	Builtin *Builtin
}{
	// global
	{"puts", &Builtin{Fn: _puts}},
	// string
	{"$string_length", &Builtin{Fn: _length}},
	// array
	{"$array_length", &Builtin{Fn: _length}},
	{"$array_head", &Builtin{Fn: _head}},
	{"$array_tail", &Builtin{Fn: _tail}},
	{"$array_init", &Builtin{Fn: _init}},
	{"$array_last", &Builtin{Fn: _last}},
	{"$array_push", &Builtin{Fn: _push}},
	// hash
	{"$hash_length", &Builtin{Fn: _length}},
}

func newError(format string, args ...any) *Error {
	return &Error{Message: fmt.Sprintf(format, args...)}
}

func GetBuiltinByName(name string) *Builtin {
	for _, def := range Builtins {
		if def.Name == name {
			return def.Builtin
		}
	}
	return nil
}

func _length(args ...Object) Object {
	if len(args) != 1 {
		return newError("wrong number of arguments. got=%d, want=1", len(args))
	}
	switch arg := args[0].(type) {
	case *Array:
		return &Integer{Value: int64(len(arg.Elements))}
	case *String:
		return &Integer{Value: int64(len(arg.Value))}
	case *Hash:
		return &Integer{Value: int64(len(arg.Pairs))}
	default:
		return newError("method 'length' not supported, got %s", args[0].Type())
	}
}

func _puts(args ...Object) Object {
	for _, arg := range args {
		fmt.Println(arg.Inspect())
	}
	return nil
}

func _head(args ...Object) Object {
	if len(args) != 1 {
		return newError("wrong number of arguments. got=%d, want=1", len(args))
	}
	if args[0].Type() != ARRAY_OBJ {
		return newError("argument to head must be ARRAY, got %s", args[0].Type())
	}
	arr := args[0].(*Array)
	if len(arr.Elements) > 0 {
		return arr.Elements[0]
	}
	return nil
}

func _tail(args ...Object) Object {
	if len(args) != 1 {
		return newError("wrong number of arguments. got=%d, want=1", len(args))
	}
	if args[0].Type() != ARRAY_OBJ {
		return newError("argument to tail must be ARRAY, got %s", args[0].Type())
	}
	arr := args[0].(*Array)
	length := len(arr.Elements)
	if length > 0 {
		newElements := make([]Object, length-1)
		copy(newElements, arr.Elements[1:length])
		return &Array{Elements: newElements}
	}
	return nil
}

func _init(args ...Object) Object {
	if len(args) != 1 {
		return newError("wrong number of arguments. got=%d, want=1", len(args))
	}
	if args[0].Type() != ARRAY_OBJ {
		return newError("argument to init must be ARRAY, got %s", args[0].Type())
	}
	arr := args[0].(*Array)
	length := len(arr.Elements)
	if length > 0 {
		newElements := make([]Object, length-1)
		copy(newElements, arr.Elements[0:length-1])
		return &Array{Elements: newElements}
	}
	return nil
}

func _last(args ...Object) Object {
	if len(args) != 1 {
		return newError("wrong number of arguments. got=%d, want=1", len(args))
	}
	if args[0].Type() != ARRAY_OBJ {
		return newError("argument to last must be ARRAY, got %s", args[0].Type())
	}
	arr := args[0].(*Array)
	if len(arr.Elements) > 0 {
		return arr.Elements[len(arr.Elements)-1]
	}
	return nil
}

func _push(args ...Object) Object {
	if len(args) != 2 {
		return newError("wrong number of arguments. got=%d, want=2", len(args))
	}
	if args[0].Type() != ARRAY_OBJ {
		return newError("argument to push must be ARRAY, got %s", args[0].Type())
	}
	arr := args[0].(*Array)
	length := len(arr.Elements)

	newElements := make([]Object, length+1)
	copy(newElements, arr.Elements)
	newElements[length] = args[1]
	return &Array{Elements: newElements}
}
