// Package circuit: three-phase power decomposition.
//
// Power in a balanced three-phase system is three times the per-phase
// power. The decomposition below follows the physical path of the energy:
// the supply delivers Pin, the stator winding and the iron circuit consume
// Pcu1 and Pfe before the gap, the remaining power crosses the gap as Pag,
// and the rotor consumes Pcu2 of it, leaving Pem for the shaft. Package
// calc turns that final difference into the torque and mechanical power.
package circuit

// Powers carries the power split across the machine for one operating point.
//
//	Input        three-phase power drawn from the supply, 3 Re(Vph I1*)
//	StatorCopper resistive loss of the stator winding, 3 I1^2 Rs
//	Iron         magnetising loss of the iron circuit, 3 I0^2 Rm
//	Airgap       power crossing the air gap, 3 I2^2 (r2/s)
//	RotorCopper  resistive loss of the referred rotor, 3 I2^2 r2
//
// The air-gap power minus the rotor copper loss is delivered at the shaft;
// package calc splits it into mechanical power and torque. The stator copper
// and iron losses are dissipated before the gap, so they never appear in
// the mechanical power.
type Powers struct {
	Input        float64
	StatorCopper float64
	Iron         float64
	Airgap       float64
	RotorCopper  float64
}

// ThreePhaseResistivePower returns 3 I^2 R for a balanced three-phase
// system in which a current of magnitude i flows through a resistance R per
// phase.
func ThreePhaseResistivePower(currentMag, resistance float64) float64 {
	return 3 * currentMag * currentMag * resistance
}

// AirgapPower is the power crossing the air gap. It equals the power that
// the slip-shifted rotor resistance r2/s dissipates, which by the rotating
// field argument is the power converted through the gap:
//
//	Pag = 3 I2^2 (r2 / s)
//
// At standstill (s = 1) it equals the rotor copper loss; as s shrinks the
// mechanical share grows. At the synchronous point the branch current I2 is
// exactly zero while r2/s is singular, so the zero-slip case is pinned to
// zero explicitly to keep the product finite.
func AirgapPower(rotorCurrentMag, r2, slip float64) float64 {
	if slip == 0 {
		return 0
	}
	return ThreePhaseResistivePower(rotorCurrentMag, r2/slip)
}

// RotorCopperLoss is the resistive loss in the referred rotor resistance
// r2. The difference Pag - Pcu2 is the electromagnetic power delivered to
// the shaft.
func RotorCopperLoss(rotorCurrentMag, r2 float64) float64 {
	return ThreePhaseResistivePower(rotorCurrentMag, r2)
}

// StatorCopperLoss is the resistive loss of the stator winding, carried by
// the full input current.
func StatorCopperLoss(inputCurrentMag, rs float64) float64 {
	return ThreePhaseResistivePower(inputCurrentMag, rs)
}

// IronLoss is the loss dissipated in the magnetising resistance Rm, carried
// by the magnetising current I0. It represents the iron (hysteresis and
// eddy-current) loss of the core and is present at every operating point,
// including no load.
func IronLoss(magnetisingCurrentMag, rm float64) float64 {
	return ThreePhaseResistivePower(magnetisingCurrentMag, rm)
}

// TotalLosses sums the stator copper, iron and rotor copper losses. It is
// the quantity that must equal the difference between the input power and
// the mechanical power, which is the macroscopic energy balance of the
// machine.
func (p Powers) TotalLosses() float64 {
	return p.StatorCopper + p.Iron + p.RotorCopper
}

// PowerFactor is the cosine of the angle between the phase voltage and the
// line current, that is the real part of the ratio of the apparent power to
// its magnitude. A purely resistive load would give exactly one.
func PowerFactor(phaseVoltage, lineCurrent complex128) float64 {
	app := ApparentPower(phaseVoltage, lineCurrent)
	if Mag(app) == 0 {
		return 0
	}
	return Real(app) / Mag(app)
}

// Powers computes the full power decomposition of a solved circuit. The slip
// argument is required only to scale the rotor resistance into the air-gap
// power; the synchronous point must use the open-rotor solution, for which
// the rotor current is zero and the air-gap power vanishes.
func (s Solution) Powers(phaseVoltage, statorImpedance complex128, rm, r2, slip float64) Powers {
	i1 := Mag(s.InputCurrent)
	i2 := Mag(s.RotorCurrent)
	i0 := Mag(s.MagnetisingCurrent)
	return Powers{
		Input:        s.InputPower(phaseVoltage),
		StatorCopper: StatorCopperLoss(i1, Resistance(statorImpedance)),
		Iron:         IronLoss(i0, rm),
		Airgap:       AirgapPower(i2, r2, slip),
		RotorCopper:  RotorCopperLoss(i2, r2),
	}
}
