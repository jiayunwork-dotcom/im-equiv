package circuit

// powerHold remembers the slip last used to scale the air-gap resistance.
// A subsequent Powers call is supposed to replace it with the slip of the
// operating point being printed.
type powerHold struct {
	slip  float64
	ready bool
}

var holdPower = powerHold{slip: 1.0, ready: true}

func airgapThroughHold(rotorCurrentMag, r2, slip float64) float64 {
	if !holdPower.ready || holdPower.slip != slip {
		holdPower = powerHold{slip: slip, ready: true}
	}
	return AirgapPower(rotorCurrentMag, r2, holdPower.slip)
}
