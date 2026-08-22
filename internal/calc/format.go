// Package calc: rendering of results.
//
// The same Results value can be rendered three ways:
//
//   - FormatResults, the aligned table meant for a terminal;
//   - FormatCSV, one row per quantity, meant for spreadsheets and diffs;
//   - FormatJSON, the structured form, meant for other programs.
//
// All three share the number rendering helper so a value never appears as
// scientific notation or as a raw IEEE special value.
package calc

import (
	"fmt"
	"math"
	"strings"
)

// FormatResults renders the computed quantities as an aligned, human-readable
// table. Line labels are fixed width so the values line up in a terminal,
// and every quantity carries its unit. Non-finite values (e.g. an infinite
// maximum-torque slip) are shown as +Inf rather than as the raw IEEE
// spelling.
func FormatResults(res Results) string {
	var b strings.Builder
	write := func(label, unit string, v float64) {
		fmt.Fprintf(&b, "%-22s %10s  %s\n", label, number(v), unit)
	}
	writeText := func(label, text, unit string) {
		fmt.Fprintf(&b, "%-22s %10s  %s\n", label, text, unit)
	}

	write("input current I1", "A", res.InputCurrentA)
	write("rotor current I2", "A", res.RotorCurrentA)
	write("magnetising I0", "A", res.MagnetisingA)
	write("node voltage Vn", "V", res.NodeVoltageV)
	write("sync speed ws", "rad/s", res.SyncSpeedRadS)
	write("rotor speed wm", "rad/s", res.RotorSpeedRadS)
	write("sync speed", "rpm", res.SyncSpeedRPM)
	write("rotor speed", "rpm", res.RotorSpeedRPM)
	write("input power Pin", "W", res.InputPowerW)
	write("stator copper Pcu1", "W", res.StatorCopperW)
	write("iron loss Pfe", "W", res.IronW)
	write("airgap power Pag", "W", res.AirgapPowerW)
	write("rotor copper Pcu2", "W", res.RotorCopperW)
	write("mechanical Pem", "W", res.MechanicalW)
	if res.Mode == "motoring" {
		write("efficiency", "%", res.Efficiency*100)
	} else {
		// Efficiency is only meaningful in the motoring mode, where power
		// flows from the supply to the shaft.
		writeText("efficiency", "-", "%")
	}
	write("power factor", "", res.PowerFactor)
	write("torque T", "Nm", res.TorqueNm)
	write("max-torque slip sm", "", res.MaxSlip)
	write("max torque Tmax", "Nm", res.MaxTorqueNm)
	fmt.Fprintf(&b, "%-22s %10s  %s\n", "mode", "", res.Mode)
	return b.String()
}

// number renders a float with four decimals, mapping infinities to +Inf and
// -Inf and NaN to NaN so the output never contains scientific notation.
func number(v float64) string {
	switch {
	case math.IsInf(v, 1):
		return "+Inf"
	case math.IsInf(v, -1):
		return "-Inf"
	case math.IsNaN(v):
		return "NaN"
	default:
		return fmt.Sprintf("%.4f", v)
	}
}
