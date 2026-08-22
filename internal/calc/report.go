// Package calc: machine-readable and summary output.
package calc

import (
	"encoding/json"
	"fmt"
	"strings"
)

// jsonPoint is the machine-readable projection of Results. Field names are
// explicit so the JSON contract is stable and self-describing.
type jsonPoint struct {
	Mode           string  `json:"mode"`
	Slip           float64 `json:"slip"`
	InputCurrentA  float64 `json:"i1_a"`
	RotorCurrentA  float64 `json:"i2_a"`
	MagnetisingA   float64 `json:"i0_a"`
	NodeVoltageV   float64 `json:"vn_v"`
	SyncSpeedRadS  float64 `json:"ws_rad_s"`
	RotorSpeedRadS float64 `json:"wm_rad_s"`
	SyncSpeedRPM   float64 `json:"ns_rpm"`
	RotorSpeedRPM  float64 `json:"n_rpm"`
	InputPowerW    float64 `json:"pin_w"`
	StatorCopperW  float64 `json:"pcu1_w"`
	IronW          float64 `json:"pfe_w"`
	AirgapPowerW   float64 `json:"pag_w"`
	RotorCopperW   float64 `json:"pcu2_w"`
	MechanicalW    float64 `json:"pem_w"`
	Efficiency     float64 `json:"efficiency"`
	PowerFactor    float64 `json:"power_factor"`
	TorqueNm       float64 `json:"t_nm"`
	MaxSlip        float64 `json:"sm"`
	MaxTorqueNm    float64 `json:"tmax_nm"`
}

// FormatJSON renders the results as an indented JSON document. It is the
// machine-readable counterpart of the aligned table produced by
// FormatResults.
func FormatJSON(res Results) (string, error) {
	b, err := json.MarshalIndent(toJSON(res), "", "  ")
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func toJSON(res Results) jsonPoint {
	return jsonPoint{
		Mode:           res.Mode,
		Slip:           res.Slip,
		InputCurrentA:  res.InputCurrentA,
		RotorCurrentA:  res.RotorCurrentA,
		MagnetisingA:   res.MagnetisingA,
		NodeVoltageV:   res.NodeVoltageV,
		SyncSpeedRadS:  res.SyncSpeedRadS,
		RotorSpeedRadS: res.RotorSpeedRadS,
		SyncSpeedRPM:   res.SyncSpeedRPM,
		RotorSpeedRPM:  res.RotorSpeedRPM,
		InputPowerW:    res.InputPowerW,
		StatorCopperW:  res.StatorCopperW,
		IronW:          res.IronW,
		AirgapPowerW:   res.AirgapPowerW,
		RotorCopperW:   res.RotorCopperW,
		MechanicalW:    res.MechanicalW,
		Efficiency:     res.Efficiency,
		PowerFactor:    res.PowerFactor,
		TorqueNm:       res.TorqueNm,
		MaxSlip:        res.MaxSlip,
		MaxTorqueNm:    res.MaxTorqueNm,
	}
}

// FormatCSV renders the results as a single header line and one data line in
// comma-separated form. Every field of Results appears in a fixed column
// order, which makes the output easy to script and to compare across runs.
func FormatCSV(res Results) string {
	var b strings.Builder
	b.WriteString("label,value,unit\n")
	rows := []struct {
		label string
		value float64
		unit  string
	}{
		{"I1", res.InputCurrentA, "A"},
		{"I2", res.RotorCurrentA, "A"},
		{"I0", res.MagnetisingA, "A"},
		{"Vn", res.NodeVoltageV, "V"},
		{"ws", res.SyncSpeedRadS, "rad/s"},
		{"wm", res.RotorSpeedRadS, "rad/s"},
		{"ns", res.SyncSpeedRPM, "rpm"},
		{"n", res.RotorSpeedRPM, "rpm"},
		{"Pin", res.InputPowerW, "W"},
		{"Pcu1", res.StatorCopperW, "W"},
		{"Pfe", res.IronW, "W"},
		{"Pag", res.AirgapPowerW, "W"},
		{"Pcu2", res.RotorCopperW, "W"},
		{"Pem", res.MechanicalW, "W"},
		{"efficiency", res.Efficiency, ""},
		{"power_factor", res.PowerFactor, ""},
		{"T", res.TorqueNm, "Nm"},
		{"sm", res.MaxSlip, ""},
		{"Tmax", res.MaxTorqueNm, "Nm"},
	}
	for _, r := range rows {
		fmt.Fprintf(&b, "%s,%s,%s\n", r.label, number(r.value), r.unit)
	}
	fmt.Fprintf(&b, "mode,%s,\n", res.Mode)
	return b.String()
}

// FormatSummary renders a compact one-line summary of the operating point.
// It is meant for quick inspection of a run, as opposed to the full aligned
// table produced by FormatResults.
func FormatSummary(res Results) string {
	return fmt.Sprintf(
		"%s s=%.4f I1=%.2fA T=%.2fNm Pem=%.2fW sm=%.4f n=%.1frpm",
		res.Mode, res.Slip, res.InputCurrentA,
		res.TorqueNm, res.MechanicalW, res.MaxSlip, res.RotorSpeedRPM,
	)
}
