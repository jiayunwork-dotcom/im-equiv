// Package calc: output rendering tests.
package calc

import (
	"strings"
	"testing"
)

func TestFormatResultsContainsQuantities(t *testing.T) {
	in := example()
	res, err := Compute(in)
	if err != nil {
		t.Fatalf("Compute: %v", err)
	}
	out := FormatResults(res)

	for _, label := range []string{
		"input current I1",
		"rotor current I2",
		"torque T",
		"mechanical Pem",
		"max-torque slip sm",
		"airgap power Pag",
		"rotor copper Pcu2",
	} {
		if !strings.Contains(out, label) {
			t.Errorf("formatted output missing label %q", label)
		}
	}
}

func TestFormatResultsHasPlausibleNumbers(t *testing.T) {
	in := example()
	res, err := Compute(in)
	if err != nil {
		t.Fatalf("Compute: %v", err)
	}
	out := FormatResults(res)
	if !strings.Contains(out, "motoring") {
		t.Errorf("formatted output missing operating mode, got:\n%s", out)
	}
	for _, bad := range []string{"NaN", "+Inf", "-Inf"} {
		if strings.Contains(out, bad) {
			t.Errorf("formatted output contains %q:\n%s", bad, out)
		}
	}
}

func TestFormatCSVColumns(t *testing.T) {
	in := example()
	res, err := Compute(in)
	if err != nil {
		t.Fatalf("Compute: %v", err)
	}
	out := FormatCSV(res)
	if !strings.HasPrefix(out, "label,value,unit\n") {
		t.Errorf("CSV header = %q, want label,value,unit", out)
	}
	for _, label := range []string{"I1,", "T,", "Pem,", "sm,", "mode,"} {
		if !strings.Contains(out, label) {
			t.Errorf("CSV output missing row for %q", label)
		}
	}
}

func TestFormatSummaryOneLine(t *testing.T) {
	in := example()
	res, err := Compute(in)
	if err != nil {
		t.Fatalf("Compute: %v", err)
	}
	out := FormatSummary(res)
	if strings.Count(out, "\n") != 0 {
		t.Errorf("summary must be a single line, got %q", out)
	}
	if !strings.Contains(out, "motoring") {
		t.Errorf("summary missing mode, got %q", out)
	}
}

func TestFormatJSONValid(t *testing.T) {
	in := example()
	res, err := Compute(in)
	if err != nil {
		t.Fatalf("Compute: %v", err)
	}
	out, err := FormatJSON(res)
	if err != nil {
		t.Fatalf("FormatJSON: %v", err)
	}
	for _, key := range []string{`"t_nm"`, `"pem_w"`, `"sm"`, `"mode"`} {
		if !strings.Contains(out, key) {
			t.Errorf("JSON output missing key %q", key)
		}
	}
}
