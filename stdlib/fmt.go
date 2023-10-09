package stdlib

import (
	"fmt"
	"io"
	"os"

	"github.com/thingsme/thingscript/object"
)

type Option func(p any)

func FmtPackage(opts ...Option) object.Package {
	ret := &fmtPkg{
		out: os.Stdout,
	}
	for _, o := range opts {
		o(ret)
	}
	return ret
}

func WithWriter(w io.Writer) Option {
	return func(pkg any) {
		if p, ok := pkg.(*fmtPkg); ok {
			p.out = w
		}
	}
}

type fmtPkg struct {
	out io.Writer
}

func (fp *fmtPkg) Name() string {
	return "fmt"
}

func (fp *fmtPkg) Member(member string) func(object.Object, ...object.Object) object.Object {
	switch member {
	case "println":
		return func(receiver object.Object, args ...object.Object) object.Object {
			params := make([]any, len(args))
			for i, a := range args {
				switch raw := a.(type) {
				case *object.String:
					params[i] = raw.Value
				case *object.Integer:
					params[i] = raw.Value
				case *object.Boolean:
					params[i] = raw.Value
				default:
					params[i] = a.Inspect()
				}
			}
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
			params := make([]any, len(args)-1)
			for i, a := range args[1:] {
				switch raw := a.(type) {
				case *object.String:
					params[i] = raw.Value
				case *object.Integer:
					params[i] = raw.Value
				case *object.Boolean:
					params[i] = raw.Value
				default:
					params[i] = a.Inspect()
				}
			}
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
