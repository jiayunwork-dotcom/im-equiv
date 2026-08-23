// Package config: validation tests covering the failure modes required by
// the problem statement: non-positive pole count, non-positive frequency and
// negative circuit parameters must all produce an error.
package config

import (
	"math"
	"strings"
	"testing"
)

func TestValidateRejectsNonPositivePoleCount(t *testing.T) {
	cases := []struct {
		name string
		p    int
	}{
		{"zero", 0},
		{"negative", -1},
	}
	for _, tc := range cases {
		in := valid()
		in.P = tc.p
		if err := Validate(in); err == nil {
			t.Errorf("p=%d: expected error, got none", tc.p)
		}
	}
}

func TestValidateRejectsNonPositiveFrequency(t *testing.T) {
	cases := []struct {
		name string
		f    float64
	}{
		{"zero", 0},
		{"negative", -50},
	}
	for _, tc := range cases {
		in := valid()
		in.F = tc.f
		if err := Validate(in); err == nil {
			t.Errorf("f=%g: expected error, got none", tc.f)
		}
	}
}

func TestValidateRejectsNegativeResistance(t *testing.T) {
	fields := []struct {
		name string
		set  func(*Input, float64)
	}{
		{"Rs", func(in *Input, v float64) { in.Rs = v }},
		{"Rm", func(in *Input, v float64) { in.Rm = v }},
		{"r2", func(in *Input, v float64) { in.R2 = v }},
	}
	for _, f := range fields {
		in := valid()
		f.set(&in, -0.5)
		if err := Validate(in); err == nil {
			t.Errorf("%s = -0.5: expected error, got none", f.name)
		}
	}
}

func TestValidateRejectsNegativeReactance(t *testing.T) {
	fields := []struct {
		name string
		set  func(*Input, float64)
	}{
		{"Xs", func(in *Input, v float64) { in.Xs = v }},
		{"Xm", func(in *Input, v float64) { in.Xm = v }},
		{"x2", func(in *Input, v float64) { in.X2 = v }},
	}
	for _, f := range fields {
		in := valid()
		f.set(&in, -0.5)
		if err := Validate(in); err == nil {
			t.Errorf("%s = -0.5: expected error, got none", f.name)
		}
	}
}

func TestValidateRejectsNonPositiveLineVoltage(t *testing.T) {
	in := valid()
	in.V = 0
	if err := Validate(in); err == nil {
		t.Error("V=0: expected error, got none")
	}
}

func TestValidateAcceptsGeneratingSlip(t *testing.T) {
	in := valid()
	in.S = -0.02
	if err := Validate(in); err != nil {
		t.Errorf("negative slip must be legal (generating), got error: %v", err)
	}
}

func TestValidateAcceptsSynchronousSlip(t *testing.T) {
	in := valid()
	in.S = 0
	if err := Validate(in); err != nil {
		t.Errorf("zero slip must be legal (synchronous), got error: %v", err)
	}
}

func TestValidateRejectsNaN(t *testing.T) {
	in := valid()
	in.S = math.NaN()
	if err := Validate(in); err == nil {
		t.Error("s=NaN: expected error, got none")
	}
}

func TestValidateAllAggregatesEveryError(t *testing.T) {
	// Several defects at once must all be reported by ValidateAll, while
	// Validate stops at the first one.
	bad := Input{P: 0, F: 0, V: 0, Rs: -1, R2: -1}
	err := ValidateAll(bad)
	if err == nil {
		t.Fatal("ValidateAll(bad): expected errors, got none")
	}
	msg := err.Error()
	for _, want := range []string{"pole count", "frequency", "line voltage", "Rs", "r2"} {
		if !strings.Contains(msg, want) {
			t.Errorf("ValidateAll error %q missing %q", msg, want)
		}
	}
	first := Validate(bad)
	if first == nil {
		t.Error("Validate(bad): expected first error, got none")
	}
}

// valid returns an input that passes every validation check. p = 4 is the
// pole count of the shipped four-pole example machine.
func valid() Input {
	return Input{P: 4, F: 50, V: 380, Rs: 1.5, Xs: 2.0, Rm: 30.0, Xm: 20.0, R2: 0.8, X2: 1.5, S: 0.03}
}
