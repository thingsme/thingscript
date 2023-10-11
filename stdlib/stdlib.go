package stdlib

import (
	"github.com/thingsme/thingscript/eval"
	"github.com/thingsme/thingscript/object"
)

func Packages() []object.Package {
	return []object.Package{
		FmtPackage(),
	}
}

func init() {
	object.StringMemberFunc = Strings
	object.IntegerMemberFunc = Integers
	object.FloatMemberFunc = Floats
	object.BooleanMemberFunc = Booleans
	object.ArrayMemberFunc = Arrays
	object.HashMapMemberFunc = HashMaps
}

func Integers(member string) object.MemberFunc {
	switch member {
	case "type":
		return func(receiver object.Object, args ...object.Object) object.Object {
			if len(args) != 0 {
				return object.Errorf("wrong number of arguments. got=%d, want=0", len(args))
			}
			return &object.String{Value: "integer"}
		}
	default:
		return nil
	}
}

func Floats(member string) object.MemberFunc {
	switch member {
	case "type":
		return func(receiver object.Object, args ...object.Object) object.Object {
			if len(args) != 0 {
				return object.Errorf("wrong number of arguments. got=%d, want=0", len(args))
			}
			return &object.String{Value: "float"}
		}
	default:
		return nil
	}
}

func Booleans(member string) object.MemberFunc {
	switch member {
	case "type":
		return func(receiver object.Object, args ...object.Object) object.Object {
			if len(args) != 0 {
				return object.Errorf("wrong number of arguments. got=%d, want=0", len(args))
			}
			return &object.String{Value: "boolean"}
		}
	default:
		return nil
	}
}

func Strings(member string) object.MemberFunc {
	switch member {
	case "type":
		return func(receiver object.Object, args ...object.Object) object.Object {
			if len(args) != 0 {
				return object.Errorf("wrong number of arguments. got=%d, want=0", len(args))
			}
			return &object.String{Value: "string"}
		}
	case "length":
		return func(receiver object.Object, args ...object.Object) object.Object {
			if len(args) != 0 {
				return object.Errorf("wrong number of arguments. got=%d, want=0", len(args))
			}
			str := receiver.(*object.String)
			return &object.Integer{Value: int64(len(str.Value))}
		}
	default:
		return nil
	}
}

func HashMaps(member string) object.MemberFunc {
	switch member {
	case "type":
		return func(receiver object.Object, args ...object.Object) object.Object {
			if len(args) != 0 {
				return object.Errorf("wrong number of arguments. got=%d, want=0", len(args))
			}
			return &object.String{Value: "hashmap"}
		}
	case "length":
		return func(receiver object.Object, args ...object.Object) object.Object {
			if len(args) != 0 {
				return object.Errorf("wrong number of arguments. got=%d, want=0", len(args))
			}
			h := receiver.(*object.HashMap)
			return &object.Integer{Value: int64(len(h.Pairs))}
		}
	default:
		return nil
	}
}

func Arrays(member string) object.MemberFunc {
	switch member {
	case "type":
		return func(receiver object.Object, args ...object.Object) object.Object {
			if len(args) != 0 {
				return object.Errorf("wrong number of arguments. got=%d, want=0", len(args))
			}
			return &object.String{Value: "array"}
		}
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
	case "foreach":
		return func(receiver object.Object, args ...object.Object) object.Object {
			if len(args) != 1 {
				return object.Errorf("wrong number of arguments. got=%d, want=1", len(args))
			}
			fn := args[0].(*object.Function)
			if len(fn.Parameters) != 2 {
				return object.Errorf("wrong number of arguments. got=%d, want=2", len(fn.Parameters))
			}
			arr := receiver.(*object.Array)
			for i, elm := range arr.Elements {
				env := object.NewEnclosedEnvironment(fn.Env)
				env.Set(fn.Parameters[0].Value, &object.Integer{Value: int64(i)})
				env.Set(fn.Parameters[1].Value, elm)
				ret := eval.Eval(fn.Body, env)
				if ret != nil {
					if ret.Type() == object.ERROR_OBJ {
						return ret
					} else if ret.Type() == object.BREAK_OBJ {
						break
					}
				}
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
