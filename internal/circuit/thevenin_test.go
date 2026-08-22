// Package circuit: Thevenin equivalent and maximum-slip tests.
package circuit

import (
	"math"
	"testing"
)

func TestTheveninOpenCircuitVoltage(t *testing.T) {
	vs := complex(220, 0)
	zs := Impedance(1.5, 2.0)
	zm := Impedance(30, 20)
	eq := TheveninOf(vs, zs, zm)

	// Vth = Vph * Zm / (Zs + Zm). The magnetising branch dominates the
	// divider, so Vth is close to the phase voltage.
	vth := vs * zm / (zs + zm)
	if d := Mag(eq.V - vth); d > 1e-9 {
		t.Errorf("Vth = %v, want %v", eq.V, vth)
	}
	// Zth = Zs || Zm.
	zth := Parallel(zs, zm)
	if d := Mag(eq.Z - zth); d > 1e-9 {
		t.Errorf("Zth = %v, want %v", eq.Z, zth)
	}
}

func TestMaxSlipFormula(t *testing.T) {
	eq := TheveninOf(complex(219.392, 0), Impedance(1.5, 2.0), Impedance(30, 20))
	got := MaxSlip(0.8, 1.5, eq)
	rth := Resistance(eq.Z)
	xth := Reactance(eq.Z) + 1.5
	want := 0.8 / math.Hypot(rth, xth)
	if math.Abs(got-want) > 1e-12 {
		t.Errorf("MaxSlip = %v, want %v", got, want)
	}
	if got <= 0 || got >= 1 {
		t.Errorf("MaxSlip = %v, want in (0,1)", got)
	}
}

func TestMaxSlipDoublesWithR2(t *testing.T) {
	eq := TheveninOf(complex(219.392, 0), Impedance(1.5, 2.0), Impedance(30, 20))
	base := MaxSlip(0.8, 1.5, eq)
	doubled := MaxSlip(1.6, 1.5, eq)
	// r2 doubled, the denominator unchanged: sm must double exactly.
	if math.Abs(doubled-base*2) > 1e-12 {
		t.Errorf("MaxSlip(r2*2) = %v, want %v", doubled, base*2)
	}
}

func TestTorqueTheveninMatchesPowerMethod(t *testing.T) {
	// The closed-form Thevenin torque must agree with the torque derived
	// from the air-gap power at the same slip.
	vs := complex(219.392, 0)
	zs := Impedance(1.5, 2.0)
	zm := Impedance(30, 20)
	zr := RotorImpedance(0.8, 1.5, 0.03)
	c := Circuit{Vs: vs, Zs: zs, Zm: zm, Zr: zr}
	sol := c.Solve()
	pow := sol.Powers(vs, zs, 30, 0.8, 0.03)

	ws := 2 * math.Pi * 50
	polePairs := 2.0
	fromPower := polePairs * pow.Airgap / ws

	eq := TheveninOf(vs, zs, zm)
	fromThevenin := TorqueFromThevenin(polePairs, ws, 0.8, 1.5, eq, 0.03)

	if d := math.Abs(fromThevenin - fromPower); d/fromPower > 1e-9 {
		t.Errorf("Thevenin torque %v differs from power-method torque %v", fromThevenin, fromPower)
	}
}

func TestTorqueTendsToZeroAtTinySlip(t *testing.T) {
	eq := TheveninOf(complex(219.392, 0), Impedance(1.5, 2.0), Impedance(30, 20))
	t1 := TorqueFromThevenin(2, 2*math.Pi*50, 0.8, 1.5, eq, 1e-3)
	t2 := TorqueFromThevenin(2, 2*math.Pi*50, 0.8, 1.5, eq, 1e-6)
	if !(t2 < t1) {
		t.Errorf("torque must fall with slip: T(1e-6)=%v, T(1e-3)=%v", t2, t1)
	}
}

func TestBreakdownTorqueMatchesMaxTorque(t *testing.T) {
	eq := TheveninOf(complex(219.392, 0), Impedance(1.5, 2.0), Impedance(30, 20))
	ws := 2 * math.Pi * 50
	viaSlip := MaxTorque(2, ws, 0.8, 1.5, eq)
	closedForm := BreakdownTorque(2, ws, 0.8, 1.5, eq)
	if d := math.Abs(viaSlip - closedForm); d/closedForm > 1e-12 {
		t.Errorf("BreakdownTorque %v differs from MaxTorque at sm %v", closedForm, viaSlip)
	}
}
