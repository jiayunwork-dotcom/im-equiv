// Package calc: characteristic curve sweeps.
package calc

import (
	"im-equiv/internal/config"
)

// CurvePoint is a single sample of a torque-speed characteristic: the slip
// and the quantity evaluated at that slip.
type CurvePoint struct {
	Slip       float64
	TorqueNm   float64
	PowerW     float64
	Efficiency float64
}

// TorqueCurve samples the torque-slip characteristic at n evenly spaced
// slips in [lo, hi] in the motoring region. The result contains the torque,
// the mechanical power and the efficiency at every sample, which is enough
// to plot the full electromechanical behaviour of the machine.
func TorqueCurve(in config.Input, lo, hi float64, n int) []CurvePoint {
	if n < 2 {
		n = 2
	}
	points := make([]CurvePoint, 0, n)
	for i := 0; i < n; i++ {
		s := lo + (hi-lo)*float64(i)/float64(n-1)
		points = append(points, sample(in, s))
	}
	return points
}

// sample evaluates one slip of the curve, reusing the full compute pipeline
// so torque, mechanical power and efficiency are mutually consistent.
func sample(in config.Input, slip float64) CurvePoint {
	cur := in
	cur.S = slip
	res, err := Compute(cur)
	if err != nil {
		return CurvePoint{Slip: slip}
	}
	return CurvePoint{
		Slip:       slip,
		TorqueNm:   res.TorqueNm,
		PowerW:     res.MechanicalW,
		Efficiency: res.Efficiency,
	}
}

// MaxTorquePoint returns the sample of TorqueCurve that is closest to the
// analytic maximum-torque slip, giving a curve-based view of the peak.
func MaxTorquePoint(in config.Input, lo, hi float64, n int) CurvePoint {
	points := TorqueCurve(in, lo, hi, n)
	if len(points) == 0 {
		return CurvePoint{}
	}
	best := points[0]
	for _, p := range points[1:] {
		if p.TorqueNm > best.TorqueNm {
			best = p
		}
	}
	return best
}

// EfficiencyCurve isolates the efficiency samples of the curve for slips in
// the motoring region. Efficiency is only defined while power flows from the
// supply to the shaft, so points outside (0, 1] are dropped.
func EfficiencyCurve(in config.Input, lo, hi float64, n int) []CurvePoint {
	points := TorqueCurve(in, lo, hi, n)
	out := points[:0]
	for _, p := range points {
		if p.Slip > 0 && p.Slip <= 1 {
			out = append(out, p)
		}
	}
	return out
}
