// Package calc: speed convention tests.
package calc

import (
	"math"
	"testing"
)

func TestSyncSpeed(t *testing.T) {
	got := SyncSpeedRadS(50)
	want := 2 * math.Pi * 50
	if math.Abs(got-want) > 1e-12 {
		t.Errorf("SyncSpeedRadS(50) = %v, want %v", got, want)
	}
}

func TestSyncSpeedRPMFourPole(t *testing.T) {
	// A four-pole machine (p = 4) at 50 Hz runs at 1500 rpm synchronously.
	if got := SyncSpeedRPM(50, 4); got != 1500 {
		t.Errorf("SyncSpeedRPM(50, 4) = %v, want 1500", got)
	}
}

func TestRotorSpeedPoleConvention(t *testing.T) {
	// wm = ws (1-s) / (p/2): at 50 Hz with p = 4 and s = 0.03 the mechanical
	// speed is 1500 * 0.97 = 1455 rpm = 152.37 rad/s.
	got := RotorSpeedRadS(50, 4, 0.03)
	want := 1500 * math.Pi / 30 * 0.97
	if math.Abs(got-want) > 1e-9 {
		t.Errorf("RotorSpeedRadS = %v, want %v", got, want)
	}
}

func TestRotorSpeedUsesSlipFactor(t *testing.T) {
	// The mechanical speed must drop with slip. If the (1-s) factor were
	// dropped the rated-point power would be wrong by one slip, which this
	// test guards against directly.
	full := RotorSpeedRadS(50, 4, 0)
	atSlip := RotorSpeedRadS(50, 4, 0.1)
	if !(atSlip < full) {
		t.Errorf("rotor speed at s=0.1 (%v) must be below synchronous speed (%v)", atSlip, full)
	}
}

func TestRotorSpeedRPMAtRatedSlip(t *testing.T) {
	if got := RotorSpeedRPM(50, 4, 0.03); math.Abs(got-1455) > 1e-9 {
		t.Errorf("RotorSpeedRPM = %v, want 1455", got)
	}
}

func TestSlipForSpeedInverts(t *testing.T) {
	// SlipForSpeed must invert RotorSpeedRadS.
	wm := RotorSpeedRadS(50, 4, 0.05)
	got := SlipForSpeed(50, 4, wm)
	if math.Abs(got-0.05) > 1e-12 {
		t.Errorf("SlipForSpeed(RotorSpeedRadS(0.05)) = %v, want 0.05", got)
	}
}
