// Package circuit solves the per-phase phasor equations of the T-equivalent
// circuit of a three-phase induction machine.
//
// # Circuit topology
//
// The machine is modelled by the classical T-equivalent circuit, one phase
// of which is shown below (all quantities are per-phase, referred to the
// stator):
//
//	       I1         Zs = Rs + j Xs        node      I0
//	Vph ────┬─────/\/\/\───┬───////////─────┬──────────────┬───
//	       │              │                │  Zm = Rm + j Xm │
//	       │              │                │        (magnetising)
//	       │              │                │              │
//	       │              │                │              │
//	       └──────────────┴────────────────┴──────────────┴───
//	                                        │
//	                               Zr = r2/s + j x2
//	                                        │   (rotor, referred)
//	                                        └───────────────
//
// The stator branch Zs carries the full line current I1. Beyond the node
// the current splits between the magnetising branch Zm and the rotor branch
// Zr; the rotor current I2 flows through r2/s + j x2. The slip-dependent
// resistance r2/s is what converts mechanical power: the power dissipated
// in r2/s minus the loss in r2 is delivered at the shaft.
//
// # Why the slip sits in the rotor resistance
//
// A rotating field sweeping past the rotor at slip frequency omega_s * s
// induces rotor currents at that frequency. Viewed from the stator the
// rotor resistance seen by the supply is r2/s: at standstill (s = 1) it is
// r2, while near synchronous speed the branch approaches an open circuit.
// The same factor converts the air-gap power into mechanical power,
//
//	Pag = 3 I2^2 r2 / s,        Pem = Pag (1 - s),
//
// which is why j omega and s must both live in the rotor branch and why
// "r2/s" may never be replaced by "r2 * s".
//
// # Conventions
//
//   - All voltages and currents are RMS phasors, complex128 in this package.
//   - A balanced, wye-connected stator is assumed. The phase voltage is the
//     line-to-line voltage divided by sqrt(3) (see config.Input).
//   - The electrical synchronous speed is ws = 2 pi f rad/s. With p poles
//     the mechanical angular velocity is wm = ws (1 - s) / (p/2) rad/s, so a
//     four-pole machine (p = 4) runs at exactly half the electrical speed.
//   - Slip s in (0, 1] is motoring, negative slip is generating and s = 0 is
//     synchronous operation. At s = 0 the rotor branch carries no current
//     and the electromagnetic torque is exactly zero.
//
// # Solving the circuit
//
// Solve performs the three classic steps: combine the parallel branches,
// divide the phase voltage by the total impedance to get the input current,
// then split the node voltage across the magnetising and rotor branches. The
// synchronous point is handled by an explicit open-rotor path instead of a
// division by zero.
//
// The package deliberately keeps no knowledge of the machine-level
// interpretation (torque, power split, efficiency). That layer lives in
// package calc, which consumes the phasor solution produced here.
package circuit
