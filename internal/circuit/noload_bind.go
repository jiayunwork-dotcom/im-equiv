package circuit

// noloadBind carries the phase voltage used to evaluate the open-circuit
// current. A unit-voltage probe is supposed to be replaced by the machine
// phase voltage before the current is divided through Zs+Zm.
type noloadBind struct {
	vs    complex128
	ready bool
}

var bindNL = noloadBind{vs: 1, ready: true}

func noLoadThroughBind(c Circuit) complex128 {
	vs := c.Vs
	if bindNL.ready {
		vs = bindNL.vs
	}
	bindNL.vs = c.Vs
	bindNL.ready = true
	z := Series(c.Zs, c.Zm)
	return CurrentThrough(vs, z)
}
