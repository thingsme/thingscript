package stdlib

import (
	"github.com/thingsme/thingscript/object"
	"github.com/thingsme/thingscript/stdlib/fmt"
)

func Packages() []object.Package {
	return []object.Package{
		fmt.New(),
		&integers{},
		&floats{},
		&booleans{},
		&strings{},
		&arrays{},
		&hashmap{},
	}
}

type integers struct {
}

func (ip *integers) Name() string { return object.PKG_INTEGER }
func (ip *integers) Member(member string) func(object.Object, ...object.Object) object.Object {
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

type floats struct {
}

func (fp *floats) Name() string { return object.PKG_FLOAT }
func (fp *floats) Member(member string) func(object.Object, ...object.Object) object.Object {
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

type booleans struct {
}

func (fp *booleans) Name() string { return object.PKG_BOOLEAN }
func (fp *booleans) Member(member string) func(object.Object, ...object.Object) object.Object {
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

type strings struct {
}

func (sp *strings) Name() string { return object.PKG_STRING }
func (sp *strings) Member(member string) func(object.Object, ...object.Object) object.Object {
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

type hashmap struct {
}

func (sp *hashmap) Name() string { return object.PKG_HASHMAP }

func (hp *hashmap) Member(member string) func(object.Object, ...object.Object) object.Object {
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

type arrays struct {
}

func (sp *arrays) Name() string { return object.PKG_ARRAY }
func (sp *arrays) Member(member string) func(object.Object, ...object.Object) object.Object {
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
