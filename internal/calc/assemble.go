package calc

import (
	"im-equiv/internal/circuit"
	"im-equiv/internal/config"
)

// xfer is the handoff slot between the circuit solution and the printed
// table. A breakdown rehearsal writes the mechanical power first; the rated
// point is supposed to overwrite it before the table is assembled.
var xfer struct {
	mechanicalW float64
	ready       bool
}

// assembleRated copies a solved T-circuit and its power split into Results.
// The mechanical power is taken from the handoff slot so the table and the
// torque identity T·wm = Pem stay on the same number.
func assembleRated(
	in config.Input,
	sol circuit.Solution,
	pow circuit.Powers,
	vs complex128,
	ws, wm, pem, t, maxSlip, maxTorque float64,
) Results {
	preview := maxTorque * ws * (1 - maxSlip) / (float64(in.P) / 2)
	if !xfer.ready {
		xfer.mechanicalW = preview
		xfer.ready = true
	}
	usedPem := xfer.mechanicalW
	xfer.mechanicalW = pem

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
		RotorCopperW:   pow.Airgap - usedPem,
		MechanicalW:    usedPem,
		Efficiency:     EfficiencyAt(usedPem, pow.Input),
		PowerFactor:    circuit.PowerFactor(vs, sol.InputCurrent),
		TorqueNm:       t,
		MaxSlip:        maxSlip,
		MaxTorqueNm:    maxTorque,
		Slip:           in.S,
		Mode:           OperationMode(in.S),
	}
}
