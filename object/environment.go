package object

type Environment struct {
	store    map[string]Object
	outer    *Environment
	packages map[string]Package
}

func NewEnvironment() *Environment {
	env := &Environment{
		store:    make(map[string]Object),
		packages: make(map[string]Package),
	}
	return env
}

func NewEnclosedEnvironment(outer *Environment) *Environment {
	env := NewEnvironment()
	env.outer = outer
	return env
}

func (e *Environment) Get(name string) (Object, bool) {
	obj, ok := e.store[name]
	if !ok && e.outer != nil {
		obj, ok = e.outer.Get(name)
	}
	return obj, ok
}

func (e *Environment) Set(name string, val Object) Object {
	e.store[name] = val
	return val
}

func (e *Environment) RegisterPackages(pkgs ...Package) {
	for _, p := range pkgs {
		e.packages[p.Name()] = p
	}
}

func (e *Environment) Import(name string) (Package, bool) {
	p, ok := e.packages[name]
	if !ok && e.outer != nil {
		p, ok = e.outer.Import(name)
	}
	return p, ok
}

type Package interface {
	Name() string
	Member(string) func(receiver Object, args ...Object) Object
}
