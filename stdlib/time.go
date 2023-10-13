package stdlib

import (
	"fmt"
	"time"

	"github.com/thingsme/thingscript/object"
)

type timePkg struct {
	timeProvider func() time.Time
}

var _ object.Package = &timePkg{}

func (tp *timePkg) Type() object.ObjectType { return object.PACKAGE_OBJ }

func (tp *timePkg) Inspect() string { return "package time" }

func (tp *timePkg) Name() string { return "time" }

func (tp *timePkg) OnLoad(env *object.Environment) {
	if env.TimeProvider != nil {
		tp.timeProvider = env.TimeProvider
	} else {
		tp.timeProvider = func() time.Time { return time.Now() }
	}
}

func (tp *timePkg) Member(name string) object.MemberFunc {
	switch name {
	case "Time":
		return func(receiver object.Object, args ...object.Object) object.Object {
			if len(args) == 1 {
				switch v := args[0].(type) {
				case *TimeObj:
					return &TimeObj{tm: v.tm}
				case *object.Integer:
					return &TimeObj{tm: time.Unix(0, v.Value)}
				default:
					return nil
				}
			}
			return &TimeObj{tm: time.Unix(0, 0)}
		}
	case "Now":
		return func(receiver object.Object, args ...object.Object) object.Object {
			return &TimeObj{tm: tp.timeProvider()}
		}
	default:
		return nil
	}
}

type TimeObj struct {
	tm time.Time
}

var _ object.Object = &TimeObj{}

func (to *TimeObj) Type() object.ObjectType {
	return "time.Time"
}

func (to *TimeObj) Inspect() string {
	return fmt.Sprintf("time.Time(%s)", to.tm)
}

func (to *TimeObj) Member(name string) object.MemberFunc {
	switch name {
	case "=":
		return func(receiver object.Object, args ...object.Object) object.Object {
			if len(args) != 1 {
				return object.Errorf("wrong number of arguments. got=%d, want=1", len(args))
			}
			if left, ok := receiver.(*TimeObj); ok {
				switch v := args[0].(type) {
				case *TimeObj:
					left.tm = v.tm
					return left
				}
			}
			return nil
		}
	default:
		return nil
	}
}
