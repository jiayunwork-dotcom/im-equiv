// Package calc: maximum-slip location tests.
package calc

import (
	"math"
	"testing"
)

func TestNumericMaxSlipMatchesFormula(t *testing.T) {
	in := example()
	res, err := Compute(in)
	if err != nil {
		t.Fatalf("Compute: %v", err)
	}
	numeric := FindMaxTorqueSlip(in, 1e-4, 1.0)
	if d := math.Abs(numeric - res.MaxSlip); d > 1e-6 {
		t.Errorf("numeric sm = %v, analytic sm = %v", numeric, res.MaxSlip)
	}
}

func TestMaxSlipIsTorquePeak(t *testing.T) {
	in := example()
	res, err := Compute(in)
	if err != nil {
		t.Fatalf("Compute: %v", err)
	}
	// Evaluate the torque a small distance either side of the analytic
	// maximum; both must be lower.
	sm := res.MaxSlip
	step := 1e-3
	at := TorqueAtSlip(in, sm)
	left := TorqueAtSlip(in, math.Max(sm-step, 1e-4))
	right := TorqueAtSlip(in, sm+step)
	if at < left || at < right {
		t.Errorf("torque at sm (%v) not above neighbours (%v, %v)", at, left, right)
	}
}

func TestDoubleR2ApproxDoublesMaxSlip(t *testing.T) {
	in := example()
	base, err := Compute(in)
	if err != nil {
		t.Fatalf("Compute(base): %v", err)
	}
	doubled := in
	doubled.R2 = 2 * in.R2
	res2, err := Compute(doubled)
	if err != nil {
		t.Fatalf("Compute(r2*2): %v", err)
	}
	// The Thevenin denominator is independent of r2, so doubling r2 doubles
	// the maximum-torque slip exactly.
	if d := math.Abs(res2.MaxSlip - 2*base.MaxSlip); d > 1e-12 {
		t.Errorf("sm(r2*2) = %v, want 2*sm(r2) = %v", res2.MaxSlip, 2*base.MaxSlip)
	}
}

func TestTorqueSweepSpansInterval(t *testing.T) {
	in := example()
	points := TorqueSweep(in, 0.01, 1.0, 10)
	if len(points) != 10 {
		t.Fatalf("TorqueSweep returned %d points, want 10", len(points))
	}
	if first, last := points[0][0], points[len(points)-1][0]; first != 0.01 || last != 1.0 {
		t.Errorf("sweep ends at (%v, %v), want (0.01, 1.0)", first, last)
	}
	for _, p := range points {
		if p[1] < 0 {
			t.Errorf("sweep torque at s=%v is negative: %v", p[0], p[1])
		}
	}
}
