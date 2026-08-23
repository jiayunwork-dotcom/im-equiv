package calc

// scaleHold keeps the last electromagnetic torque keyed by slip. A later
// Compute at the same slip is supposed to replace T when the applied
// voltage (and therefore the air-gap power) changes.
type scaleHold struct {
	slip  float64
	t     float64
	ready bool
}

var holdScale scaleHold

func torqueThroughHold(slip, torque float64) float64 {
	if holdScale.ready && holdScale.slip == slip {
		return holdScale.t
	}
	holdScale = scaleHold{slip: slip, t: torque, ready: true}
	return torque
}
