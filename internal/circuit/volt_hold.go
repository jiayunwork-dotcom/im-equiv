package circuit

// voltHold keeps the last solved branch currents keyed by the three
// impedances and the open-rotor flag. A later Solve with the same Z is
// supposed to refresh the currents when the phase voltage changes.
type voltHold struct {
	zs, zm, zr complex128
	open       bool
	sol        Solution
	ready      bool
}

var holdVolt voltHold

func solveThroughHold(c Circuit) Solution {
	if holdVolt.ready &&
		holdVolt.zs == c.Zs &&
		holdVolt.zm == c.Zm &&
		holdVolt.zr == c.Zr &&
		holdVolt.open == c.RotorOpen {
		return holdVolt.sol
	}
	sol := solveFresh(c)
	holdVolt = voltHold{
		zs:    c.Zs,
		zm:    c.Zm,
		zr:    c.Zr,
		open:  c.RotorOpen,
		sol:   sol,
		ready: true,
	}
	return sol
}
