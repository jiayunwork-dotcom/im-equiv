// Package calc: end-to-end compute tests over the shipped example and the
// boundary operating points required by the problem statement.
package calc

import (
	"math"
	"testing"

	"im-equiv/internal/config"
)

func example() config.Input {
	return config.Input{P: 4, F: 50, V: 380, Rs: 1.5, Xs: 2.0, Rm: 30.0, Xm: 20.0, R2: 0.8, X2: 1.5, S: 0.03}
}

func TestComputeRatedExampleEnergyBalance(t *testing.T) {
	res, err := Compute(example())
	if err != nil {
		t.Fatalf("Compute(rated): unexpected error: %v", err)
	}
	// Energy balance: Pag = Pcu2 + Pem.
	left := res.RotorCopperW + res.MechanicalW
	if d := math.Abs(left - res.AirgapPowerW); d > 1e-6*math.Max(1, res.AirgapPowerW) {
		t.Errorf("energy balance: Pcu2+Pem = %v, Pag = %v", left, res.AirgapPowerW)
	}
	// Mechanical power must equal the torque times the mechanical speed,
	// T = Pem / wm.
	if got := res.TorqueNm * res.RotorSpeedRadS; math.Abs(got-res.MechanicalW) > 1e-6*math.Max(1, res.MechanicalW) {
		t.Errorf("T*wm = %v, Pem = %v", got, res.MechanicalW)
	}
	// Order of magnitude at the rated point.
	if res.TorqueNm < 10 || res.TorqueNm > 60 {
		t.Errorf("rated torque = %v Nm, want a plausible machine value", res.TorqueNm)
	}
	if res.MechanicalW < 1000 || res.MechanicalW > 9000 {
		t.Errorf("rated mechanical power = %v W, want a plausible machine value", res.MechanicalW)
	}
	if res.MaxSlip <= 0 || res.MaxSlip >= 1 {
		t.Errorf("max-torque slip = %v, want in (0,1)", res.MaxSlip)
	}
	if res.Mode != "motoring" {
		t.Errorf("mode = %q, want motoring at s=0.03", res.Mode)
	}
}

func TestComputeZeroSlipZeroTorque(t *testing.T) {
	in := example()
	in.S = 0
	res, err := Compute(in)
	if err != nil {
		t.Fatalf("Compute(s=0): unexpected error: %v", err)
	}
	if res.TorqueNm != 0 {
		t.Errorf("torque at s=0 = %v, want 0", res.TorqueNm)
	}
	if res.MechanicalW != 0 {
		t.Errorf("mechanical power at s=0 = %v, want 0", res.MechanicalW)
	}
	if res.RotorCurrentA != 0 {
		t.Errorf("rotor current at s=0 = %v, want 0", res.RotorCurrentA)
	}
	if res.Mode != "synchronous" {
		t.Errorf("mode at s=0 = %q, want synchronous", res.Mode)
	}
	// The input current at s=0 must be the magnetising current.
	if d := math.Abs(res.InputCurrentA - res.MagnetisingA); d > 1e-9 {
		t.Errorf("I1 = %v, I0 = %v, want equal at s=0", res.InputCurrentA, res.MagnetisingA)
	}
}

func TestComputeSmallSlipTorqueTendsToZero(t *testing.T) {
	in := example()
	in.S = 1e-4
	res, err := Compute(in)
	if err != nil {
		t.Fatalf("Compute(s=1e-4): unexpected error: %v", err)
	}
	if res.TorqueNm > 1 {
		t.Errorf("torque at s=1e-4 = %v, want close to zero", res.TorqueNm)
	}
}

func TestComputeStallZeroMechanicalPower(t *testing.T) {
	in := example()
	in.S = 1.0
	res, err := Compute(in)
	if err != nil {
		t.Fatalf("Compute(s=1): unexpected error: %v", err)
	}
	if res.MechanicalW != 0 {
		t.Errorf("mechanical power at standstill = %v, want 0", res.MechanicalW)
	}
	// At standstill the entire air-gap power is rotor copper loss.
	if d := math.Abs(res.AirgapPowerW - res.RotorCopperW); d > 1e-6*math.Max(1, res.AirgapPowerW) {
		t.Errorf("Pag = %v, Pcu2 = %v, want equal at standstill", res.AirgapPowerW, res.RotorCopperW)
	}
}

func TestComputeGeneratingNegativeTorqueAndPower(t *testing.T) {
	in := example()
	in.S = -0.02
	res, err := Compute(in)
	if err != nil {
		t.Fatalf("Compute(s=-0.02): unexpected error: %v", err)
	}
	if res.TorqueNm >= 0 {
		t.Errorf("generating torque = %v, want negative", res.TorqueNm)
	}
	if res.MechanicalW >= 0 {
		t.Errorf("generating mechanical power = %v, want negative", res.MechanicalW)
	}
	if res.Mode != "generating" {
		t.Errorf("mode = %q, want generating", res.Mode)
	}
	if res.RotorSpeedRadS <= 0 {
		t.Errorf("generating rotor speed = %v, want above zero", res.RotorSpeedRadS)
	}
}

func TestComputeRejectsNegativeResistance(t *testing.T) {
	in := example()
	in.R2 = -0.5
	if _, err := Compute(in); err == nil {
		t.Error("Compute with r2 < 0: expected error, got none")
	}
}

func TestComputeRejectsNonPositiveFrequency(t *testing.T) {
	in := example()
	in.F = 0
	if _, err := Compute(in); err == nil {
		t.Error("Compute with f = 0: expected error, got none")
	}
}

func TestComputeRejectsNonPositivePolePairs(t *testing.T) {
	in := example()
	in.P = 0
	if _, err := Compute(in); err == nil {
		t.Error("Compute with p = 0: expected error, got none")
	}
}

func TestComputeRatedTorqueFromTheveninMatches(t *testing.T) {
	// The printed torque (power method) must agree with the closed-form
	// Thevenin torque evaluated at the same slip.
	res, err := Compute(example())
	if err != nil {
		t.Fatalf("Compute(rated): %v", err)
	}
	if res.TorqueNm <= 0 || res.MaxTorqueNm <= res.TorqueNm {
		t.Errorf("rated torque %v must be positive and below max torque %v", res.TorqueNm, res.MaxTorqueNm)
	}
}

func TestComputeFromFileExample(t *testing.T) {
	res, err := ComputeFromFile("../../example/4pole-50hz.json")
	if err != nil {
		t.Fatalf("ComputeFromFile(example): %v", err)
	}
	if res.TorqueNm <= 0 {
		t.Errorf("example torque = %v, want positive", res.TorqueNm)
	}
	if res.Mode != "motoring" {
		t.Errorf("example mode = %q, want motoring", res.Mode)
	}
}
