// Package config: JSON input model for the induction-machine calculator.
//
// Parsing is deliberately separate from validation: Parse answers "is this a
// well-formed JSON document with the right shape", while Validate answers
// "is this operating point physically meaningful". Keeping the two apart
// means a malformed file and a physically absurd file get distinct, precise
// error messages instead of one vague failure.
package config

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
)

// Parse decodes a JSON document into an Input. The decoder is strict in two
// ways that matter for a calculator:
//
//   - unknown fields are rejected, so a misspelled key such as "voltage"
//     fails loudly instead of silently defaulting to zero;
//   - any data after the closing brace is rejected, so a truncated or
//     concatenated document does not pass as valid.
//
// Missing fields are not rejected here: a missing "p" decodes to zero and is
// caught by Validate, whereas a missing "Rs" simply means a zero stator
// resistance, which is physically meaningful. Keeping the parser and the
// physical validation separate makes each failure mode easy to test.
func Parse(data []byte) (Input, error) {
	dec := json.NewDecoder(bytes.NewReader(data))
	dec.DisallowUnknownFields()

	var in Input
	if err := dec.Decode(&in); err != nil {
		return Input{}, fmt.Errorf("invalid JSON input: %w", err)
	}
	if err := rejectTrailing(dec); err != nil {
		return Input{}, err
	}
	return in, nil
}

// rejectTrailing fails when the document contains anything after the single
// top-level JSON value that was just decoded.
func rejectTrailing(dec *json.Decoder) error {
	var extra any
	err := dec.Decode(&extra)
	if err == io.EOF {
		return nil
	}
	if err != nil {
		return fmt.Errorf("invalid JSON input: %w", err)
	}
	return fmt.Errorf("invalid JSON input: unexpected data after the document")
}
