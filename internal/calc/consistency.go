// Package calc: physical consistency checks.
//
// A T-equivalent-circuit result is only worth printing if it obeys the few
// invariants the model cannot violate. The checks here pin those invariants
// down and are called from Compute as well as exercised directly by the
// tests:
//
//	s = 0   => T = 0                    (synchronous point)
//	Pag     = Pcu2 + Pem                (air-gap energy balance)
//	scaled V => T and Pem scale with V^2 (linear magnetisation)
//
// Each check is deliberately tiny: the value of a guard is that it is cheap
// to run on every computation and impossible to forget.
package calc

import (
	"errors"
	"fmt"
	"math"
)

// ErrZeroSlipTorque is returned when a zero-slip solution carries a non-zero
// electromagnetic torque.
var ErrZeroSlipTorque = errors.New("zero slip must give zero electromagnetic torque")

// ErrEnergyBalance is returned when the air-gap power does not equal the sum
// of the rotor copper loss and the mechanical power.
var ErrEnergyBalance = errors.New("air-gap power does not balance rotor copper loss and mechanical power")

// CheckZeroSlipTorque enforces the synchronous-speed constraint: at a slip
// of exactly zero the electromagnetic torque must vanish. A non-zero torque
// there would mean the rotor branch contributed a slip-dependent term where
// none can exist — the symptom of a slip applied in the wrong place.
func CheckZeroSlipTorque(slip, torqueNm float64) error {
	if slip == 0 && math.Abs(torqueNm) > 1e-9 {
		return fmt.Errorf("%w: got %.6g Nm at s = 0", ErrZeroSlipTorque, torqueNm)
	}
	return nil
}

// CheckEnergyBalance enforces the air-gap power split: the rotor copper loss
// plus the mechanical power must equal the air-gap power. Any drift beyond a
// small relative tolerance indicates a numerical or modelling error and is
// reported instead of being printed as a subtly wrong result.
func CheckEnergyBalance(airgapW, rotorCopperW, mechanicalW float64) error {
	left := rotorCopperW + mechanicalW
	scale := math.Max(1.0, math.Abs(airgapW))
	if math.Abs(left-airgapW) > 1e-9*scale {
		return fmt.Errorf(
			"%w: Pag = %.6g W but Pcu2 + Pem = %.6g W",
			ErrEnergyBalance, airgapW, left,
		)
	}
	return nil
}

// CheckVoltageScaling verifies that a change in the applied line voltage
// scales torque and mechanical power with the square of the voltage ratio,
// as required by the linear-magnetisation model. The check is used by the
// tests to exercise the cross rule "double the voltage, quadruple the
// torque" and by any caller who wants the rule asserted explicitly.
func CheckVoltageScaling(before, after, voltageRatio float64) error {
	expected := before * voltageRatio * voltageRatio
	if math.Abs(after-expected) > 1e-6*math.Max(1.0, math.Abs(expected)) {
		return fmt.Errorf(
			"voltage scaling broken: expected %.6g, got %.6g",
			expected, after,
		)
	}
	return nil
}
