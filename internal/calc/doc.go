// Package calc interprets the phasor solution of the T-equivalent circuit
// as machine quantities: synchronous and rotor speeds, air-gap power, rotor
// copper loss, electromagnetic torque, mechanical power and the slip of
// maximum torque.
//
// # Model equations
//
// From the per-phase circuit solution (package circuit) the rotor branch
// current I2 is known. With r2 the referred rotor resistance, p the pole
// count and s the slip:
//
//	air-gap power   Pag = 3 I2^2 (r2 / s)
//	rotor copper    Pcu2 = 3 I2^2 r2
//	mechanical      Pem = Pag (1 - s) = Pag - Pcu2
//	rotor speed     wm = ws (1 - s) / (p/2),   ws = 2 pi f
//	torque          T = Pem / wm = (p/2) Pag / ws
//
// The pole count p is the number of poles, so a four-pole machine carries
// p = 4 and runs at half the electrical speed. Computing T from the air-gap
// power and the electrical synchronous speed keeps the expression finite at
// synchronous speed, where both Pem and wm vanish. The alternative Thevenin
// expression,
//
//	T = (3 p'/ws) Vth^2 (r2/s) / ((Rth + r2/s)^2 + (Xth + x2)^2),
//
// uses the pole-pair count p' = p/2 and is provided by package circuit; the
// tests cross-check it against the power method.
//
// # Operating regions
//
//   - s in (0, 1]: motoring. Power flows from the supply to the shaft,
//     torque acts in the direction of rotation.
//   - s = 0: synchronous. No slip means no induced rotor current, so the
//     torque is exactly zero and the input current is the magnetising
//     current.
//   - s < 0: generating. The rotor runs above synchronous speed, torque
//     opposes rotation and power flows back to the supply.
//   - s > 1: plugging. The rotor is driven against a reversed field and
//     dissipates the slip power.
//
// # Cross-checks
//
// The package exposes explicit checks for the rules that must hold for any
// consistent result:
//
//   - CheckZeroSlipTorque — s = 0 forces T = 0.
//   - CheckEnergyBalance — Pag = Pcu2 + Pem.
//   - CheckVoltageScaling — doubling V quadruples T and Pem.
//
// These are evaluated on every Compute call where applicable and are also
// exercised directly by the tests, together with the golden rules of the
// T-equivalent circuit: a maximum torque point exists, and it moves with
// the rotor resistance.
package calc
