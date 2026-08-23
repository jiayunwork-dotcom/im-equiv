// Package calc: cross-check rule tests.
//
// These exercise the rules that any consistent T-equivalent-circuit model
// must satisfy: zero slip gives zero torque, and torque and mechanical power
// scale with the square of the applied voltage.
package calc

import (
	"math"
	"testing"
)

func TestDoubleVoltageQuadruplesTorque(t *testing.T) {
	base := example()
	res1, err := Compute(base)
	if err != nil {
		t.Fatalf("Compute(base): %v", err)
	}
	doubled := base
	doubled.V = 2 * base.V
	res2, err := Compute(doubled)
	if err != nil {
		t.Fatalf("Compute(V*2): %v", err)
	}
	if got := res2.TorqueNm / res1.TorqueNm; math.Abs(got-4) > 1e-9 {
		t.Errorf("torque ratio with doubled voltage = %v, want 4", got)
	}
}

func TestDoubleVoltageQuadruplesMechanicalPower(t *testing.T) {
	base := example()
	res1, err := Compute(base)
	if err != nil {
		t.Fatalf("Compute(base): %v", err)
	}
	doubled := base
	doubled.V = 2 * base.V
	res2, err := Compute(doubled)
	if err != nil {
		t.Fatalf("Compute(V*2): %v", err)
	}
	if got := res2.MechanicalW / res1.MechanicalW; math.Abs(got-4) > 1e-9 {
		t.Errorf("mechanical power ratio with doubled voltage = %v, want 4", got)
	}
}

func TestDoubleVoltageScalingCheck(t *testing.T) {
	base := example()
	res1, err := Compute(base)
	if err != nil {
		t.Fatalf("Compute(base): %v", err)
	}
	doubled := base
	doubled.V = 2 * base.V
	res2, err := Compute(doubled)
	if err != nil {
		t.Fatalf("Compute(V*2): %v", err)
	}
	if err := CheckVoltageScaling(res1.TorqueNm, res2.TorqueNm, 2.0); err != nil {
		t.Errorf("voltage scaling of torque: %v", err)
	}
	if err := CheckVoltageScaling(res1.MechanicalW, res2.MechanicalW, 2.0); err != nil {
		t.Errorf("voltage scaling of mechanical power: %v", err)
	}
}

func TestZeroSlipZeroTorqueAtAnyVoltage(t *testing.T) {
	for _, v := range []float64{100, 380, 760} {
		in := example()
		in.V = v
		in.S = 0
		res, err := Compute(in)
		if err != nil {
			t.Fatalf("Compute(V=%v, s=0): %v", v, err)
		}
		if res.TorqueNm != 0 {
			t.Errorf("torque at s=0 with V=%v = %v, want 0", v, res.TorqueNm)
		}
	}
}

func TestQuadScalingAirgapPower(t *testing.T) {
	base := example()
	res1, err := Compute(base)
	if err != nil {
		t.Fatalf("Compute(base): %v", err)
	}
	doubled := base
	doubled.V = 2 * base.V
	res2, err := Compute(doubled)
	if err != nil {
		t.Fatalf("Compute(V*2): %v", err)
	}
	if got := res2.AirgapPowerW / res1.AirgapPowerW; math.Abs(got-4) > 1e-9 {
		t.Errorf("air-gap power ratio with doubled voltage = %v, want 4", got)
	}
}
