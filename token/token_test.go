package token

import "testing"

func TestTokens(t *testing.T) {
	tests := []struct {
		input    string
		expected TokenType
	}{
		{"int", TYPEDECL},
		{"xyz", IDENT},
	}

	for _, tt := range tests {
		ret := LookupIdent(tt.input)
		if ret != tt.expected {
			t.Errorf("ERR expected %q, got=%q", tt.expected, ret)
		}
	}
}

func TestPosition(t *testing.T) {
	p := Position{1, 2}
	expected := "Ln 1, Col 2"
	if p.String() != expected {
		t.Errorf("ERR expected %q, got=%q", expected, p.String())
	}
}
