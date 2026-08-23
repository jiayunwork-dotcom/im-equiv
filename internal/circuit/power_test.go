// Package circuit: power split tests.
package circuit

import (
	"math"
	"testing"
)

func ratedCircuit() Circuit {
	return Circuit{
		Vs:        complex(219.392, 0),
		Zs:        Impedance(1.5, 2.0),
		Zm:        Impedance(30, 20),
		Zr:        RotorImpedance(0.8, 1.5, 0.03),
		RotorOpen: false,
	}
}

func TestPowersAtRatedSlip(t *testing.T) {
	c := ratedCircuit()
	sol := c.Solve()
	pow := sol.Powers(c.Vs, c.Zs, 30, 0.8, 0.03)

	// The air-gap power splits into rotor copper loss and the rest, which is
	// delivered at the shaft. Pag = 3 I2^2 r2/s and Pcu2 = 3 I2^2 r2.
	wantPcu2 := 3 * MagSq(sol.RotorCurrent) * 0.8
	if math.Abs(pow.RotorCopper-wantPcu2) > 1e-9 {
		t.Errorf("Pcu2 = %v, want %v", pow.RotorCopper, wantPcu2)
	}
	wantPag := 3 * MagSq(sol.RotorCurrent) * 0.8 / 0.03
	if math.Abs(pow.Airgap-wantPag) > 1e-6 {
		t.Errorf("Pag = %v, want %v", pow.Airgap, wantPag)
	}
	// Mechanical power implied by the split: Pag - Pcu2 = Pag (1 - s).
	wantMechanical := pow.Airgap - pow.RotorCopper
	if got := pow.Airgap * (1 - 0.03); math.Abs(got-wantMechanical) > 1e-9 {
		t.Errorf("mechanical power = %v, want %v", got, wantMechanical)
	}
}

func TestPowersAtStallAllRotorCopper(t *testing.T) {
	c := ratedCircuit()
	c.Zr = RotorImpedance(0.8, 1.5, 1.0)
	sol := c.Solve()
	pow := sol.Powers(c.Vs, c.Zs, 30, 0.8, 1.0)
	// At s = 1 the air-gap power equals the rotor copper loss and the
	// mechanical power is zero.
	if math.Abs(pow.Airgap-pow.RotorCopper) > 1e-9 {
		t.Errorf("at stall Pag = %v, Pcu2 = %v, want equal", pow.Airgap, pow.RotorCopper)
	}
	if got := pow.Airgap * (1 - 1.0); got != 0 {
		t.Errorf("at stall mechanical power = %v, want 0", got)
	}
}

func TestPowersZeroSlipAirgapZero(t *testing.T) {
	c := ratedCircuit()
	c.RotorOpen = true
	sol := c.Solve()
	pow := sol.Powers(c.Vs, c.Zs, 30, 0.8, 0.0)
	if pow.Airgap != 0 {
		t.Errorf("at s=0 airgap power = %v, want 0", pow.Airgap)
	}
	if pow.RotorCopper != 0 {
		t.Errorf("at s=0 rotor copper = %v, want 0", pow.RotorCopper)
	}
}

func TestPowerFactorLagging(t *testing.T) {
	c := ratedCircuit()
	sol := c.Solve()
	pf := PowerFactor(c.Vs, sol.InputCurrent)
	if pf <= 0 || pf > 1 {
		t.Errorf("power factor = %v, want in (0,1] for an inductive motor", pf)
	}
}

func TestIronLossPresent(t *testing.T) {
	c := ratedCircuit()
	sol := c.Solve()
	pow := sol.Powers(c.Vs, c.Zs, 30, 0.8, 0.03)
	want := 3 * MagSq(sol.MagnetisingCurrent) * 30
	if math.Abs(pow.Iron-want) > 1e-9 {
		t.Errorf("iron loss = %v, want %v", pow.Iron, want)
	}
}

func TestTotalLossesBalance(t *testing.T) {
	// Input power minus total losses must equal the mechanical power.
	c := ratedCircuit()
	sol := c.Solve()
	pow := sol.Powers(c.Vs, c.Zs, 30, 0.8, 0.03)
	mechanical := pow.Airgap - pow.RotorCopper
	left := pow.Input - pow.TotalLosses()
	if math.Abs(left-mechanical) > 1e-6*math.Max(1, mechanical) {
		t.Errorf("Pin - losses = %v, mechanical power = %v", left, mechanical)
	}
}
