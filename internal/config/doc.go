// Package config defines the JSON input of the induction-machine calculator
// and validates it before any arithmetic runs.
//
// # Input document
//
// The command-line tool reads a single JSON object with the fields
// documented on Input. The example below is the shipped four-pole, 50 Hz
// machine (example/4pole-50hz.json):
//
//	{
//	  "p":  4,
//	  "f":  50,
//	  "V":  380,
//	  "Rs": 1.5,
//	  "Xs": 2.0,
//	  "Rm": 30.0,
//	  "Xm": 20.0,
//	  "r2": 0.8,
//	  "x2": 1.5,
//	  "s":  0.03
//	}
//
// All quantities are per-phase and RMS. The stator is assumed to be wye
// connected, so the per-phase voltage is V / sqrt(3). The pole count p is
// the number of poles: a four-pole machine has p = 4 and two pole pairs.
//
// # Failure modes
//
// Loading and parsing are strict: a malformed document, a trailing JSON
// value, or a misspelled field name all produce an error instead of being
// silently ignored. Validation then rejects inputs that make the physical
// problem ill posed:
//
//   - pole count p <= 0;
//   - supply frequency f <= 0;
//   - line voltage V <= 0;
//   - any negative resistance or reactance;
//   - non-finite numbers (NaN, +/-Inf) in any field.
//
// Validate reports the first problem, while ValidateAll aggregates every
// problem into a single error so the user can fix the whole file at once.
//
// A slip of zero is legal (synchronous operation), as is a negative slip
// (generating). The consistency of the computed results — including the rule
// that zero slip must give zero torque — is enforced by package calc.
package config
