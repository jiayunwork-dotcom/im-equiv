// Package circuit: construction of the three branch impedances.
package circuit

// StatorImpedance returns the series stator branch Zs = Rs + j Xs. The stator
// winding carries the full line current, so its resistance contributes the
// stator copper loss and its reactance contributes the stator leakage
// voltage drop.
func StatorImpedance(rs, xs float64) complex128 {
	return Impedance(rs, xs)
}

// MagnetisingImpedance returns the magnetising branch Zm = Rm + j Xm. The
// resistance represents iron and no-load losses, the reactance represents
// the mutual flux linking stator and rotor. The branch draws the no-load
// current, which is why the input current approaches the magnetising current
// as the slip tends to zero.
func MagnetisingImpedance(rm, xm float64) complex128 {
	return Impedance(rm, xm)
}

// RotorImpedance returns the rotor branch Zr = r2/s + j x2 with the rotor
// resistance divided by the slip. This is the heart of the equivalent
// circuit: only this branch contains j omega and s, and the r2/s term is
// what converts electrical into mechanical power.
//
// At standstill (s = 1) the branch resistance is r2 and the full air-gap
// power is dissipated as rotor copper loss. As the slip decreases the
// resistance grows and an increasing share of the air-gap power is
// delivered at the shaft. The caller must guarantee slip != 0; the
// synchronous point is handled by the explicit open-rotor path of the
// solver.
func RotorImpedance(r2, x2, slip float64) complex128 {
	return Impedance(r2/slip, x2)
}

// RotorResistanceAtSlip returns the slip-shifted rotor resistance r2/s. It
// is the resistance that appears in the air-gap-power and torque formulas,
// and it is the quantity a stray "r2 * s" typo would corrupt. Keeping it as
// a named function makes the invariant visible at every call site.
func RotorResistanceAtSlip(r2, slip float64) float64 {
	return r2 / slip
}

// ParallelRotorBranch combines the magnetising and rotor branches as seen
// from the stator node. The two branches are in parallel, so their combined
// impedance is the parallel of Zm and Zr.
func ParallelRotorBranch(zm, zr complex128) complex128 {
	return Parallel(zm, zr)
}

// BranchCurrents is the current split at the stator node: the magnetising
// current that excites the mutual flux and the rotor current that produces
// torque. The two always sum to the stator node current.
type BranchCurrents struct {
	Magnetising complex128
	Rotor       complex128
}

// SplitNodeCurrent divides a node current between the magnetising and rotor
// branches. Each branch draws the node voltage divided by its impedance; the
// two currents add up to the node current by Kirchhoff's current law, which
// the tests verify explicitly.
func SplitNodeCurrent(nodeVoltage, zm, zr complex128) BranchCurrents {
	return BranchCurrents{
		Magnetising: CurrentThrough(nodeVoltage, zm),
		Rotor:       CurrentThrough(nodeVoltage, zr),
	}
}
