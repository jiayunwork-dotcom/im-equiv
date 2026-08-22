// Package config: JSON input model for the induction-machine calculator.
package config

import (
	"fmt"
	"os"
)

// LoadFile reads a JSON input document from disk and parses it. File I/O
// errors and parse errors are both surfaced as ordinary errors so the
// command-line front end can print them on stderr and exit non-zero.
func LoadFile(path string) (Input, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return Input{}, fmt.Errorf("cannot read input file %s: %w", path, err)
	}
	in, err := Parse(data)
	if err != nil {
		return Input{}, fmt.Errorf("input file %s: %w", path, err)
	}
	return in, nil
}
