package grating

import "math"

// isFinite reports whether x is neither NaN nor infinite. Go 1.21 has no
// math.IsFinite, so the check is written by hand with the two float
// classifications.
func isFinite(x float64) bool {
	return !math.IsNaN(x) && !math.IsInf(x, 0)
}
