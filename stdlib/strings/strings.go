package strings

import "github.com/thingsme/thingscript/object"

func New() object.Package {
	return &pkg{}
}

type pkg struct {
}

func (sp *pkg) Name() string { return "$string" }
func (sp *pkg) Member(member string) func(object.Object, ...object.Object) object.Object {
	switch member {
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
