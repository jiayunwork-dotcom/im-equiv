// Package calc: machine-level results.
//
// Results is the single data structure handed back by Compute. Every field
// carries its unit in the name so callers cannot confuse a torque with a
// power or a speed, and every field is populated for every operating point
// including the synchronous and generating extremes.
package calc

// Results is the complete set of quantities computed for one operating point
// and printed by the torque subcommand. Every value carries its unit in the
// field name so a client of the package cannot mistake the meaning.
type Results struct {
	// Currents.
	InputCurrentA float64 // line current of the stator, I1
	RotorCurrentA float64 // rotor branch current referred to the stator, I2
	MagnetisingA  float64 // no-load magnetising current, I0
	NodeVoltageV  float64 // voltage at the junction of the two parallel branches

	// Speeds.
	SyncSpeedRadS  float64 // electrical synchronous speed ws = 2 pi f
	RotorSpeedRadS float64 // mechanical rotor speed wm = ws (1-s) / (p/2)
	SyncSpeedRPM   float64 // synchronous speed in revolutions per minute
	RotorSpeedRPM  float64 // rotor speed in revolutions per minute

	// Powers.
	InputPowerW   float64 // electrical power drawn from the supply
	StatorCopperW float64 // stator copper loss 3 I1^2 Rs
	IronW         float64 // iron loss in the magnetising branch 3 I0^2 Rm
	AirgapPowerW  float64 // air-gap power Pag = 3 I2^2 (r2/s)
	RotorCopperW  float64 // rotor copper loss 3 I2^2 r2
	MechanicalW   float64 // electromagnetic (shaft) power Pem = Pag (1-s)
	Efficiency    float64 // Pem / Pin at the operating point
	PowerFactor   float64 // cosine of the angle between Vph and I1

	// Torque.
	TorqueNm    float64 // electromagnetic torque T = (p/2) Pag / ws
	MaxSlip     float64 // slip of maximum torque, from the Thevenin model
	MaxTorqueNm float64 // peak electromagnetic torque

	// Operating classification.
	Slip float64 // slip at which the point was computed
	Mode string  // motoring, generating, synchronous or plugging
}
