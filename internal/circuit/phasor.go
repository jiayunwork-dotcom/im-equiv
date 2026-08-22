// Package circuit: phasor arithmetic helpers.
//
// The per-phase analysis of a balanced three-phase system reduces every
// steady-state quantity to a single complex number: a voltage phasor, a
// current phasor or an impedance. These helpers give the magnitude, angle
// and component extraction used throughout the solver, plus the named
// product, quotient and rotation operations that keep phasor identities
// readable.
package circuit

import (
	"math"
)

// All branch voltages and currents in this package are RMS phasors stored
// as complex128 values. These small helpers keep the arithmetic readable and
// give every quantity a named meaning instead of scattering complex()
// literals through the solver.

// Mag returns the RMS magnitude of a phasor.
func Mag(z complex128) float64 {
	return math.Hypot(real(z), imag(z))
}

// MagSq returns the squared magnitude without computing a square root. It is
// used where only ratios matter, such as the torque formula of the Thevenin
// equivalent.
func MagSq(z complex128) float64 {
	return real(z)*real(z) + imag(z)*imag(z)
}

// Phase returns the phasor angle in radians, in the range (-pi, pi].
func Phase(z complex128) float64 {
	return math.Atan2(imag(z), real(z))
}

// PhaseDegrees is Phase converted to degrees for human-readable output.
func PhaseDegrees(z complex128) float64 {
	return Phase(z) * 180 / math.Pi
}

// Conjugate returns the complex conjugate of a phasor.
func Conjugate(z complex128) complex128 {
	return complex(real(z), -imag(z))
}

// FromPolar builds a phasor from a magnitude and an angle in radians.
func FromPolar(mag, angleRad float64) complex128 {
	return complex(mag*math.Cos(angleRad), mag*math.Sin(angleRad))
}

// FromPolarDegrees is FromPolar with the angle given in degrees.
func FromPolarDegrees(mag, angleDeg float64) complex128 {
	return FromPolar(mag, angleDeg*math.Pi/180)
}

// Real returns the real (resistive) component of a phasor.
func Real(z complex128) float64 {
	return real(z)
}

// Imag returns the imaginary (reactive) component of a phasor.
func Imag(z complex128) float64 {
	return imag(z)
}

// Resistance extracts the resistive part of an impedance.
func Resistance(z complex128) float64 {
	return real(z)
}

// Reactance extracts the reactive part of an impedance.
func Reactance(z complex128) float64 {
	return imag(z)
}

// Ratio returns the ratio of two magnitudes, guarding against a zero
// denominator. It is used for relative comparisons such as the voltage-
// scaling rule.
func Ratio(a, b complex128) float64 {
	mb := Mag(b)
	if mb == 0 {
		return math.Inf(1)
	}
	return Mag(a) / mb
}

// Product multiplies two phasors. It is a named form of the built-in
// operator that keeps phasor identities explicit at call sites.
func Product(a, b complex128) complex128 {
	return a * b
}

// Quotient divides two phasors. Like Product it is a named form of the
// built-in division, used when the meaning is "the ratio of these two
// phasors" rather than a generic arithmetic step.
func Quotient(a, b complex128) complex128 {
	return a / b
}

// Rotate rotates a phasor by an angle in radians, multiplying by e^{jθ}.
func Rotate(z complex128, angleRad float64) complex128 {
	return z * complex(math.Cos(angleRad), math.Sin(angleRad))
}

// PhaseShift is Rotate for an angle given in degrees.
func PhaseShift(z complex128, angleDeg float64) complex128 {
	return Rotate(z, angleDeg*math.Pi/180)
}
