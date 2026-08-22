// Package calc: ratings and derived quantity tests.
package calc

import (
	"math"
	"testing"
)

func TestPullOutRatioTypicalRange(t *testing.T) {
	in := example()
	res, err := Compute(in)
	if err != nil {
		t.Fatalf("Compute: %v", err)
	}
	ratio := PullOutRatio(res.MaxTorqueNm, res.TorqueNm)
	if ratio < 2 || ratio > 4 {
		t.Errorf("pull-out ratio = %v, want typical 2-3", ratio)
	}
}

func TestPullOutRatioZeroTorqueInfinite(t *testing.T) {
	if r := PullOutRatio(50, 0); !math.IsInf(r, 1) {
		t.Errorf("PullOutRatio(50, 0) = %v, want +Inf", r)
	}
}

func TestSlipAtSpeed(t *testing.T) {
	// At 1455 rpm for a four-pole 50 Hz machine the slip is 0.03.
	if got := SlipAtSpeed(50, 4, 1455); math.Abs(got-0.03) > 1e-9 {
		t.Errorf("SlipAtSpeed(50, 4, 1455) = %v, want 0.03", got)
	}
}

func TestStartingTorquePositive(t *testing.T) {
	in := example()
	if torque := StartingTorqueNm(in); torque <= 0 {
		t.Errorf("starting torque = %v, want positive", torque)
	}
}

func TestStartingCurrentRatioAboveOne(t *testing.T) {
	in := example()
	if r := StartingCurrentRatio(in); r <= 1 {
		t.Errorf("starting current ratio = %v, want above 1", r)
	}
}

func TestBreakdownTorqueMatchesCompute(t *testing.T) {
	in := example()
	res, err := Compute(in)
	if err != nil {
		t.Fatalf("Compute: %v", err)
	}
	if got := BreakdownTorqueNm(in); math.Abs(got-res.MaxTorqueNm) > 1e-9 {
		t.Errorf("BreakdownTorqueNm = %v, want %v", got, res.MaxTorqueNm)
	}
}
