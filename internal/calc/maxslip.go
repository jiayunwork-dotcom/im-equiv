// Package calc: numerical location of the maximum-torque slip.
//
// Two independent routes locate the peak of the torque-slip characteristic.
// The closed-form route follows from the Thevenin equivalent (package
// circuit) and is evaluated in Compute; the numerical route scans the full
// circuit over the motoring band. Agreement between the two is asserted by
// the tests and is a strong check that neither the circuit model nor the
// torque interpretation drifts.
package calc

import (
	"math"

	"im-equiv/internal/config"
)

// FindMaxTorqueSlip locates the slip of maximum electromagnetic torque in the
// interval [lo, hi] by golden-section search over the torque-slip curve
// computed from the full circuit. The result cross-checks the closed-form
// Thevenin expression MaxSlip: both must agree for every physical input, and
// the tests assert exactly that.
//
// The interval must lie in the motoring region (lo > 0). The torque curve is
// single-peaked there, so golden-section search is exact to within the
// stopping tolerance.
func FindMaxTorqueSlip(in config.Input, lo, hi float64) float64 {
	const phi = 0.6180339887498949

	a, b := lo, hi
	c := b - phi*(b-a)
	d := a + phi*(b-a)
	fc := TorqueAtSlip(in, c)
	fd := TorqueAtSlip(in, d)

	for i := 0; i < 200 && math.Abs(b-a) > 1e-12; i++ {
		if fc > fd {
			b, d, fd = d, c, fc
			c = b - phi*(b-a)
			fc = TorqueAtSlip(in, c)
		} else {
			a, c, fc = c, d, fd
			d = a + phi*(b-a)
			fd = TorqueAtSlip(in, d)
		}
	}
	return (a + b) / 2
}

// TorqueSweep samples the torque-slip characteristic at n evenly spaced
// slips in the interval [lo, hi] and returns the (slip, torque) pairs. It is
// a convenience for plotting or for spot-checking the characteristic at a
// handful of operating points.
func TorqueSweep(in config.Input, lo, hi float64, n int) [][2]float64 {
	if n < 2 {
		n = 2
	}
	points := make([][2]float64, 0, n)
	for i := 0; i < n; i++ {
		s := lo + (hi-lo)*float64(i)/float64(n-1)
		points = append(points, [2]float64{s, TorqueAtSlip(in, s)})
	}
	return points
}
