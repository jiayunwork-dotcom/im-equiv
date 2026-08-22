// Package calc: power-share and frequency tests.
package calc

import (
	"math"
	"testing"
)

func TestRotorCopperShareEqualsSlip(t *testing.T) {
	for _, s := range []float64{0.03, 0.5, 1.0} {
		if got := RotorCopperShare(s); got != s {
			t.Errorf("RotorCopperShare(%v) = %v, want %v", s, got, s)
		}
	}
}

func TestMechanicalShareComplements(t *testing.T) {
	for _, s := range []float64{0.03, 0.5, 1.0} {
		if got := MechanicalShare(s); math.Abs(got-(1-s)) > 1e-15 {
			t.Errorf("MechanicalShare(%v) = %v, want %v", s, got, 1-s)
		}
	}
}

func TestSharesSumToUnity(t *testing.T) {
	if got := RotorCopperShare(0.04) + MechanicalShare(0.04); math.Abs(got-1) > 1e-15 {
		t.Errorf("shares sum = %v, want 1", got)
	}
}

func TestElectricalFrequencyAtRotor(t *testing.T) {
	// At 3% slip on a 50 Hz supply the rotor currents run at 1.5 Hz.
	if got := ElectricalFrequencyAtRotor(50, 0.03); math.Abs(got-1.5) > 1e-12 {
		t.Errorf("ElectricalFrequencyAtRotor(50, 0.03) = %v, want 1.5", got)
	}
}

func TestRotorI2RLossMatchesCopperLoss(t *testing.T) {
	// RotorI2RLoss must agree with the copper-loss formula used by the
	// circuit package.
	in := example()
	res, err := Compute(in)
	if err != nil {
		t.Fatalf("Compute: %v", err)
	}
	if got := RotorI2RLoss(res.RotorCurrentA, 0.8); math.Abs(got-res.RotorCopperW) > 1e-9 {
		t.Errorf("RotorI2RLoss = %v, want %v", got, res.RotorCopperW)
	}
}
