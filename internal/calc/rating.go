// Package calc: ratings and derived quantities.
package calc

import (
	"math"

	"im-equiv/internal/circuit"
	"im-equiv/internal/config"
)

// PullOutRatio is the ratio of the maximum (breakdown) torque to the torque
// at the given operating point. It is the standard overload capability of
// the machine at that slip; induction motors typically show a pull-out ratio
// between two and three.
func PullOutRatio(maxTorqueNm, operatingTorqueNm float64) float64 {
	if operatingTorqueNm == 0 {
		return math.Inf(1)
	}
	return maxTorqueNm / operatingTorqueNm
}

// BreakdownTorqueNm returns the maximum torque the machine can develop from
// the Thevenin model, independent of any single operating point. It is
// exposed as a convenience over the field of the same name in Results.
func BreakdownTorqueNm(in config.Input) float64 {
	res, err := Compute(in)
	if err != nil {
		return 0
	}
	return res.MaxTorqueNm
}

// SlipAtSpeed returns the slip needed to run at a given rotor speed in rpm.
// It is the RPM-based companion of SlipForSpeed: at the synchronous speed
// the slip is zero, and below it the slip is positive.
func SlipAtSpeed(f float64, poleCount int, rpm float64) float64 {
	sync := SyncSpeedRPM(f, poleCount)
	if sync == 0 {
		return 1
	}
	return 1 - rpm/sync
}

// StartingTorqueNm is the torque developed at standstill, where the slip is
// one. It is evaluated through the general torque curve and is the quantity
// that decides whether the motor can accelerate its load from rest.
func StartingTorqueNm(in config.Input) float64 {
	return TorqueAtSlip(in, 1.0)
}

// StartingCurrentRatio returns how many times the no-load current the
// standstill current is. It reuses the circuit-level characteristic points
// so the ratio is computed against the same phasor model as every other
// quantity.
func StartingCurrentRatio(in config.Input) float64 {
	c := buildCircuit(in)
	return circuit.CurrentRatio(c)
}

// buildCircuit assembles the per-phase circuit for an input at its stated
// slip, mirroring the construction inside compute without revalidating.
func buildCircuit(in config.Input) circuit.Circuit {
	vs := complex(in.PhaseVoltage(), 0)
	c := circuit.Circuit{
		Vs:        vs,
		Zs:        circuit.StatorImpedance(in.Rs, in.Xs),
		Zm:        circuit.MagnetisingImpedance(in.Rm, in.Xm),
		RotorOpen: in.S == 0,
	}
	if !c.RotorOpen {
		c.Zr = circuit.RotorImpedance(in.R2, in.X2, in.S)
	}
	return c
}
