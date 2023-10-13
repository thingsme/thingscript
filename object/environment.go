package object

import (
	"io"
	"time"
)

type Environment struct {
	outer    *Environment
	store    map[string]Object
	packages map[string]PackageImpl

	Stdout       io.Writer
	TimeProvider func() time.Time
}

func NewEnvironment() *Environment {
	env := &Environment{
		store:    make(map[string]Object),
		packages: make(map[string]PackageImpl),
	}
	return env
}

func NewEnclosedEnvironment(outer *Environment) *Environment {
	env := NewEnvironment()
	env.outer = outer
	return env
}

func (e *Environment) Builtin(name string) *Builtin {
	switch name {
	case "import":
		return &Builtin{Func: func(args ...Object) Object {
			if len(args) != 1 {
				return Errorf("wrong number of arguements. got=%d, want=1", len(args))
			}
			name, ok := args[0].(*String)
			if !ok {
				return Errorf("argument to import must be string, got %s", args[0].Type())
			}
			if pkg, ok := e.packages[name.Value]; ok {
				return &Package{pkg: pkg}
			} else {
				return Errorf("package %q not found", name.Value)
			}
		}}
	default:
		pkg, ok := e.Import("")
		if !ok {
			return nil
		}
		memberFunc := pkg.Member(name)
		if memberFunc == nil {
			return nil
		}
		return &Builtin{Func: func(args ...Object) Object {
			return memberFunc(nil, args...)
		}}
	}
}

func (e *Environment) Type(pkgName string, name string, initial Object) Object {
	pkg, ok := e.packages[pkgName]
	if !ok {
		return Errorf("unknown %q", pkgName)
	}
	memberFunc := pkg.Member(name)
	if memberFunc == nil {
		return Errorf("unknown %s.%s", pkgName, name)
	}
	var ret Object
	if initial != nil {
		ret = memberFunc(nil, initial)
		if ret == nil {
			return Errorf("unknown %s.%s(%s)", pkgName, name, initial.Type())
		}
	} else {
		ret = memberFunc(nil)
		if ret == nil {
			return Errorf("unknown %s.%s()", pkgName, name)
		}
	}
	return ret
}

func (e *Environment) Get(name string) (Object, bool) {
	obj, ok := e.store[name]
	if !ok && e.outer != nil {
		obj, ok = e.outer.Get(name)
	}
	if !ok {
		if pkg, ok := e.packages[name]; ok {
			obj = &Package{pkg: pkg}
		}
	}
	return obj, ok
}

func (e *Environment) Set(name string, val Object) Object {
	e.store[name] = val
	return val
}

func (e *Environment) RegisterPackages(pkgs ...PackageImpl) {
	for _, p := range pkgs {
		p.OnLoad(e)
		e.packages[p.Name()] = p
	}
}

func (e *Environment) Import(name string) (PackageImpl, bool) {
	p, ok := e.packages[name]
	if !ok && e.outer != nil {
		p, ok = e.outer.Import(name)
	}
	return p, ok
}
