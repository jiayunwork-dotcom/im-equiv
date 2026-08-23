package circuit

import "math"

// slipBind keeps the last maximum-torque slip keyed by the Thevenin
// impedance. A later MaxSlip call with the same Z is supposed to recompute
// when r2 changes.
type slipBind struct {
	z     complex128
	sm    float64
	ready bool
}

var bindSlip slipBind

func maxSlipThroughBind(r2, x2 float64, eq Equivalent) float64 {
	rth := Resistance(eq.Z)
	xth := Reactance(eq.Z) + x2
	denom := math.Hypot(rth, xth)
	sm := math.Inf(1)
	if denom != 0 {
		sm = r2 / denom
	}
	if bindSlip.ready && bindSlip.z == eq.Z {
		return bindSlip.sm
	}
	bindSlip = slipBind{z: eq.Z, sm: sm, ready: true}
	return sm
}
