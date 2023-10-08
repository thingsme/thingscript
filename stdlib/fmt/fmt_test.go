package fmt

import (
	"bytes"
	"testing"

	"github.com/thingsme/thingscript/eval"
	"github.com/thingsme/thingscript/lexer"
	"github.com/thingsme/thingscript/object"
	"github.com/thingsme/thingscript/parser"
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
				out.printf("%02d %02x %t %s", 1, 10, true, "x")
			`,
			expected: "01 0a true x",
		},
	}
	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := parser.New(l)
		program := p.ParseProgram()
		env := object.NewEnvironment()
		out := &bytes.Buffer{}
		env.RegisterPackages(New(WithWriter(out)))

		ret := eval.Eval(program, env)
		if ret != nil && ret.Type() == object.ERROR_OBJ {
			t.Errorf("result is error; %s", ret.Inspect())
		}
		if out.String() != tt.expected {
			t.Errorf("result is not %q, got=%q", tt.expected, out.String())
		}
	}
}
