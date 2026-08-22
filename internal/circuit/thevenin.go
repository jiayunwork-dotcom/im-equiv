// Package circuit: Thevenin equivalent and the closed-form torque formula.
//
// The Thevenin view is the classical shortcut of the T-equivalent circuit:
// seen from the rotor terminals, the stator and magnetising branches reduce
// to a fixed source Vth in series with Zth. The rotor branch then sees a
// simple series circuit, which turns the torque expression into a single
// rational function of the slip and yields closed forms for the
// maximum-torque slip and the breakdown torque.
package circuit

import (
	"math"
)

// Equivalent is the Thevenin equivalent seen from the rotor terminals of the
// T-equivalent circuit: the stator branch in series with the magnetising
// branch, with the rotor branch itself removed.
//
//	Zth = Zs || Zm
//	Vth = Vph Zm / (Zs + Zm)
//
// The rotor branch then sees a fixed source Vth in series with Zth, which
// reduces the torque expression to a single rational function of the slip.
type Equivalent struct {
	V complex128 // open-circuit voltage seen by the rotor
	Z complex128 // equivalent impedance seen by the rotor
}

// TheveninOf computes the equivalent seen by the rotor for a phase voltage
// vs, a stator branch zs and a magnetising branch zm.
func TheveninOf(vs, zs, zm complex128) Equivalent {
	return Equivalent{
		V: vs * zm / (zs + zm),
		Z: Parallel(zs, zm),
	}
}

// MaxSlip returns the slip at which the machine develops its maximum
// electromagnetic torque, using the standard closed-form result
//
//	sm = r2 / sqrt(Rth^2 + (Xth + x2)^2)
//
// with Rth + j Xth the Thevenin impedance and x2 the referred rotor leakage
// reactance. Doubling r2 approximately doubles sm, which is one of the
// cross-check rules the tests exercise. If the denominator vanishes the
// torque is monotonic and there is no finite maximum, so +Inf is returned.
func MaxSlip(r2, x2 float64, eq Equivalent) float64 {
	rth := Resistance(eq.Z)
	xth := Reactance(eq.Z) + x2
	denom := math.Hypot(rth, xth)
	if denom == 0 {
		return math.Inf(1)
	}
	return r2 / denom
}

// TorqueFromThevenin evaluates the electromagnetic torque from the Thevenin
// equivalent with the standard formula
//
//	T = (3 p / ws) * Vth^2 (r2/s) / ((Rth + r2/s)^2 + (Xth + x2)^2)
//
// where ws is the electrical synchronous speed in rad/s and p the number of
// pole pairs. The slip must be non-zero; near s = 0 the result correctly
// tends to zero as r2/s dominates the denominator.
func TorqueFromThevenin(polePairs float64, omegaSyncRadS float64, r2, x2 float64, eq Equivalent, slip float64) float64 {
	r := Resistance(eq.Z) + r2/slip
	x := Reactance(eq.Z) + x2
	return 3 * polePairs / omegaSyncRadS * MagSq(eq.V) * (r2 / slip) / (r*r + x*x)
}

// MaxTorque returns the peak torque developed at the maximum-torque slip,
// evaluated on the Thevenin curve. An infinite MaxSlip yields zero, meaning
// the torque characteristic has no interior peak.
func MaxTorque(polePairs float64, omegaSyncRadS float64, r2, x2 float64, eq Equivalent) float64 {
	sm := MaxSlip(r2, x2, eq)
	if math.IsInf(sm, 1) {
		return 0
	}
	return TorqueFromThevenin(polePairs, omegaSyncRadS, r2, x2, eq, sm)
}

// BreakdownTorque returns the closed-form peak torque without first locating
// the maximum-torque slip:
//
//	Tmax = (3 p' / (2 ws)) Vth^2 / (Rth + sqrt(Rth^2 + (Xth + x2)^2))
//
// It is algebraically identical to evaluating TorqueFromThevenin at sm and
// is provided as the standard "pull-out torque" expression found in every
// induction-machine textbook.
func BreakdownTorque(polePairs float64, omegaSyncRadS float64, r2, x2 float64, eq Equivalent) float64 {
	rth := Resistance(eq.Z)
	xth := Reactance(eq.Z) + x2
	denom := rth + math.Hypot(rth, xth)
	return 3 * polePairs / (2 * omegaSyncRadS) * MagSq(eq.V) / denom
}
