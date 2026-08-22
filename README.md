# im-equiv

A command-line calculator for the **T-equivalent circuit** of a three-phase
induction machine. Given the per-phase circuit parameters — stator
resistance and leakage reactance, magnetising branch, referred rotor
resistance and leakage reactance — and the slip, it computes the input
current, the electromagnetic torque, the mechanical power and the slip of
maximum torque.

This is an equivalent-circuit *kernel*: it models the electromechanical
power conversion through the classical T circuit, nothing more. There is no
motor catalogue, no drive, no service-ticket logic.

## What it computes

For one operating point the tool prints an aligned table with:

- currents: line current `I1`, rotor current `I2`, magnetising current `I0`;
- speeds: synchronous speed `ws`, rotor speed `wm` (both rad/s and rpm);
- power split: input power, stator copper loss, iron loss, air-gap power
  `Pag`, rotor copper loss `Pcu2`, mechanical power `Pem`, efficiency;
- torque: electromagnetic torque `T`, maximum-torque slip `sm`, peak torque;
- operating mode: motoring / generating / synchronous / plugging.

The model equations are the standard ones:

```
Pag = 3 I2^2 (r2 / s)          air-gap power
Pem = Pag (1 - s)              mechanical power, Pag - Pcu2
wm  = ws (1 - s) / (p/2)       mechanical speed, ws = 2 pi f
T   = (p/2) Pag / ws           electromagnetic torque
```

The torque is also cross-checked against the Thevenin closed form in the
test suite.

## Input file

The input is a single JSON object. The shipped example
`example/4pole-50hz.json` describes a four-pole, 50 Hz, 380 V machine at 3%
slip:

```json
{
  "p":  4,
  "f":  50,
  "V":  380,
  "Rs": 1.5,
  "Xs": 2.0,
  "Rm": 30.0,
  "Xm": 20.0,
  "r2": 0.8,
  "x2": 1.5,
  "s":  0.03
}
```

All quantities are per-phase and RMS; the stator is assumed wye connected,
so the phase voltage is `V / sqrt(3)`. `p` is the number of poles (a
four-pole machine has `p = 4`).

## Usage

```bash
go run . torque example/4pole-50hz.json
```

A second example runs the same machine as a generator:

```bash
go run . torque example/4pole-50hz-generating.json
```

Errors — unreadable files, malformed JSON, or physically invalid inputs such
as a non-positive frequency, non-positive pole count, or negative circuit
parameters — are printed on standard error and the process exits non-zero
(1 for run errors, 2 for usage errors).

## Key conventions

- Slip `s` in `(0, 1]` is motoring, `s < 0` is generating, `s = 0` is
  synchronous operation at exactly zero torque.
- `j omega` and the slip both live in the rotor branch as `r2/s + j x2`;
  the tool never conflates `r2/s` with `r2 * s`.
- The rotor speed uses `wm = ws (1 - s) / (p/2)`; dropping the `(1 - s)`
  factor is the classic error this tool guards against.
- A consistent result always satisfies `Pag = Pcu2 + Pem`, zero slip gives
  zero torque, and doubling the voltage quadruples torque and mechanical
  power. These invariants are asserted on every computation.

## Build & test

```bash
go build ./...
go test ./...
```

The test suite covers the failure modes above plus the cross-check rules
(voltage scaling, rotor-resistance scaling of the maximum-torque slip, the
energy balance at the rated point) and the equivalence of the power method
with the Thevenin formula.
