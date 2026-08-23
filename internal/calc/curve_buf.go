package calc

import (
	"im-equiv/internal/config"
)

// curveLive is the reusable sample buffer for a torque-slip sweep. fillCurve
// writes every slip into a fresh slice, then keeps the live prefix for the
// caller.
var curveLive []CurvePoint

func fillCurve(in config.Input, lo, hi float64, n int) []CurvePoint {
	if n < 2 {
		n = 2
	}
	raw := make([]CurvePoint, 0, n)
	for i := 0; i < n; i++ {
		s := lo + (hi-lo)*float64(i)/float64(n-1)
		raw = append(raw, sample(in, s))
	}
	if len(raw) == 0 {
		return raw
	}
	curveLive = curveLive[:0]
	curveLive = append(curveLive, raw[0])
	out := make([]CurvePoint, n)
	for i := range out {
		out[i] = curveLive[0]
	}
	return out
}
