package stdlib

import (
	"github.com/thingsme/thingscript/eval"
	"github.com/thingsme/thingscript/object"
)

func Packages() []object.PackageImpl {
	return []object.PackageImpl{
		&primitives{},
		&fmtPkg{},
		&timePkg{},
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

type primitives struct {
}

var _ object.PackageImpl = &primitives{}

func (p *primitives) Name() string { return "" }

func (p *primitives) OnLoad(env *object.Environment) {}

func (p *primitives) Member(name string) object.MemberFunc {
	switch name {
	case "int":
		return func(receiver object.Object, args ...object.Object) object.Object {
			if len(args) == 1 {
				switch v := args[0].(type) {
				case *object.Integer:
					return &object.Integer{Value: v.Value}
				}
			}
			return &object.Integer{Value: 0}
		}
	case "float":
		return func(receiver object.Object, args ...object.Object) object.Object {
			if len(args) == 1 {
				switch v := args[0].(type) {
				case *object.Float:
					return &object.Float{Value: v.Value}
				case *object.Integer:
					return &object.Float{Value: float64(v.Value)}
				}
			}
			return &object.Float{Value: 0}
		}
	case "string":
		return func(receiver object.Object, args ...object.Object) object.Object {
			if len(args) == 1 {
				switch v := args[0].(type) {
				case *object.String:
					return &object.String{Value: v.Value}
				}
			}
			return &object.String{Value: ""}
		}
	case "bool":
		return func(receiver object.Object, args ...object.Object) object.Object {
			if len(args) == 1 {
				switch v := args[0].(type) {
				case *object.Boolean:
					return &object.Boolean{Value: v.Value}
				}
			}
			return &object.Boolean{Value: false}
		}
	case "array":
		return func(receiver object.Object, args ...object.Object) object.Object {
			if len(args) > 0 {
				return &object.Array{Elements: args}
			}
			return &object.Array{Elements: []object.Object{}}
		}
	}
	return nil
}

func errWrongNumberOfArguments(want int, got int) *object.Error {
	return object.Errorf("wrong number of arguments. want=%d got=%d", want, got)
}

func errTypeMismatched(left object.Object, oper string, right object.Object) *object.Error {
	return object.Errorf("type mismatch: %s %s %s", left.Type(), oper, right.Type())
}

func errUnknownOperator(left object.Object, oper string) *object.Error {
	return object.Errorf("unknown operator %s of %s", oper, left.Type())
}

func Integers(member string) object.MemberFunc {
	switch member {
	case "type":
		return func(receiver object.Object, args ...object.Object) object.Object {
			if len(args) != 0 {
				return errWrongNumberOfArguments(0, len(args))
			}
			return &object.String{Value: "int"}
		}
	case "=", "+", "-", "*", "/", "%", "<", "<=", ">", ">=", "==", "!=":
		return func(receiver object.Object, args ...object.Object) object.Object {
			if len(args) != 1 {
				return errWrongNumberOfArguments(1, len(args))
			}
			left, ok := receiver.(*object.Integer)
			if !ok {
				return nil
			}
			switch right := args[0].(type) {
			case *object.Integer:
				switch member {
				case "=":
					left.Value = right.Value
					return left
				case "+":
					return &object.Integer{Value: left.Value + right.Value}
				case "-":
					return &object.Integer{Value: left.Value - right.Value}
				case "*":
					return &object.Integer{Value: left.Value * right.Value}
				case "/":
					return &object.Integer{Value: left.Value / right.Value}
				case "%":
					return &object.Integer{Value: left.Value % right.Value}
				case "<":
					return &object.Boolean{Value: left.Value < right.Value}
				case "<=":
					return &object.Boolean{Value: left.Value <= right.Value}
				case ">":
					return &object.Boolean{Value: left.Value > right.Value}
				case ">=":
					return &object.Boolean{Value: left.Value >= right.Value}
				case "==":
					return &object.Boolean{Value: left.Value == right.Value}
				case "!=":
					return &object.Boolean{Value: left.Value != right.Value}
				default:
					return errTypeMismatched(receiver, member, args[0])
				}
			case *object.Float:
				switch member {
				case "+":
					return &object.Float{Value: float64(left.Value) + right.Value}
				case "-":
					return &object.Float{Value: float64(left.Value) - right.Value}
				case "*":
					return &object.Float{Value: float64(left.Value) * right.Value}
				case "/":
					return &object.Float{Value: float64(left.Value) / right.Value}
				case "<":
					return &object.Boolean{Value: float64(left.Value) < right.Value}
				case "<=":
					return &object.Boolean{Value: float64(left.Value) <= right.Value}
				case ">":
					return &object.Boolean{Value: float64(left.Value) > right.Value}
				case ">=":
					return &object.Boolean{Value: float64(left.Value) >= right.Value}
				case "==":
					return &object.Boolean{Value: float64(left.Value) == right.Value}
				case "!=":
					return &object.Boolean{Value: float64(left.Value) != right.Value}
				default:
					return errTypeMismatched(receiver, member, args[0])
				}
			default:
				return errTypeMismatched(receiver, member, args[0])
			}
		}
	}
	return nil
}

func Floats(member string) object.MemberFunc {
	switch member {
	case "type":
		return func(receiver object.Object, args ...object.Object) object.Object {
			if len(args) != 0 {
				return errWrongNumberOfArguments(0, len(args))
			}
			return &object.String{Value: "float"}
		}
	case "=", "+", "-", "*", "/", "<", "<=", ">", ">=", "==", "!=":
		return func(receiver object.Object, args ...object.Object) object.Object {
			if len(args) != 1 {
				return errWrongNumberOfArguments(1, len(args))
			}
			left, ok := receiver.(*object.Float)
			if !ok {
				return nil
			}
			var rightValue float64
			switch rv := args[0].(type) {
			case *object.Integer:
				rightValue = float64(rv.Value)
			case *object.Float:
				rightValue = rv.Value
			default:
				return errTypeMismatched(receiver, member, args[0])
			}
			switch member {
			case "=":
				left.Value = rightValue
				return left
			case "+":
				return &object.Float{Value: left.Value + rightValue}
			case "-":
				return &object.Float{Value: left.Value - rightValue}
			case "*":
				return &object.Float{Value: left.Value * rightValue}
			case "/":
				return &object.Float{Value: left.Value / rightValue}
			case "<":
				return &object.Boolean{Value: left.Value < rightValue}
			case "<=":
				return &object.Boolean{Value: left.Value <= rightValue}
			case ">":
				return &object.Boolean{Value: left.Value > rightValue}
			case ">=":
				return &object.Boolean{Value: left.Value >= rightValue}
			case "==":
				return &object.Boolean{Value: left.Value == rightValue}
			case "!=":
				return &object.Boolean{Value: left.Value != rightValue}
			default:
				return errUnknownOperator(receiver, member)
			}
		}
	}
	return nil
}

func Booleans(member string) object.MemberFunc {
	switch member {
	case "type":
		return func(receiver object.Object, args ...object.Object) object.Object {
			if len(args) != 0 {
				return errWrongNumberOfArguments(0, len(args))
			}
			return &object.String{Value: "bool"}
		}
	case "=", "==", "!=":
		return func(receiver object.Object, args ...object.Object) object.Object {
			if len(args) != 1 {
				return errWrongNumberOfArguments(1, len(args))
			}
			left, ok := receiver.(*object.Boolean)
			if !ok {
				return nil
			}
			switch right := args[0].(type) {
			case *object.Boolean:
				switch member {
				case "=":
					left.Value = right.Value
					return left
				case "==":
					return &object.Boolean{Value: left.Value == right.Value}
				case "!=":
					return &object.Boolean{Value: left.Value != right.Value}
				}
			default:
				return errTypeMismatched(receiver, member, args[0])
			}
			return errUnknownOperator(receiver, member)
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
				return errWrongNumberOfArguments(0, len(args))
			}
			return &object.String{Value: "string"}
		}
	case "=", "+", "<", "<=", ">", ">=", "==", "!=":
		return func(receiver object.Object, args ...object.Object) object.Object {
			if len(args) != 1 {
				return errWrongNumberOfArguments(1, len(args))
			}
			left, ok := receiver.(*object.String)
			if !ok {
				return nil
			}
			switch right := args[0].(type) {
			case *object.String:
				switch member {
				case "=":
					left.Value = right.Value
					return left
				case "+":
					return &object.String{Value: left.Value + right.Value}
				case "<":
					return &object.Boolean{Value: left.Value < right.Value}
				case "<=":
					return &object.Boolean{Value: left.Value <= right.Value}
				case ">":
					return &object.Boolean{Value: left.Value > right.Value}
				case ">=":
					return &object.Boolean{Value: left.Value >= right.Value}
				case "==":
					return &object.Boolean{Value: left.Value == right.Value}
				case "!=":
					return &object.Boolean{Value: left.Value != right.Value}
				}
			default:
				return errTypeMismatched(receiver, member, args[0])
			}
			return errUnknownOperator(receiver, member)
		}
	case "length":
		return func(receiver object.Object, args ...object.Object) object.Object {
			if len(args) != 0 {
				return errWrongNumberOfArguments(0, len(args))
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
				return errWrongNumberOfArguments(0, len(args))
			}
			return &object.String{Value: "map"}
		}
	case "=":
		return func(receiver object.Object, args ...object.Object) object.Object {
			if len(args) != 1 {
				return errWrongNumberOfArguments(1, len(args))
			}
			if lv, ok := receiver.(*object.HashMap); ok {
				switch rv := args[0].(type) {
				case *object.HashMap:
					lv.Pairs = rv.Pairs
					return lv
				}
			}
			return nil
		}
	case "[": // index oper
		return func(receiver object.Object, args ...object.Object) object.Object {
			if len(args) != 1 {
				return errWrongNumberOfArguments(1, len(args))
			}
			h := receiver.(*object.HashMap)
			key, ok := args[0].(object.Hashable)
			if !ok {
				return object.Errorf("unusable as hash key: %s", args[0].Type())
			}
			pair, ok := h.Pairs[key.HashKey()]
			if !ok {
				return eval.NULL
			}
			return pair.Value
		}
	case "length":
		return func(receiver object.Object, args ...object.Object) object.Object {
			if len(args) != 0 {
				return errWrongNumberOfArguments(0, len(args))
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
				return errWrongNumberOfArguments(0, len(args))
			}
			return &object.String{Value: "array"}
		}
	case "=":
		return func(receiver object.Object, args ...object.Object) object.Object {
			if len(args) != 1 {
				return errWrongNumberOfArguments(1, len(args))
			}
			if lv, ok := receiver.(*object.Array); ok {
				switch rv := args[0].(type) {
				case *object.Array:
					lv.Elements = rv.Elements
					return lv
				}
			}
			return nil
		}
	case "[": // index oper
		return func(receiver object.Object, args ...object.Object) object.Object {
			if len(args) != 1 {
				return errWrongNumberOfArguments(1, len(args))
			}
			arr := receiver.(*object.Array)
			max := int64(len(arr.Elements) - 1)
			if rv, ok := args[0].(*object.Integer); ok {
				idx := rv.Value
				if idx < 0 || idx > max {
					return eval.NULL
				}
				return arr.Elements[idx]
			}
			return nil
		}
	case "length":
		return func(receiver object.Object, args ...object.Object) object.Object {
			if len(args) != 0 {
				return errWrongNumberOfArguments(0, len(args))
			}
			arr := receiver.(*object.Array)
			return &object.Integer{Value: int64(len(arr.Elements))}
		}
	case "head":
		return func(receiver object.Object, args ...object.Object) object.Object {
			if len(args) != 0 {
				return errWrongNumberOfArguments(0, len(args))
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
				return errWrongNumberOfArguments(0, len(args))
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
				return errWrongNumberOfArguments(0, len(args))
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
				return errWrongNumberOfArguments(0, len(args))
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
				return errWrongNumberOfArguments(1, len(args))
			}
			fn := args[0].(*object.Function)
			if len(fn.Parameters) != 2 {
				return errWrongNumberOfArguments(2, len(args))
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
				return errWrongNumberOfArguments(1, len(args))
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
