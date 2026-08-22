// Package config: JSON input model for the induction-machine calculator.
package config

import (
	"errors"
	"fmt"
	"math"
)

// Validate checks every physical constraint that makes the input usable for
// the T-equivalent circuit. A non-positive pole count, a non-positive
// frequency, a non-positive line voltage or a negative circuit parameter all
// make the problem ill posed and must stop the computation before any
// arithmetic runs.
//
// The checks mirror the failure modes listed in the problem statement:
//
//	p <= 0     -> error (a machine needs at least two poles)
//	f  <= 0    -> error (a static field induces nothing)
//	V  <= 0    -> error (no excitation, no torque)
//	R < 0      -> error (negative resistance is unphysical)
//	X < 0      -> error (the T-equivalent circuit is purely inductive)
//
// A negative slip is deliberately accepted: s < 0 is the generating mode and
// the calculation is required to handle it. Validate reports the first
// problem found; ValidateAll reports every problem at once.
func Validate(in Input) error {
	errs := validateAll(in)
	if len(errs) == 0 {
		return nil
	}
	return errs[0]
}

// ValidateAll reports every validation failure of the input in a single
// error, joined with errors.Join. It is the right tool for a user who wants
// to fix the whole input file in one pass instead of discovering each bad
// field one run at a time.
func ValidateAll(in Input) error {
	return errors.Join(validateAll(in)...)
}

// validateAll collects the independent checks. Each check stands alone so a
// caller can see every defective field, and the shared error text keeps the
// single-field Validate consistent with the aggregate ValidateAll.
func validateAll(in Input) []error {
	var errs []error

	add := func(bad bool, msg string, args ...any) {
		if bad {
			errs = append(errs, fmt.Errorf(msg, args...))
		}
	}

	add(in.P <= 0, "pole count p must be positive, got %d", in.P)
	add(in.F <= 0, "frequency f must be positive, got %g Hz", in.F)
	add(in.V <= 0, "line voltage V must be positive, got %g V", in.V)
	add(in.Rs < 0, "stator resistance Rs must not be negative, got %g", in.Rs)
	add(in.Xs < 0, "stator leakage reactance Xs must not be negative, got %g", in.Xs)
	add(in.Rm < 0, "magnetising resistance Rm must not be negative, got %g", in.Rm)
	add(in.Xm < 0, "magnetising reactance Xm must not be negative, got %g", in.Xm)
	add(in.R2 < 0, "rotor resistance r2 must not be negative, got %g", in.R2)
	add(in.X2 < 0, "rotor leakage reactance x2 must not be negative, got %g", in.X2)

	values := []struct {
		name string
		v    float64
	}{
		{"f", in.F}, {"V", in.V},
		{"Rs", in.Rs}, {"Xs", in.Xs},
		{"Rm", in.Rm}, {"Xm", in.Xm},
		{"r2", in.R2}, {"x2", in.X2},
		{"s", in.S},
	}
	for _, item := range values {
		add(math.IsNaN(item.v) || math.IsInf(item.v, 0),
			"%s must be a finite number, got %g", item.name, item.v)
	}

	return errs
}
