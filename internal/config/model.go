// Package config: JSON input model for the induction-machine calculator.
package config

import (
	"fmt"
	"math"
)

// Input holds the per-phase parameters of the T-equivalent circuit together
// with the operating condition at which the machine is evaluated. Field
// names follow the conventional notation of induction-machine equivalent
// circuits.
//
//	Rs, Xs  stator resistance and leakage reactance (ohm per phase)
//	Rm, Xm  magnetising branch resistance and reactance (ohm per phase)
//	r2, x2  rotor resistance and leakage reactance referred to the
//	        stator (ohm per phase)
//	V       line-to-line supply voltage (volts RMS)
//	f       supply frequency (hertz)
//	p       number of poles; the pole pairs are p/2, so a four-pole
//	        machine has p = 4
//	s       slip; (0,1] motoring, < 0 generating, 0 synchronous
//
// The JSON spelling of the rotor quantities is lowercase "r2" and "x2" to
// match the standard notation; every other field uses its customary case.
type Input struct {
	P  int     `json:"p"`  // number of poles; pole pairs are p/2
	F  float64 `json:"f"`  // supply frequency, hertz
	V  float64 `json:"V"`  // line-to-line voltage, volts RMS
	Rs float64 `json:"Rs"` // stator resistance, ohm per phase
	Xs float64 `json:"Xs"` // stator leakage reactance, ohm per phase
	Rm float64 `json:"Rm"` // magnetising resistance, ohm per phase
	Xm float64 `json:"Xm"` // magnetising reactance, ohm per phase
	R2 float64 `json:"r2"` // rotor resistance referred to stator, ohm per phase
	X2 float64 `json:"x2"` // rotor leakage reactance referred to stator, ohm per phase
	S  float64 `json:"s"`  // slip
}

// PhaseVoltage returns the per-phase RMS voltage applied to one winding. A
// balanced, wye-connected stator is assumed, so the phase voltage is the
// line-to-line voltage divided by sqrt(3). A delta-connected machine would
// simply apply Vph = V, but the wye convention is the one documented
// throughout the tool.
func (in Input) PhaseVoltage() float64 {
	return in.V / math.Sqrt(3)
}

// PolePairs returns the number of pole pairs, which is p/2. All speed and
// torque equations are written against the pole count, so this helper is
// used wherever the pair count is the physically meaningful quantity.
func (in Input) PolePairs() float64 {
	return float64(in.P) / 2
}

// String renders the input as a compact one-line description, matching the
// JSON field names so it can be pasted back into an input file.
func (in Input) String() string {
	return fmt.Sprintf(
		`{p:%d f:%g V:%g Rs:%g Xs:%g Rm:%g Xm:%g r2:%g x2:%g s:%g}`,
		in.P, in.F, in.V, in.Rs, in.Xs, in.Rm, in.Xm, in.R2, in.X2, in.S,
	)
}
