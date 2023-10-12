package object

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"hash/fnv"
	"math"
	"strings"

	"github.com/thingsme/thingscript/ast"
)

type ObjectType string

const (
	INTEGER_OBJ      = "INTEGER"
	FLOAT_OBJ        = "FLOAT"
	BOOLEAN_OBJ      = "BOOLEAN"
	STRING_OBJ       = "STRING"
	ARRAY_OBJ        = "ARRAY"
	HASHMAP_OBJ      = "HASHMAP"
	FUNCTION_OBJ     = "FUNCTION"
	BUILTIN_OBJ      = "BUILTIN"
	PACKAGE_OBJ      = "PACKAGE"
	NULL_OBJ         = "NULL"
	RETURN_VALUE_OBJ = "RETURN_VALUE"
	BREAK_OBJ        = "BREAK"
	ERROR_OBJ        = "ERROR"
)

type MemberFunc func(receiver Object, args ...Object) Object
type FunctionFunc func(args ...Object) Object

type Object interface {
	Type() ObjectType
	Member(name string) MemberFunc
	Inspect() string
}

type Null struct{}

func (n *Null) Type() ObjectType              { return NULL_OBJ }
func (n *Null) Inspect() string               { return "null" }
func (n *Null) Member(name string) MemberFunc { return nil }

type Error struct {
	Message string
}

func (e *Error) Type() ObjectType              { return ERROR_OBJ }
func (e *Error) Inspect() string               { return "ERROR: " + e.Message }
func (e *Error) Member(name string) MemberFunc { return nil }

func Errorf(format string, args ...any) *Error {
	return &Error{Message: fmt.Sprintf(format, args...)}
}

type Integer struct {
	Value int64
}

func (i *Integer) Type() ObjectType { return INTEGER_OBJ }
func (i *Integer) Inspect() string  { return fmt.Sprintf("%d", i.Value) }
func (i *Integer) Member(name string) MemberFunc {
	if IntegerMemberFunc != nil {
		return IntegerMemberFunc(name)
	}
	return nil
}

var IntegerMemberFunc func(string) MemberFunc

type Float struct {
	Value float64
}

func (f *Float) Type() ObjectType { return FLOAT_OBJ }
func (f *Float) Inspect() string  { return fmt.Sprintf("%f", f.Value) }
func (f *Float) Member(name string) MemberFunc {
	if FloatMemberFunc != nil {
		return FloatMemberFunc(name)
	}
	return nil
}

var FloatMemberFunc func(string) MemberFunc

type Boolean struct {
	Value bool
}

func (b *Boolean) Type() ObjectType { return BOOLEAN_OBJ }
func (b *Boolean) Inspect() string  { return fmt.Sprintf("%t", b.Value) }
func (b *Boolean) Member(name string) MemberFunc {
	if BooleanMemberFunc != nil {
		return BooleanMemberFunc(name)
	}
	return nil
}

var BooleanMemberFunc func(string) MemberFunc

type String struct {
	Value string
}

func (s *String) Type() ObjectType { return STRING_OBJ }
func (s *String) Inspect() string  { return s.Value }
func (b *String) Member(name string) MemberFunc {
	if StringMemberFunc != nil {
		return StringMemberFunc(name)
	}
	return nil
}

var StringMemberFunc func(string) MemberFunc

type ReturnValue struct {
	Value Object
}

func (rv *ReturnValue) Type() ObjectType              { return RETURN_VALUE_OBJ }
func (rv *ReturnValue) Inspect() string               { return rv.Value.Inspect() }
func (rv *ReturnValue) Member(name string) MemberFunc { return nil }

type Break struct {
}

func (br *Break) Type() ObjectType              { return BREAK_OBJ }
func (br *Break) Inspect() string               { return "break" }
func (br *Break) Member(name string) MemberFunc { return nil }

type Function struct {
	Parameters []*ast.Identifier
	Body       *ast.BlockStatement
	Env        *Environment
}

