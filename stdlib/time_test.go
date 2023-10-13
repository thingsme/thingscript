package stdlib_test

import (
	"bytes"
	"fmt"
	"strconv"
	"testing"
	"time"

	"github.com/thingsme/thingscript/eval"
	"github.com/thingsme/thingscript/lexer"
	"github.com/thingsme/thingscript/object"
	"github.com/thingsme/thingscript/parser"
	"github.com/thingsme/thingscript/stdlib"
)

func TestTime(t *testing.T) {
	timing := time.Now()
	tests := []struct {
		input    string
		expected string
	}{
		{
			input: `
				out := import("fmt")
				time := import("time")
				out.printf("%v", time.Now())
			`,
			expected: fmt.Sprintf("time.Time(%v)", timing),
		},
		{
			input: `
				out := import("fmt")
				time := import("time")
				var t1 = time.Now()
				var t2 = time.Time(t1)  // constructor with time.Time
				out.printf("%v", t2)
			`,
			expected: fmt.Sprintf("time.Time(%v)", timing),
		},
		{
			input: `
				out := import("fmt")
				time := import("time")
				var t_epoch = time.Time(` + strconv.FormatInt(timing.UnixNano(), 10) + `) // constructor with epoch nano
				out.printf("%v", t_epoch)
			`,
			expected: fmt.Sprintf("time.Time(%v)", time.Unix(0, timing.UnixNano())),
		},
		{
			input: `
				out := import("fmt")
				time := import("time")
				var t1 = time.Now()
				var t2 = t1
				out.printf("%v", t2)
			`,
			expected: fmt.Sprintf("time.Time(%v)", timing),
		},
		{
			input: `
				time := import("time")
				var tick time.Time
				tick = time.Now()
				import("fmt").println("time:", tick)
			`,
			expected: fmt.Sprintf("time: time.Time(%s)\n", timing),
		},
	}
	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := parser.New(l)
		program := p.ParseProgram()
		for _, err := range p.Errors() {
			t.Errorf("parse error: %s", err)
		}
		env := object.NewEnvironment()
		out := &bytes.Buffer{}
		env.Stdout = out
		env.TimeProvider = func() time.Time { return timing }
		env.RegisterPackages(stdlib.Packages()...)

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
