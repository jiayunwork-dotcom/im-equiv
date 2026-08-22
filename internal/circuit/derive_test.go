// Package circuit: parameter referral tests.
package circuit

import (
	"math"
	"testing"
)

func TestReferToStatorScalesWithSquare(t *testing.T) {
	r2, x2 := ReferToStator(2, 1.5, 2.5)
	// k = 2 scales both by 4.
	if math.Abs(r2-6) > 1e-12 {
		t.Errorf("referred r2 = %v, want 6", r2)
	}
	if math.Abs(x2-10) > 1e-12 {
		t.Errorf("referred x2 = %v, want 10", x2)
	}
}

func TestReferredRatioInverse(t *testing.T) {
	referred, _ := ReferToStator(3, 0.8, 1.5)
	if k := ReferredRatio(referred, 0.8); math.Abs(k-3) > 1e-12 {
		t.Errorf("ReferredRatio = %v, want 3", k)
	}
}

func TestReferredRatioZeroRotorResistance(t *testing.T) {
	if k := ReferredRatio(10, 0); k != 0 {
		t.Errorf("ReferredRatio with zero rotor resistance = %v, want 0", k)
	}
}

func TestPerUnitRealPower(t *testing.T) {
	if pu := PerUnitRealPower(4100, 4100); math.Abs(pu-1) > 1e-12 {
		t.Errorf("PerUnitRealPower(4100, 4100) = %v, want 1", pu)
	}
	if pu := PerUnitRealPower(2050, 4100); math.Abs(pu-0.5) > 1e-12 {
		t.Errorf("PerUnitRealPower(2050, 4100) = %v, want 0.5", pu)
	}
}

func TestPerUnitRealPowerZeroBase(t *testing.T) {
	if pu := PerUnitRealPower(100, 0); pu != 0 {
		t.Errorf("PerUnitRealPower(100, 0) = %v, want 0", pu)
	}
}
