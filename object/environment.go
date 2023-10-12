package object

import "io"

type Environment struct {
	Stdout io.Writer

	store    map[string]Object
	outer    *Environment
	packages map[string]PackageImpl
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
		return nil
	}
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
