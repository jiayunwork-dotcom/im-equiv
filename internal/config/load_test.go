// Package config: file loading tests.
package config

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestLoadExampleFourPole(t *testing.T) {
	in, err := LoadFile("../../example/4pole-50hz.json")
	if err != nil {
		t.Fatalf("LoadFile(example): unexpected error: %v", err)
	}
	want := Input{P: 4, F: 50, V: 380, Rs: 1.5, Xs: 2.0, Rm: 30.0, Xm: 20.0, R2: 0.8, X2: 1.5, S: 0.03}
	if in != want {
		t.Errorf("loaded example: got %+v, want %+v", in, want)
	}
}

func TestLoadExampleGenerating(t *testing.T) {
	in, err := LoadFile("../../example/4pole-50hz-generating.json")
	if err != nil {
		t.Fatalf("LoadFile(generating): unexpected error: %v", err)
	}
	if in.S >= 0 {
		t.Errorf("generating example slip = %v, want negative", in.S)
	}
}

func TestLoadMissingFileErrors(t *testing.T) {
	if _, err := LoadFile("../../example/does-not-exist.json"); err == nil {
		t.Error("LoadFile(missing): expected error, got none")
	}
}

func TestLoadFileWithInvalidContent(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "bad.json")
	if err := os.WriteFile(path, []byte(`{"p": "not-a-number"}`), 0o644); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}
	if _, err := LoadFile(path); err == nil {
		t.Error("LoadFile(bad type): expected error, got none")
	} else if !strings.Contains(err.Error(), "invalid JSON") {
		t.Errorf("LoadFile(bad type): error = %q, want mention of invalid JSON", err)
	}
}
