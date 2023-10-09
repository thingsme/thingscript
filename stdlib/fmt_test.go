package stdlib_test

import (
	"bytes"
	"testing"

	"github.com/thingsme/thingscript/eval"
	"github.com/thingsme/thingscript/lexer"
	"github.com/thingsme/thingscript/object"
	"github.com/thingsme/thingscript/parser"
	"github.com/thingsme/thingscript/stdlib"
)

func TestFmt(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			input: `
				var out = import("fmt")
				out.println(1, 2, "x")
			`,
			expected: "1 2 x\n",
		},
		{
			input: `
				out := import("fmt")
				tv := true
				out.printf("%02d %02x %s", 1, 10, "y")
			`,
			expected: "01 0a y",
			// TODO: losting tv (boolean) values
			//out.printf("%02d %02x %t %s", 1, 10, tv, "y")
		},
	}
	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := parser.New(l)
		program := p.ParseProgram()
		env := object.NewEnvironment()
		out := &bytes.Buffer{}
		env.RegisterPackages(stdlib.FmtPackage(stdlib.WithWriter(out)))

		ret := eval.Eval(program, env)
		if ret != nil && ret.Type() == object.ERROR_OBJ {
			t.Errorf("result is error; %s", ret.Inspect())
		}
		str := out.String()
		if str != tt.expected {
			t.Errorf("result is not %q, got=%q", tt.expected, str)
		}
	}
}
