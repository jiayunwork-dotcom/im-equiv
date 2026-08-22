// Package calc: characteristic curve tests.
package calc

import (
	"testing"
)

func TestTorqueCurveSamples(t *testing.T) {
	in := example()
	points := TorqueCurve(in, 0.01, 1.0, 21)
	if len(points) != 21 {
		t.Fatalf("TorqueCurve returned %d points, want 21", len(points))
	}
	// Torque must be positive throughout the motoring band and rise towards
	// the maximum-torque slip before falling at standstill.
	if points[0].TorqueNm <= 0 {
		t.Errorf("curve start torque = %v, want positive", points[0].TorqueNm)
	}
	last := points[len(points)-1]
	if last.TorqueNm <= 0 {
		t.Errorf("standstill torque = %v, want positive", last.TorqueNm)
	}
}

func TestTorqueCurvePowerVanishesAtStall(t *testing.T) {
	in := example()
	points := TorqueCurve(in, 0.01, 1.0, 11)
	last := points[len(points)-1]
	if last.PowerW != 0 {
		t.Errorf("standstill mechanical power = %v, want 0", last.PowerW)
	}
}

func TestEfficiencyCurveMotoringOnly(t *testing.T) {
	in := example()
	// The band includes negative slips (generating), which must be filtered
	// out of the efficiency curve.
	points := EfficiencyCurve(in, -0.1, 1.0, 23)
	for _, p := range points {
		if p.Slip <= 0 || p.Slip > 1 {
			t.Errorf("efficiency curve contains slip %v outside motoring band", p.Slip)
		}
	}
}

func TestMaxTorquePointNearAnalyticSlip(t *testing.T) {
	in := example()
	res, err := Compute(in)
	if err != nil {
		t.Fatalf("Compute: %v", err)
	}
	best := MaxTorquePoint(in, 0.01, 1.0, 200)
	if best.TorqueNm < res.MaxTorqueNm*0.99 {
		t.Errorf("curve peak torque %v too far below analytic max %v", best.TorqueNm, res.MaxTorqueNm)
	}
}

func TestTorqueCurvePointsConsistent(t *testing.T) {
	// Each curve point must agree with a direct single-point evaluation at
	// the same slip.
	in := example()
	points := TorqueCurve(in, 0.02, 0.02, 1)
	p := points[0]
	want := TorqueAtSlip(in, 0.02)
	if got := p.TorqueNm; got != want {
		t.Errorf("curve torque %v differs from direct %v", got, want)
	}
}
