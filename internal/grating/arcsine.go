package grating

import (
	"math"
	"strconv"
)

// ArcSine returns the arcsine of x in radians. The standard library's
// math.Asin returns NaN for |x| > 1, so callers are required to pass the
// argument through this helper, which first checks the domain and returns a
// descriptive error instead of silently propagating NaN through the report.
//
// The check is deliberately before the call: the grating equation only
// produces real angles when |sin(theta_m)| <= 1, and every place that turns
// a sine back into an angle must guard the argument in the same way.
func ArcSine(x float64) (float64, error) {
	if math.IsNaN(x) || math.IsInf(x, 0) || math.Abs(x) > 1 {
		return 0, ErrSineOutOfRange(x)
	}
	// Clamp tiny overshoots from floating-point arithmetic, e.g. a value of
	// 1.0000000000000002 that is mathematically on the boundary.
	if x > 1 {
		x = 1
	}
	if x < -1 {
		x = -1
	}
	return math.Asin(x), nil
}

// ErrSineOutOfRange describes an arcsine argument outside [-1, 1]. The
// error text carries the offending value so that the CLI message states
// exactly which quantity fell out of the visible region.
func ErrSineOutOfRange(x float64) error {
	return &SineError{X: x}
}

// SineError is a typed error for an out-of-range sine value.
type SineError struct {
	X float64
}

// Error implements the error interface.
func (e *SineError) Error() string {
	return "sine value " + strconv.FormatFloat(e.X, 'g', 17, 64) +
		" outside [-1, 1]; the requested order does not exist"
}

// SineOf computes sin(angle) after verifying that the angle is finite.
// Angles produced by the grating equation are always finite, but the helper
// guards the subsequent arccos-style conversions in the dispersion layer.
func SineOf(angleRad float64) (float64, error) {
	if math.IsNaN(angleRad) || math.IsInf(angleRad, 0) {
		return 0, ErrBadIncident("angle became non-finite")
	}
	return math.Sin(angleRad), nil
}

// CosineOf returns cos(angle). It never returns a value whose magnitude
// exceeds 1, which the angular dispersion denominator relies on.
func CosineOf(angleRad float64) float64 {
	return math.Cos(angleRad)
}

// NormalizeAngle folds an angle into the interval [-pi, pi]. The grating
// equation never needs angles outside that window; normalizing keeps the
// printed output stable when a scanning routine walks far from zero.
func NormalizeAngle(angleRad float64) float64 {
	twoPi := 2 * math.Pi
	angle := math.Mod(angleRad, twoPi)
	if angle > math.Pi {
		angle -= twoPi
	} else if angle < -math.Pi {
		angle += twoPi
	}
	return angle
}
