package stdlib

import (
	"fmt"
	"io"
	"os"

	"github.com/thingsme/thingscript/object"
)

type fmtPkg struct {
	out io.Writer
}

var _ object.PackageImpl = &fmtPkg{}

func (fp *fmtPkg) Name() string { return "fmt" }

func (fp *fmtPkg) OnLoad(env *object.Environment) {
	if env.Stdout != nil {
		fp.out = env.Stdout
	} else {
		fp.out = os.Stdout
	}
}

func (fp *fmtPkg) Member(member string) object.MemberFunc {
	switch member {
	case "println":
		return func(receiver object.Object, args ...object.Object) object.Object {
			params := object2native(args)
			n, err := fmt.Fprintln(fp.out, params...)
			if err != nil {
				return object.Errorf(err.Error())
			}
			return &object.Integer{Value: int64(n)}
		}
	case "printf":
		return func(receiver object.Object, args ...object.Object) object.Object {
			if len(args) == 0 {
				return object.Errorf("wrong number of arguments. got=%d, want >= 1", len(args))
			}
			format := args[0].Inspect()
			params := object2native(args[1:])
			n, err := fmt.Fprintf(fp.out, format, params...)
			if err != nil {
				return object.Errorf(err.Error())
			}
			return &object.Integer{Value: int64(n)}
		}
	default:
		return nil
	}
}

func object2native(args []object.Object) []any {
	params := make([]any, len(args))
	for i, a := range args {
		switch raw := a.(type) {
		case *object.String:
			params[i] = raw.Value
		case *object.Integer:
			params[i] = raw.Value
		case *object.Boolean:
			params[i] = raw.Value
		case *object.Float:
			params[i] = raw.Value
		default:
			params[i] = a.Inspect()
		}
	}
	return params
}
