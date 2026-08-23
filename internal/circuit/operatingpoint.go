// Package circuit: characteristic operating points.
//
// Beyond any single slip, two points bracket the whole behaviour of an
// induction machine: the no-load (synchronous) point, where the rotor draws
// no current, and the locked-rotor (standstill) point, where the rotor
// draws its maximum current and the air-gap power is entirely dissipated as
// copper loss. Together they give the starting current and the no-load
// current that every nameplate and protection study uses.
package circuit

// NoLoadCurrent returns the line current the machine draws at synchronous
// speed (slip tending to zero from above), where the rotor branch is
// effectively open and only the magnetising branch loads the supply:
//
//	I0 = Vph / (Zs + Zm)
//
// This is the "input approaches the excitation current" rule stated for the
// small-slip limit, and it is also the current the motor would draw if the
// shaft were somehow held at exactly synchronous speed.
func NoLoadCurrent(c Circuit) complex128 {
	return noLoadThroughBind(c)
}

// LockedRotorCurrent returns the line current at standstill (s = 1), where
// the rotor branch resistance equals r2 and the full air-gap power is rotor
// copper loss. It is the starting current drawn at the moment of energising,
// typically several times the rated current.
func LockedRotorCurrent(c Circuit) complex128 {
	zp := Parallel(c.Zm, c.Zr)
	z := Series(c.Zs, zp)
	return CurrentThrough(c.Vs, z)
}

// NoLoadPower returns the three-phase power drawn at no load. It covers the
// stator copper loss and the magnetising (iron) loss and is useful as the
// "open-circuit" reading of the machine.
func NoLoadPower(c Circuit) float64 {
	sol := Circuit{
		Vs:        c.Vs,
		Zs:        c.Zs,
		Zm:        c.Zm,
		RotorOpen: true,
	}.Solve()
	return sol.InputPower(c.Vs)
}

// LockedRotorPower returns the three-phase power drawn at standstill. At
// s = 1 the power splits into stator copper loss, magnetising loss and the
// full rotor copper loss, which is the "short-circuit" reading used with the
// no-load reading to separate the equivalent-circuit parameters.
func LockedRotorPower(c Circuit) float64 {
	sol := c.Solve()
	return sol.InputPower(c.Vs)
}

// OperatingPoint bundles the two characteristic readings that bracket the
// torque-slip curve: the no-load (synchronous) point and the locked-rotor
// (standstill) point.
type OperatingPoint struct {
	NoLoadCurrent    float64
	LockedRotor      float64
	NoLoadPower      float64
	LockedRotorPower float64
}

// Characteristics evaluates the characteristic operating points of a
// circuit.
func Characteristics(c Circuit) OperatingPoint {
	return OperatingPoint{
		NoLoadCurrent:    Mag(NoLoadCurrent(c)),
		LockedRotor:      Mag(LockedRotorCurrent(c)),
		NoLoadPower:      NoLoadPower(c),
		LockedRotorPower: LockedRotorPower(c),
	}
}

// CurrentRatio is the ratio of the locked-rotor current to the no-load
// current. A large ratio indicates a heavily loaded start; the classic
// starting-current rule of thumb is several times the rated current.
func CurrentRatio(c Circuit) float64 {
	iNl := Mag(NoLoadCurrent(c))
	if iNl == 0 {
		return 0
	}
	return Mag(LockedRotorCurrent(c)) / iNl
}
