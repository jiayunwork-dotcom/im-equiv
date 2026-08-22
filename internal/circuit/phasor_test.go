// Package circuit: phasor and impedance arithmetic tests.
package circuit

import (
	"math"
	"testing"
)

func TestPhasorMagnitudeAndPhase(t *testing.T) {
	z := complex(3, 4)
	if got, want := Mag(z), 5.0; math.Abs(got-want) > 1e-12 {
		t.Errorf("Mag(3+4j) = %v, want %v", got, want)
	}
	if got, want := MagSq(z), 25.0; math.Abs(got-want) > 1e-12 {
		t.Errorf("MagSq(3+4j) = %v, want %v", got, want)
	}
	wantPhase := math.Atan2(4, 3)
	if got := Phase(z); math.Abs(got-wantPhase) > 1e-12 {
		t.Errorf("Phase(3+4j) = %v, want %v", got, wantPhase)
	}
}

func TestFromPolarDegrees(t *testing.T) {
	z := FromPolarDegrees(10, 90)
	if got, want := Real(z), 0.0; math.Abs(got-want) > 1e-12 {
		t.Errorf("FromPolarDegrees(10,90).real = %v, want %v", got, want)
	}
	if got, want := Imag(z), 10.0; math.Abs(got-want) > 1e-12 {
		t.Errorf("FromPolarDegrees(10,90).imag = %v, want %v", got, want)
	}
}

func TestSeriesSum(t *testing.T) {
	total := Series(Impedance(1, 2), Impedance(3, -1))
	if got, want := Resistance(total), 4.0; math.Abs(got-want) > 1e-12 {
		t.Errorf("Series resistance = %v, want %v", got, want)
	}
	if got, want := Reactance(total), 1.0; math.Abs(got-want) > 1e-12 {
		t.Errorf("Series reactance = %v, want %v", got, want)
	}
}

func TestParallelOpenCircuit(t *testing.T) {
	// A branch with infinite impedance must leave the other branch unchanged.
	z := Parallel(Impedance(4, 3), complex(math.Inf(1), 0))
	if got, want := Resistance(z), 4.0; math.Abs(got-want) > 1e-9 {
		t.Errorf("Parallel(4+3j, inf) resistance = %v, want %v", got, want)
	}
}

func TestParallelShortCircuit(t *testing.T) {
	if z := Parallel(Impedance(4, 3), 0); z != 0 {
		t.Errorf("Parallel(4+3j, 0) = %v, want 0", z)
	}
}

func TestApparentPowerAngle(t *testing.T) {
	v := complex(220, 0)
	// Current lagging the voltage by 30 degrees: power factor 0.866.
	i := FromPolarDegrees(10, -30)
	app := ApparentPower(v, i)
	expected := 3 * 220 * 10 * math.Cos(math.Pi/6)
	if got := Real(app); math.Abs(got-expected) > 1e-9 {
		t.Errorf("Real(ApparentPower) = %v, want %v", got, expected)
	}
}
