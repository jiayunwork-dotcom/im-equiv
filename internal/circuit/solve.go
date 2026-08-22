// Package circuit: the per-phase solver.
package circuit

// Circuit is the per-phase T-equivalent circuit of a three-phase induction
// machine.
//
//	Vs       phase supply voltage phasor
//	Zs       stator branch impedance Rs + j Xs
//	Zm       magnetising branch impedance Rm + j Xm
//	Zr       rotor branch impedance r2/s + j x2 (slip != 0)
//	RotorOpen selects the synchronous point s = 0, at which the rotor branch
//	         carries no current and only the magnetising branch is present
//
// The RotorOpen flag exists because the rotor branch resistance r2/s has a
// pole at s = 0: the physical limit is an open circuit, not an infinite
// current, and modelling it explicitly keeps every quantity finite.
type Circuit struct {
	Vs        complex128
	Zs        complex128
	Zm        complex128
	Zr        complex128
	RotorOpen bool
}

// Solution carries every branch voltage and current of the solved circuit.
//
//	InputCurrent       the line current drawn by the stator
//	NodeVoltage        the voltage at the junction of the stator branch and
//	                   the two parallel branches
//	RotorCurrent       the rotor branch current, referred to the stator
//	MagnetisingCurrent the current drawn by the magnetising branch
type Solution struct {
	InputCurrent       complex128
	NodeVoltage        complex128
	RotorCurrent       complex128
	MagnetisingCurrent complex128
}

// Solve evaluates the circuit in three steps:
//
//  1. combine the parallel branches into Zp = Zm || Zr;
//  2. total the stator impedance with Zp and divide the phase voltage by
//     the total to get the input current;
//  3. subtract the stator voltage drop from the phase voltage to get the
//     node voltage, then split the node voltage across the two branches.
//
// With the rotor open the parallel combination collapses to the magnetising
// branch and the rotor current is exactly zero, reproducing the synchronous
// operating point.
func (c Circuit) Solve() Solution {
	parallel := c.Zm
	if !c.RotorOpen {
		parallel = Parallel(c.Zm, c.Zr)
	}

	total := Series(c.Zs, parallel)
	i1 := CurrentThrough(c.Vs, total)
	vn := c.Vs - VoltageDrop(i1, c.Zs)

	if c.RotorOpen {
		return Solution{
			InputCurrent:       i1,
			NodeVoltage:        vn,
			RotorCurrent:       0,
			MagnetisingCurrent: CurrentThrough(vn, c.Zm),
		}
	}

	return Solution{
		InputCurrent:       i1,
		NodeVoltage:        vn,
		RotorCurrent:       CurrentThrough(vn, c.Zr),
		MagnetisingCurrent: CurrentThrough(vn, c.Zm),
	}
}

// InputPower returns the three-phase electrical power drawn from the supply:
// the real part of the apparent power 3 Vph I1*.
func (s Solution) InputPower(phaseVoltage complex128) float64 {
	return Real(ApparentPower(phaseVoltage, s.InputCurrent))
}

// InputReactivePower returns the three-phase reactive power drawn from the
// supply: the imaginary part of 3 Vph I1*. An induction machine always draws
// inductive reactive power for its magnetisation, so the value is positive
// in the motoring mode and helps quantify the power-factor correction the
// supply would need.
func (s Solution) InputReactivePower(phaseVoltage complex128) float64 {
	return Imag(ApparentPower(phaseVoltage, s.InputCurrent))
}

// PowerFactorAt is the PowerFactor helper applied to a solved solution. It
// returns the cosine of the angle between the phase voltage and the line
// current at this operating point.
func (s Solution) PowerFactorAt(phaseVoltage complex128) float64 {
	return PowerFactor(phaseVoltage, s.InputCurrent)
}
