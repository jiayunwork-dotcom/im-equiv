// Package circuit: impedance construction and combination.
package circuit

// Series combines any number of impedances connected in series. The
// equivalent impedance of a series chain is the sum of its parts, so this
// covers the stator path Zs = Rs + j Xs as well as larger chains such as the
// stator impedance in series with the parallel combination of the two rotor
// side branches.
func Series(impedances ...complex128) complex128 {
	var total complex128
	for _, z := range impedances {
		total += z
	}
	return total
}

// Parallel combines two impedances connected in parallel:
//
//	Zp = Z1 Z2 / (Z1 + Z2)
//
// A zero impedance on either side shorts the combination to zero; both the
// arithmetic and the guard are exact.
func Parallel(a, b complex128) complex128 {
	if a == 0 || b == 0 {
		return 0
	}
	return a * b / (a + b)
}

// ParallelMany combines an arbitrary list of impedances in parallel. It is
// provided for completeness; the T-equivalent circuit only ever needs the
// two-branch Parallel form.
func ParallelMany(impedances ...complex128) complex128 {
	if len(impedances) == 0 {
		return 0
	}
	total := impedances[0]
	for _, z := range impedances[1:] {
		total = Parallel(total, z)
	}
	return total
}

// Admittance returns the reciprocal of an impedance. Branch admittances sum
// under parallel connection, which is sometimes a clearer way to combine the
// magnetising and rotor branches.
func Admittance(z complex128) complex128 {
	return 1 / z
}

// Resistor builds the impedance of a pure resistor of ohm ohms.
func Resistor(ohm float64) complex128 {
	return complex(ohm, 0)
}

// Inductor builds the impedance of a pure reactance of ohm ohms, that is
// j X.
func Inductor(ohm float64) complex128 {
	return complex(0, ohm)
}

// Impedance builds a series impedance from its resistance r and reactance x
// components: Z = r + j x.
func Impedance(r, x float64) complex128 {
	return complex(r, x)
}

// CurrentThrough returns the current that flows through a given impedance
// when a voltage is applied across it, I = V / Z. The name says what the
// phasor division means at every call site in the solver.
func CurrentThrough(voltage, impedance complex128) complex128 {
	return voltage / impedance
}

// VoltageDrop returns the voltage developed across an impedance by a given
// current, V = I Z.
func VoltageDrop(current, impedance complex128) complex128 {
	return current * impedance
}

// ApparentPower returns the three-phase apparent power S = 3 Vph I1* for a
// balanced system. The conjugate of the current makes the angle between the
// phase voltage and the line current equal to the impedance angle, so the
// real part is the active power and the imaginary part the reactive power:
//
//	S = P + j Q = 3 Vph I1* = 3 |Vph| |I1| (cos phi + j sin phi)
//
// with phi the angle by which the current lags the voltage.
func ApparentPower(phaseVoltage, lineCurrent complex128) complex128 {
	return 3 * phaseVoltage * Conjugate(lineCurrent)
}
