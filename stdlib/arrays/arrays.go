package arrays

import "github.com/thingsme/thingscript/object"

func New() object.Package {
	return &pkg{}
}

type pkg struct {
}

func (sp *pkg) Name() string { return "$array" }
func (sp *pkg) Member(member string) func(object.Object, ...object.Object) object.Object {
	switch member {
	case "length":
		return func(receiver object.Object, args ...object.Object) object.Object {
			if len(args) != 0 {
				return object.Errorf("wrong number of arguments. got=%d, want=0", len(args))
			}
			arr := receiver.(*object.Array)
			return &object.Integer{Value: int64(len(arr.Elements))}
		}
	case "head":
		return func(receiver object.Object, args ...object.Object) object.Object {
			if len(args) != 0 {
				return object.Errorf("wrong number of arguments. got=%d, want=0", len(args))
			}
			arr := receiver.(*object.Array)
			if len(arr.Elements) > 0 {
				return arr.Elements[0]
			}
			return nil
		}
	case "tail":
		return func(receiver object.Object, args ...object.Object) object.Object {
			if len(args) != 0 {
				return object.Errorf("wrong number of arguments. got=%d, want=0", len(args))
			}
			arr := receiver.(*object.Array)
			length := len(arr.Elements)
			if length > 0 {
				newElements := make([]object.Object, length-1)
				copy(newElements, arr.Elements[1:length])
				return &object.Array{Elements: newElements}
			}
			return nil
		}
	case "init":
		return func(receiver object.Object, args ...object.Object) object.Object {
			if len(args) != 0 {
				return object.Errorf("wrong number of arguments. got=%d, want=0", len(args))
			}
			arr := receiver.(*object.Array)
			length := len(arr.Elements)
			if length > 0 {
				newElements := make([]object.Object, length-1)
				copy(newElements, arr.Elements[0:length-1])
				return &object.Array{Elements: newElements}
			}
			return nil
		}
	case "last":
		return func(receiver object.Object, args ...object.Object) object.Object {
			if len(args) != 0 {
				return object.Errorf("wrong number of arguments. got=%d, want=0", len(args))
			}
			arr := receiver.(*object.Array)
			if len(arr.Elements) > 0 {
				return arr.Elements[len(arr.Elements)-1]
			}
			return nil
		}
	case "push":
		return func(receiver object.Object, args ...object.Object) object.Object {
			if len(args) != 1 {
				return object.Errorf("wrong number of arguments. got=%d, want=1", len(args))
			}
			arr := receiver.(*object.Array)
			length := len(arr.Elements)
			newElements := make([]object.Object, length+1)
			copy(newElements, arr.Elements)
			newElements[length] = args[0]
			return &object.Array{Elements: newElements}
		}
	default:
		return nil
	}
}
