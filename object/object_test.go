package object

import (
	"testing"

	"github.com/thingsme/thingscript/ast"
)

func TestInteger(t *testing.T) {
	obj := &Integer{Value: 123}
	fn := obj.Member("+")
	if fn == nil {
		t.Errorf("error op")
	}
}

func TestFloat(t *testing.T) {
	obj := &Float{Value: 1.234}
	fn := obj.Member("+")
	if fn == nil {
		t.Errorf("error op")
	}
}

func TestBool(t *testing.T) {
	obj := &Boolean{Value: false}
	fn := obj.Member("==")
	if fn == nil {
		t.Errorf("error op")
	}
}

func TestString(t *testing.T) {
	obj := &String{}
	fn := obj.Member("+")
	if fn == nil {
		t.Errorf("error op")
	}
}

func TestArray(t *testing.T) {
	obj := &Array{}
	fn := obj.Member("[")
	if fn == nil {
		t.Errorf("error op")
	}
}

func TestHashMap(t *testing.T) {
	obj := &HashMap{}
	fn := obj.Member("[")
	if fn == nil {
		t.Errorf("error op")
	}
}

func TestInspects(t *testing.T) {
	tests := []struct {
		obj      Object
		typ      ObjectType
		expected string
	}{
		{&Null{}, NULL_OBJ, "null"},
		{&Error{Message: "some error message"}, ERROR_OBJ, "ERROR: some error message"},
		{&Integer{Value: 123}, INTEGER_OBJ, "123"},
		{&Float{Value: 3.14}, FLOAT_OBJ, "3.140000"},
		{&Boolean{Value: true}, BOOLEAN_OBJ, "true"},
		{&String{Value: "text"}, STRING_OBJ, "text"},
		{&ReturnValue{Value: &String{Value: "result"}}, RETURN_VALUE_OBJ, "result"},
		{&Break{}, BREAK_OBJ, "break"},
		{&Function{Parameters: []*ast.Identifier{}, Body: &ast.BlockStatement{}}, FUNCTION_OBJ, "func() {\n}"},
		{&Builtin{}, BUILTIN_OBJ, "builtin"},
		{&Array{Elements: []Object{&Integer{Value: 1}, &Integer{Value: 2}, &Integer{Value: 3}}}, ARRAY_OBJ, "[1, 2, 3]"},
		{&HashMap{Pairs: map[HashKey]HashPair{(&String{Value: "key"}).HashKey(): {Key: &String{Value: "key"}, Value: &String{Value: "value"}}}}, HASHMAP_OBJ, "{key: value}"},
	}

	for _, tt := range tests {
		if tt.obj.Inspect() != tt.expected {
			t.Errorf("wrong inspect %q, got=%q", tt.expected, tt.obj.Inspect())
		}
		if tt.typ != tt.obj.Type() {
			t.Errorf("wrong type %q, got=%q", tt.typ, tt.obj.Type())
		}
	}
}

func TestStringHashKey(t *testing.T) {
	hello1 := &String{Value: "Hello World"}
	hello2 := &String{Value: "Hello World"}
	diff1 := &String{Value: "My name is johnny"}
	diff2 := &String{Value: "My name is johnny"}

	if hello1.HashKey() != hello2.HashKey() {
		t.Errorf("strings with same content have different hash keys")
	}

	if diff1.HashKey() != diff2.HashKey() {
		t.Errorf("strings with same content have different hash keys")
	}

	if hello1.HashKey() == diff1.HashKey() {
		t.Errorf("strings with different content have same hash keys")
	}
}

func TestIntegerHashKey(t *testing.T) {
	origin1 := &Integer{Value: 1234567}
	origin2 := &Integer{Value: 1234567}
	diff1 := &Integer{Value: 3141592}
	diff2 := &Integer{Value: 3141592}

	if origin1.HashKey() != origin2.HashKey() {
		t.Errorf("integers with same value have different hash keys")
	}

	if diff1.HashKey() != diff2.HashKey() {
		t.Errorf("integers with same value have different hash keys")
	}

	if origin1.HashKey() == diff1.HashKey() {
		t.Errorf("integers with different value have same hash keys")
	}
}

func TestFloatHashKey(t *testing.T) {
	origin1 := &Float{Value: 1.234567}
	origin2 := &Float{Value: 1.234567}
	diff1 := &Float{Value: 3.141592}
	diff2 := &Float{Value: 3.141592}

	if origin1.HashKey() != origin2.HashKey() {
		t.Errorf("floats with same value have different hash keys")
	}

	if diff1.HashKey() != diff2.HashKey() {
		t.Errorf("floats with same value have different hash keys")
	}

	if origin1.HashKey() == diff1.HashKey() {
		t.Errorf("floats with different value have same hash keys")
	}
}

func TestBooleanHashKey(t *testing.T) {
	origin1 := &Boolean{Value: true}
	origin2 := &Boolean{Value: true}
	diff1 := &Boolean{Value: false}
	diff2 := &Boolean{Value: false}

	if origin1.HashKey() != origin2.HashKey() {
		t.Errorf("boolean with same value have different hash keys")
	}

	if diff1.HashKey() != diff2.HashKey() {
		t.Errorf("boolean with same value have different hash keys")
	}

	if origin1.HashKey() == diff1.HashKey() {
		t.Errorf("boolean with different value have same hash keys")
	}
}
