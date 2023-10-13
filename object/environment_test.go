package object_test

import (
	"testing"

	"github.com/thingsme/thingscript/object"
	"github.com/thingsme/thingscript/stdlib"
)

func TestEnclosingEnvironment(t *testing.T) {
	env := object.NewEnvironment()
	env.RegisterPackages()

	env.Set("my_var", &object.Integer{Value: 123})

	newEnv := object.NewEnclosedEnvironment(env)
	obj, ok := newEnv.Get("my_var")
	if !ok {
		t.Errorf("identifier not found")
	}
	ret, ok := obj.(*object.Integer)
	if !ok {
		t.Fatalf("wrong outer environment")
	}
	if ret.Value != 123 {
		t.Fatalf("wrong outer environment value %d, got=%d", 123, ret.Value)
	}
}

func TestBuiltin(t *testing.T) {
	env := object.NewEnvironment()
	env.RegisterPackages(stdlib.Packages()...)

	btImport := env.Builtin("import")
	if btImport == nil {
		t.Fatal("no import function found")
	}
	fmtPkg := btImport.Func(&object.String{Value: "fmt"})
	if fmtPkg.Type() != object.PACKAGE_OBJ {
		t.Fatal("no fmt package found")
	}

	invalidPkg := btImport.Func(&object.String{Value: "something_does_not_exist"})
	if invalidPkg.Type() != object.ERROR_OBJ {
		t.Fatal("import should returns error")
	}

	fn := env.Builtin("int")
	if fn == nil {
		t.Fatal("no import function found")
	}
	obj := fn.Func()
	if obj == nil {
		t.Errorf("no builtin for int")
	}
}

func TestTypes(t *testing.T) {
	env := object.NewEnvironment()
	env.RegisterPackages(stdlib.Packages()...)

	tests := []struct {
		input    string
		init     object.Object
		expected any
	}{
		{"int", nil, 0},
		{"int", &object.Integer{Value: 123}, 123},
	}

	for _, tt := range tests {
		obj := env.Type("", tt.input, tt.init)
		switch v := tt.expected.(type) {
		case int:
			ret, ok := obj.(*object.Integer)
			if !ok {
				t.Fatalf("invalid type function")
			}
			if ret.Value != int64(v) {
				t.Errorf("wrong value %d, got=%d", v, ret.Value)
			}
		default:
			t.Errorf("wrong test type %T", tt.expected)
		}
	}
}

func TestImpor(t *testing.T) {
	env := object.NewEnvironment()
	env.RegisterPackages(stdlib.Packages()...)

	pkg, ok := env.Import("fmt")
	if !ok {
		t.Errorf("import fmt failed")
	}
	if pkg.Name() != "fmt" {
		t.Errorf("wrong pkg %q, got=%q", "fmt", pkg.Name())
	}
}
