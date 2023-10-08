package hashmap

import "github.com/thingsme/thingscript/object"

func New() object.Package {
	return &pkg{}
}

type pkg struct {
}

func (sp *pkg) Name() string { return "$hash" }

func (hp *pkg) Member(member string) func(object.Object, ...object.Object) object.Object {
	switch member {
	case "length":
		return func(receiver object.Object, args ...object.Object) object.Object {
			if len(args) != 0 {
				return object.Errorf("wrong number of arguments. got=%d, want=0", len(args))
			}
			h := receiver.(*object.Hash)
			return &object.Integer{Value: int64(len(h.Pairs))}
		}
	default:
		return nil
	}
}
