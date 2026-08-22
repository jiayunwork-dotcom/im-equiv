// Package config: parser tests.
package config

import (
	"strings"
	"testing"
)

func TestParseValidInput(t *testing.T) {
	doc := `{"p":4,"f":50,"V":380,"Rs":1.5,"Xs":2.0,"Rm":30.0,"Xm":20.0,"r2":0.8,"x2":1.5,"s":0.03}`
	in, err := Parse([]byte(doc))
	if err != nil {
		t.Fatalf("Parse: unexpected error: %v", err)
	}
	want := Input{P: 4, F: 50, V: 380, Rs: 1.5, Xs: 2.0, Rm: 30.0, Xm: 20.0, R2: 0.8, X2: 1.5, S: 0.03}
	if in != want {
		t.Errorf("Parse: got %+v, want %+v", in, want)
	}
}

func TestParseRejectsUnknownField(t *testing.T) {
	doc := `{"p":4,"f":50,"V":380,"Rs":1.5,"Xs":2.0,"Rm":30.0,"Xm":20.0,"r2":0.8,"x2":1.5,"s":0.03,"voltage":380}`
	if _, err := Parse([]byte(doc)); err == nil {
		t.Error("Parse: expected error for unknown field voltage, got none")
	} else if !strings.Contains(err.Error(), "unknown field") {
		t.Errorf("Parse: error = %q, want mention of the unknown field", err)
	}
}

func TestParseRejectsTrailingGarbage(t *testing.T) {
	doc := `{"p":2,"f":50} extra`
	if _, err := Parse([]byte(doc)); err == nil {
		t.Error("Parse: expected error for trailing data, got none")
	}
}

func TestParseRejectsMalformedJSON(t *testing.T) {
	doc := `{"p":2,`
	if _, err := Parse([]byte(doc)); err == nil {
		t.Error("Parse: expected error for truncated JSON, got none")
	}
}

func TestPhaseVoltageWye(t *testing.T) {
	in := Input{V: 380}
	got := in.PhaseVoltage()
	want := 380.0 / 1.7320508075688772
	if diff := abs(got - want); diff > 1e-9 {
		t.Errorf("PhaseVoltage: got %v, want %v", got, want)
	}
}

func TestInputStringRoundTrip(t *testing.T) {
	in := Input{P: 4, F: 50, V: 380, Rs: 1.5, Xs: 2.0, Rm: 30.0, Xm: 20.0, R2: 0.8, X2: 1.5, S: 0.03}
	desc := in.String()
	for _, token := range []string{"p:4", "f:50", "r2:0.8", "x2:1.5", "s:0.03"} {
		if !strings.Contains(desc, token) {
			t.Errorf("Input.String() = %q, missing %q", desc, token)
		}
	}
}

func abs(v float64) float64 {
	if v < 0 {
		return -v
	}
	return v
}
