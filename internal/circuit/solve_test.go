// Package circuit: solver tests.
package circuit

import (
	"math"
	"testing"
)

func TestSolveMotoring(t *testing.T) {
	// Rated point of the shipped four-pole 50 Hz example.
	c := Circuit{
		Vs:        complex(219.392, 0),
		Zs:        Impedance(1.5, 2.0),
		Zm:        Impedance(30, 20),
		Zr:        RotorImpedance(0.8, 1.5, 0.03),
		RotorOpen: false,
	}
	sol := c.Solve()

	// Kirchhoff: the input current must split into rotor + magnetising.
	split := sol.MagnetisingCurrent + sol.RotorCurrent
	if d := Mag(split - sol.InputCurrent); d > 1e-9 {
		t.Errorf("node current mismatch: I1-(I0+I2) magnitude = %v, want 0", d)
	}

	if got := Mag(sol.InputCurrent); math.Abs(got-12.22) > 0.02 {
		t.Errorf("I1 = %v, want ~12.22 A", got)
	}
	if got := Mag(sol.RotorCurrent); math.Abs(got-7.27) > 0.02 {
		t.Errorf("I2 = %v, want ~7.27 A", got)
	}
}

func TestSolveZeroSlipOpensRotor(t *testing.T) {
	c := Circuit{
		Vs:        complex(219.392, 0),
		Zs:        Impedance(1.5, 2.0),
		Zm:        Impedance(30, 20),
		RotorOpen: true,
	}
	sol := c.Solve()
	if sol.RotorCurrent != 0 {
		t.Errorf("RotorCurrent = %v, want 0 at synchronous speed", sol.RotorCurrent)
	}
	// With the rotor open the only current path is the magnetising branch,
	// so the input current must equal the magnetising current.
	if d := Mag(sol.InputCurrent - sol.MagnetisingCurrent); d > 1e-12 {
		t.Errorf("I1 differs from I0 at s=0: %v, want 0", d)
	}
}

func TestSolveSmallSlipApproachesMagnetisingCurrent(t *testing.T) {
	c := Circuit{
		Vs:        complex(219.392, 0),
		Zs:        Impedance(1.5, 2.0),
		Zm:        Impedance(30, 20),
		Zr:        RotorImpedance(0.8, 1.5, 1e-6),
		RotorOpen: false,
	}
	sol := c.Solve()
	i1 := Mag(sol.InputCurrent)
	open := Circuit{
		Vs:        c.Vs,
		Zs:        c.Zs,
		Zm:        c.Zm,
		RotorOpen: true,
	}.Solve()
	i1open := Mag(open.InputCurrent)
	if math.Abs(i1-i1open) > 1e-3 {
		t.Errorf("small-slip I1 = %v, want ~ open-rotor I1 = %v", i1, i1open)
	}
}

func TestSolveNodeVoltageBelowPhaseVoltage(t *testing.T) {
	// Stator impedance must drop some voltage, so the node voltage is below
	// the phase voltage in magnitude.
	c := Circuit{
		Vs:        complex(219.392, 0),
		Zs:        Impedance(1.5, 2.0),
		Zm:        Impedance(30, 20),
		Zr:        RotorImpedance(0.8, 1.5, 0.03),
		RotorOpen: false,
	}
	sol := c.Solve()
	if vn := Mag(sol.NodeVoltage); vn >= Mag(c.Vs) {
		t.Errorf("node voltage %v must be below phase voltage %v", vn, Mag(c.Vs))
	}
}

func TestInputPowerPositiveMotoring(t *testing.T) {
	c := Circuit{
		Vs:        complex(219.392, 0),
		Zs:        Impedance(1.5, 2.0),
		Zm:        Impedance(30, 20),
		Zr:        RotorImpedance(0.8, 1.5, 0.03),
		RotorOpen: false,
	}
	sol := c.Solve()
	if p := sol.InputPower(c.Vs); p <= 0 {
		t.Errorf("motoring input power = %v, want positive", p)
	}
}

func TestInputReactivePowerPositive(t *testing.T) {
	// An induction motor draws inductive reactive power for magnetisation.
	c := ratedCircuit()
	sol := c.Solve()
	if q := sol.InputReactivePower(c.Vs); q <= 0 {
		t.Errorf("motoring reactive power = %v, want positive (inductive)", q)
	}
}

func TestPowerFactorAtMatchesHelper(t *testing.T) {
	c := ratedCircuit()
	sol := c.Solve()
	if got, want := sol.PowerFactorAt(c.Vs), PowerFactor(c.Vs, sol.InputCurrent); got != want {
		t.Errorf("PowerFactorAt = %v, PowerFactor = %v", got, want)
	}
}