func (f *Function) Type() ObjectType { return FUNCTION_OBJ }
func (f *Function) Inspect() string {
	var out bytes.Buffer
	params := []string{}
	for _, p := range f.Parameters {
		params = append(params, p.String())
	}
	out.WriteString("func")
	out.WriteString("(")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(") {\n")
	out.WriteString(f.Body.String())
	out.WriteString("}")
	return out.String()
}
func (br *Function) Member(name string) MemberFunc { return nil }

type Builtin struct {
	Func FunctionFunc
}

func (b *Builtin) Type() ObjectType              { return BUILTIN_OBJ }
func (b *Builtin) Inspect() string               { return "builtin" }
func (b *Builtin) Member(name string) MemberFunc { return nil }

type Array struct {
	Elements []Object
}

func (ao *Array) Type() ObjectType { return ARRAY_OBJ }
func (ao *Array) Inspect() string {
	var out bytes.Buffer
	elements := []string{}
	for _, e := range ao.Elements {
		elements = append(elements, e.Inspect())
	}
	out.WriteString("[")
	out.WriteString(strings.Join(elements, ", "))
	out.WriteString("]")
	return out.String()
}

func (ao *Array) Member(name string) MemberFunc {
	if ArrayMemberFunc != nil {
		return ArrayMemberFunc(name)
	}
	return nil
}

var ArrayMemberFunc func(name string) MemberFunc

type HashKey struct {
	Type  ObjectType
	Value uint64
}

type Hashable interface {
	HashKey() HashKey
}

var _ Hashable = &Integer{}
var _ Hashable = &Float{}
var _ Hashable = &Boolean{}
var _ Hashable = &String{}

func (b *Boolean) HashKey() HashKey {
	var value uint64
	if b.Value {
		value = 1
	} else {
		value = 0
	}
	return HashKey{Type: b.Type(), Value: value}
}

func (i *Integer) HashKey() HashKey {
	return HashKey{Type: i.Type(), Value: uint64(i.Value)}
}

func (f *Float) HashKey() HashKey {
	h := fnv.New64a()
	var buf [8]byte
	binary.BigEndian.PutUint64(buf[:], math.Float64bits(f.Value))
	h.Write(buf[:])
	return HashKey{Type: f.Type(), Value: h.Sum64()}
}

func (s *String) HashKey() HashKey {
	h := fnv.New64a()
	h.Write([]byte(s.Value))
	return HashKey{Type: s.Type(), Value: h.Sum64()}
}

type HashPair struct {
	Key   Object
	Value Object
}

type HashMap struct {
	Pairs map[HashKey]HashPair
}

func (h *HashMap) Type() ObjectType { return HASHMAP_OBJ }
func (h *HashMap) Inspect() string {
	var out bytes.Buffer
	pairs := []string{}
	for _, pair := range h.Pairs {
		pairs = append(pairs, fmt.Sprintf("%s: %s",
			pair.Key.Inspect(), pair.Value.Inspect()))
	}
	out.WriteString("{")
	out.WriteString(strings.Join(pairs, ", "))
	out.WriteString("}")
	return out.String()
}

func (h *HashMap) Member(name string) MemberFunc {
	if HashMapMemberFunc != nil {
		return HashMapMemberFunc(name)
	}
	return nil
}

var HashMapMemberFunc func(name string) MemberFunc

type Package struct {
	pkg PackageImpl
}

type PackageImpl interface {
	Name() string
	Member(name string) MemberFunc
	OnLoad(*Environment)
}

func (p *Package) Type() ObjectType { return PACKAGE_OBJ }
func (p *Package) Inspect() string {
	if p.pkg == nil {
		return "import()"
	}
	return fmt.Sprintf("import(%q)", p.pkg.Name())
}
func (p *Package) Member(name string) MemberFunc {
	if p.pkg == nil {
		return nil
	}
	return p.pkg.Member(name)
}
