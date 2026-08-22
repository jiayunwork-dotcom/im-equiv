// Package calc: operating-point classification.
//
// The slip alone decides which region of the torque-speed plane the machine
// occupies. The classification is used for the Mode field of Results, for
// the efficiency display (only the motoring mode shows a percentage) and as
// a quick sanity check of an input file before the full table is printed.
package calc

// OperationMode classifies the operating point by slip:
//
//	s = 0      synchronous operation; the rotor runs with the field and no
//	           torque is developed
//	0 < s <= 1 motoring; torque acts in the direction of rotation and power
//	           flows from the supply to the shaft
//	s > 1      plugging; the rotor is driven against a reversed field, so
//	           torque opposes rotation and the rotor dissipates the slip
//	           power
//	s < 0      generating; the rotor runs above synchronous speed, torque
//	           opposes rotation and power flows from the shaft to the supply
//
// The classification drives the Mode field of Results and is also useful on
// its own to sanity-check an input before printing the full table.
func OperationMode(slip float64) string {
	switch {
	case slip == 0:
		return "synchronous"
	case slip > 1:
		return "plugging"
	case slip > 0:
		return "motoring"
	default:
		return "generating"
	}
}

// IsMotoring reports whether the slip lies in the motoring band (0, 1].
func IsMotoring(slip float64) bool {
	return slip > 0 && slip <= 1
}

// IsGenerating reports whether the slip is negative, i.e. the machine
// operates as a generator above synchronous speed.
func IsGenerating(slip float64) bool {
	return slip < 0
}
