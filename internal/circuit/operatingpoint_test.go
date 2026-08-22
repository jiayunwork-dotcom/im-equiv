// Package circuit: characteristic operating point tests.
package circuit

import (
	"testing"
)

func TestNoLoadCurrent(t *testing.T) {
	c := ratedCircuit()
	// At no load the input current equals the open-rotor solution.
	sol := Circuit{Vs: c.Vs, Zs: c.Zs, Zm: c.Zm, RotorOpen: true}.Solve()
	want := Mag(sol.InputCurrent)
	if got := Mag(NoLoadCurrent(c)); abs(got-want) > 1e-9 {
		t.Errorf("NoLoadCurrent = %v, want %v", got, want)
	}
}

func TestLockedRotorCurrent(t *testing.T) {
	c := ratedCircuit()
	c.Zr = RotorImpedance(0.8, 1.5, 1.0)
	sol := c.Solve()
	want := Mag(sol.InputCurrent)
	if got := Mag(LockedRotorCurrent(c)); abs(got-want) > 1e-9 {
		t.Errorf("LockedRotorCurrent = %v, want %v", got, want)
	}
}

func TestStartingCurrentExceedsRated(t *testing.T) {
	// Starting current of an induction motor is several times the no-load
	// current; here the locked-rotor current must be well above the
	// magnetising current.
	c := ratedCircuit()
	c.Zr = RotorImpedance(0.8, 1.5, 1.0)
	nl := Mag(NoLoadCurrent(c))
	lr := Mag(LockedRotorCurrent(c))
	if !(lr > nl) {
		t.Errorf("locked-rotor current %v must exceed no-load current %v", lr, nl)
	}
}

func TestNoLoadPowerPositive(t *testing.T) {
	c := ratedCircuit()
	if p := NoLoadPower(c); p <= 0 {
		t.Errorf("no-load power = %v, want positive", p)
	}
}

func TestCharacteristicsBundled(t *testing.T) {
	c := ratedCircuit()
	c.Zr = RotorImpedance(0.8, 1.5, 1.0)
	op := Characteristics(c)
	if op.NoLoadCurrent <= 0 || op.LockedRotor <= 0 {
		t.Errorf("characteristic currents must be positive: %+v", op)
	}
	if op.NoLoadPower <= 0 || op.LockedRotorPower <= 0 {
		t.Errorf("characteristic powers must be positive: %+v", op)
	}
}

func TestCurrentRatioAboveOne(t *testing.T) {
	c := ratedCircuit()
	c.Zr = RotorImpedance(0.8, 1.5, 1.0)
	if r := CurrentRatio(c); r <= 1 {
		t.Errorf("starting/no-load current ratio = %v, want above 1", r)
	}
}

func abs(v float64) float64 {
	if v < 0 {
		return -v
	}
	return v
}
