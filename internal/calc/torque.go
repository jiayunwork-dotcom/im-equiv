// Package calc: torque and mechanical power.
//
// Torque and mechanical power are the two quantities a designer reads off
// the machine. Both are derived from the air-gap power so that the
// synchronous point (s = 0) stays finite: at that point the air-gap power
// is zero, which forces both the mechanical power and the torque to zero,
// and the torque expression written against ws rather than wm avoids the
// 0/0 indeterminate form.
package calc

// ElectromagneticTorqueNm derives the electromagnetic torque from the
// air-gap power. By definition T = Pem / wm, and substituting Pem = Pag (1-s)
// and wm = ws (1-s) / (p/2) with p the pole count gives the slip-independent
// form
//
//	T = (p/2) Pag / ws
//
// Computing the torque from the air-gap power and the electrical synchronous
// speed keeps the expression finite at synchronous speed, where both Pem and
// wm vanish. It also makes the scaling rule explicit: T grows linearly with
// Pag, hence with the square of the applied voltage.
func ElectromagneticTorqueNm(airgapPowerW, omegaSyncRadS float64, poleCount int) float64 {
	return (float64(poleCount) / 2) * airgapPowerW / omegaSyncRadS
}

// MechanicalPowerW splits the air-gap power into mechanical power and rotor
// copper loss:
//
//	Pem = Pag (1 - s)
//
// At standstill (s = 1) the mechanical power is zero and the entire air-gap
// power is dissipated in the rotor. As s shrinks the mechanical share grows,
// which is the energy balance the rated-point example must satisfy.
func MechanicalPowerW(airgapPowerW, slip float64) float64 {
	return airgapPowerW * (1 - slip)
}

// EfficiencyAt returns Pem / Pin for the given mechanical and electrical
// powers. Both powers are signed: in the generating mode they are negative,
// so the ratio stays positive and slightly below one. A purely resistive
// no-loss machine would give exactly one.
func EfficiencyAt(mechanicalW, inputPowerW float64) float64 {
	if inputPowerW == 0 {
		return 0
	}
	return mechanicalW / inputPowerW
}

// RotorCopperShare returns the fraction of the air-gap power that ends up as
// rotor copper loss, which is exactly the slip:
//
//	Pcu2 / Pag = (3 I2^2 r2) / (3 I2^2 r2/s) = s
//
// At standstill the whole air-gap power is lost in the rotor; as the slip
// shrinks the mechanical share grows.
func RotorCopperShare(slip float64) float64 {
	return slip
}

// MechanicalShare returns the fraction of the air-gap power delivered at the
// shaft, which is 1 - s. Together with RotorCopperShare the two shares
// always sum to one, which is the energy split across the air gap.
func MechanicalShare(slip float64) float64 {
	return 1 - slip
}

// RotorI2RLoss is an alias-style view of the rotor copper loss formula used
// in loss breakdown tables: 3 I2^2 r2. It is provided so callers that
// present the loss budget do not have to remember which power function
// applies to which stage.
func RotorI2RLoss(rotorCurrentMag, r2 float64) float64 {
	return 3 * rotorCurrentMag * rotorCurrentMag * r2
}
