// Package circuit: parameter referral and scaling helpers.
package circuit

import "math"

// ReferToStator scales rotor parameters to the stator side by the winding
// ratio k (stator turns over rotor turns). A rotor resistance and reactance
// measured on the rotor side both scale with k^2 when referred:
//
//	r2' = k^2 r2,   x2' = k^2 x2
//
// The T-equivalent circuit of this tool already receives referred values,
// so the helper exists for users who start from the rotor-side nameplate and
// for tests that want to check the scaling law explicitly.
func ReferToStator(windingRatio, rotorR2, rotorX2 float64) (referredR2, referredX2 float64) {
	k2 := windingRatio * windingRatio
	return k2 * rotorR2, k2 * rotorX2
}

// ReferredRatio is the turns ratio k that would map a given rotor-side
// resistance onto a referred value. It is the inverse of ReferToStator.
func ReferredRatio(referredR2, rotorR2 float64) float64 {
	if rotorR2 == 0 {
		return 0
	}
	return math.Sqrt(referredR2 / rotorR2)
}

// PerUnitRealPower normalises a real power against a base power, returning
// the dimensionless per-unit value. It is a convenience for presenting
// losses and output power relative to the rated value.
func PerUnitRealPower(actualW, baseW float64) float64 {
	if baseW == 0 {
		return 0
	}
	return actualW / baseW
}
