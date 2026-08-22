// Package calc: speed conventions.
package calc

import (
	"math"
)

// SyncSpeedRadS returns the electrical synchronous speed in radians per
// second for a supply frequency f in hertz:
//
//	ws = 2 pi f
//
// This is the speed of the rotating field in electrical terms, independent
// of the pole count.
func SyncSpeedRadS(f float64) float64 {
	return 2 * math.Pi * f
}

// RotorSpeedRadS returns the mechanical angular velocity of the rotor in
// radians per second. With p poles and slip s,
//
//	wm = ws (1 - s) / (p / 2)
//
// The division by p/2 pins down the pole convention: a four-pole machine
// (p = 4) runs at exactly half the electrical speed, so at zero slip a 50 Hz
// four-pole machine turns at 1500 rpm. Dropping the (1 - s) factor is the
// classic slip error that overestimates the rated-point power and the
// efficiency together.
func RotorSpeedRadS(f float64, poleCount int, slip float64) float64 {
	ws := SyncSpeedRadS(f)
	return ws * (1 - slip) / (float64(poleCount) / 2)
}

// SyncSpeedRPM converts the synchronous speed to revolutions per minute. The
// familiar textbook form is ns = 120 f / p with p the number of poles: a
// four-pole 50 Hz machine runs synchronously at 1500 rpm.
func SyncSpeedRPM(f float64, poleCount int) float64 {
	return 120 * f / float64(poleCount)
}

// RotorSpeedRPM returns the mechanical rotor speed in revolutions per
// minute, that is (1 - s) times the synchronous speed.
func RotorSpeedRPM(f float64, poleCount int, slip float64) float64 {
	return SyncSpeedRPM(f, poleCount) * (1 - slip)
}

// SlipForSpeed returns the slip that corresponds to a given mechanical rotor
// speed in rad/s, inverting the RotorSpeedRadS relation. It is the inverse
// mapping used to answer "what slip do I need to run at this speed".
func SlipForSpeed(f float64, poleCount int, rotorSpeedRadS float64) float64 {
	ws := SyncSpeedRadS(f) / (float64(poleCount) / 2)
	return 1 - rotorSpeedRadS/ws
}

// PerUnitRotorSpeed is the rotor speed expressed as a fraction of the
// synchronous speed, that is 1 - s. It is useful for plotting the
// torque-speed characteristic against a familiar horizontal axis.
func PerUnitRotorSpeed(slip float64) float64 {
	return 1 - slip
}

// ElectricalFrequencyAtRotor returns the slip frequency of the rotor
// currents, f * s. A running machine sees its induced rotor quantities at
// this frequency rather than at the supply frequency; the "j omega and s in
// the rotor branch" rule is exactly this.
func ElectricalFrequencyAtRotor(f, slip float64) float64 {
	return f * slip
}
