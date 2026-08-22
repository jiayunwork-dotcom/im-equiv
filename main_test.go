// Package main: command-line entry tests.
//
// The tests exercise runTorque directly because main itself calls os.Exit,
// which cannot run inside a test. Every failure path the command exposes —
// wrong argument count, unreadable file, invalid input — must surface as an
// error here; the caller in main turns that error into a stderr message and
// a non-zero exit code.
package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRunTorqueExample(t *testing.T) {
	var buf bytes.Buffer
	if err := runTorque([]string{"example/4pole-50hz.json"}, &buf); err != nil {
		t.Fatalf("runTorque(example): unexpected error: %v", err)
	}
	for _, want := range []string{"torque T", "mechanical Pem", "max-torque slip sm", "motoring"} {
		if !strings.Contains(buf.String(), want) {
			t.Errorf("output missing %q, got:\n%s", want, buf.String())
		}
	}
}

func TestRunTorqueWrongArgumentCount(t *testing.T) {
	var buf bytes.Buffer
	err := runTorque(nil, &buf)
	if err == nil {
		t.Fatal("runTorque with no arguments: expected error, got none")
	}
	if !strings.Contains(err.Error(), "exactly one") {
		t.Errorf("error = %q, want mention of the argument requirement", err)
	}
}

func TestRunTorqueMissingFile(t *testing.T) {
	var buf bytes.Buffer
	err := runTorque([]string{"example/not-there.json"}, &buf)
	if err == nil {
		t.Fatal("runTorque(missing file): expected error, got none")
	}
}

func TestRunTorqueInvalidInput(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "bad.json")
	// A negative rotor resistance must be rejected by validation.
	data := `{"p":4,"f":50,"V":380,"Rs":1.5,"Xs":2.0,"Rm":30.0,"Xm":20.0,"r2":-0.8,"x2":1.5,"s":0.03}`
	if err := os.WriteFile(path, []byte(data), 0o644); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}
	var buf bytes.Buffer
	err := runTorque([]string{path}, &buf)
	if err == nil {
		t.Fatal("runTorque(negative r2): expected error, got none")
	}
	if !strings.Contains(err.Error(), "r2") {
		t.Errorf("error = %q, want mention of the offending parameter", err)
	}
}
