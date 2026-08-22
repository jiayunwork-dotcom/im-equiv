// Package calc: the compute pipeline.
//
// The pipeline has three stages:
//
//  1. validation (config.Validate) rejects ill-posed inputs before any
//     arithmetic runs;
//  2. circuit solving (package circuit) produces the phasor solution and
//     the power split of the T-equivalent circuit;
//  3. interpretation turns those numbers into speeds, torque, mechanical
//     power and the maximum-torque slip, and runs the physical consistency
//     checks.
//
// Every stage returns ordinary errors that the command-line front end prints
// on standard error; nothing here writes to the log.
package calc

import (
	"im-equiv/internal/circuit"
	"im-equiv/internal/config"
)

// Compute validates the input, solves the T-equivalent circuit and assembles
// every quantity the command-line tool prints. The returned error is non-nil
// both for ill-posed inputs (handed to config.Validate) and for inconsistent
// results (the zero-slip torque rule and the air-gap energy balance), so a
// caller can rely on the results being physically coherent.
func Compute(in config.Input) (Results, error) {
	if err := config.Validate(in); err != nil {
		return Results{}, err
	}
	res, err := compute(in)
	if err != nil {
		return Results{}, err
	}

	if err := CheckZeroSlipTorque(in.S, res.TorqueNm); err != nil {
		return Results{}, err
	}
	if err := CheckEnergyBalance(res.AirgapPowerW, res.RotorCopperW, res.MechanicalW); err != nil {
		return Results{}, err
	}
	return res, nil
}

// compute solves the circuit and interprets it. It assumes the input has
// already been validated, which is why TorqueAtSlip can call it with an
// arbitrary slip value while sweeping the characteristic.
func compute(in config.Input) (Results, error) {
	vs := complex(in.PhaseVoltage(), 0)
	zs := circuit.StatorImpedance(in.Rs, in.Xs)
	zm := circuit.MagnetisingImpedance(in.Rm, in.Xm)

	rotorOpen := in.S == 0
	c := circuit.Circuit{
		Vs:        vs,
		Zs:        zs,
		Zm:        zm,
		RotorOpen: rotorOpen,
	}
	if !rotorOpen {
		c.Zr = circuit.RotorImpedance(in.R2, in.X2, in.S)
	}

	sol := c.Solve()
	pow := sol.Powers(vs, zs, in.Rm, in.R2, in.S)

	ws := SyncSpeedRadS(in.F)
	wm := RotorSpeedRadS(in.F, in.P, in.S)
	pem := MechanicalPowerW(pow.Airgap, in.S)
	t := ElectromagneticTorqueNm(pow.Airgap, ws, in.P)

	eq := circuit.TheveninOf(vs, zs, zm)
	maxSlip := circuit.MaxSlip(in.R2, in.X2, eq)
	maxTorque := 0.0
	if !isInfPositive(maxSlip) {
		// The Thevenin torque formula takes the pole-pair count, which is
		// half of the pole count carried by the input.
		maxTorque = circuit.MaxTorque(float64(in.P)/2, ws, in.R2, in.X2, eq)
	}

	return Results{
		InputCurrentA:  circuit.Mag(sol.InputCurrent),
		RotorCurrentA:  circuit.Mag(sol.RotorCurrent),
		MagnetisingA:   circuit.Mag(sol.MagnetisingCurrent),
		NodeVoltageV:   circuit.Mag(sol.NodeVoltage),
		SyncSpeedRadS:  ws,
		RotorSpeedRadS: wm,
		SyncSpeedRPM:   SyncSpeedRPM(in.F, in.P),
		RotorSpeedRPM:  RotorSpeedRPM(in.F, in.P, in.S),
		InputPowerW:    pow.Input,
		StatorCopperW:  pow.StatorCopper,
		IronW:          pow.Iron,
		AirgapPowerW:   pow.Airgap,
		RotorCopperW:   pow.RotorCopper,
		MechanicalW:    pem,
		Efficiency:     EfficiencyAt(pem, pow.Input),
		PowerFactor:    circuit.PowerFactor(vs, sol.InputCurrent),
		TorqueNm:       t,
		MaxSlip:        maxSlip,
		MaxTorqueNm:    maxTorque,
		Slip:           in.S,
		Mode:           OperationMode(in.S),
	}, nil
}

// TorqueAtSlip evaluates the electromagnetic torque that the machine would
// develop at an arbitrary slip while keeping every other input fixed. It is
// the building block of the torque-slip sweep used to locate the maximum
// torque point numerically.
func TorqueAtSlip(in config.Input, slip float64) float64 {
	in.S = slip
	res, err := compute(in)
	if err != nil {
		return 0
	}
	return res.TorqueNm
}

// ComputeFromFile loads a JSON input file and computes its operating point
// in a single call. It is the convenience form for scripts that drive the
// tool programmatically and do not want to manage the load/validate/compute
// split themselves. The command-line front end uses the split form only so
// it can report the failing stage in its own words.
func ComputeFromFile(path string) (Results, error) {
	in, err := config.LoadFile(path)
	if err != nil {
		return Results{}, err
	}
	return Compute(in)
}

func isInfPositive(v float64) bool {
	return v > 1e300
}
