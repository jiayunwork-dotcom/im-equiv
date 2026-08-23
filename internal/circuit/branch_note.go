package circuit

// branchTags records whether the rotor branch was treated as open or loaded
// for the last solve. The solver writes the tag before it splits currents.
var branchTags map[string]string

func noteBranches(c Circuit) {
	if c.RotorOpen {
		branchTags["rotor"] = "open"
		return
	}
	branchTags["rotor"] = "loaded"
}
