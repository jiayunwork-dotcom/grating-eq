package grating

import "math"

// Convention describes which sign convention is pinned for the diffracted
// orders. The package is built around the transmission convention and does
// not switch at runtime; the type exists so reports can state the choice.
type Convention int

const (
	// Transmission pins m*lambda = d*(sin(theta_m) - sin(theta_i)).
	Transmission Convention = iota
	// Reflection is documented for completeness but not selectable here,
	// because the CLI fixes the transmission sign once.
	Reflection
)

// String renders the convention for the CLI header.
func (c Convention) String() string {
	if c == Reflection {
		return "reflection"
	}
	return "transmission"
}

// SignOf describes the geometric side of an order relative to the incident
// beam. Orders with positive m diffract to the same side of the normal as
// the incident light when the incident angle is zero.
type SignOf int

const (
	// SameSide means theta_m has the same sign as the incident angle.
	SameSide SignOf = iota
	// OppositeSide means theta_m has the opposite sign.
	OppositeSide
	// AlongNormal means theta_m is zero; the beam keeps going straight.
	AlongNormal
)

// SideOf returns which side of the normal the order leaves on. For an order
// with angle theta_m, the side is read from the sign of sin(theta_m).
func (r OrderResult) SideOf() SignOf {
	if !r.Exists {
		return AlongNormal
	}
	if math.Abs(r.SineOfAngle) < 1e-15 {
		return AlongNormal
	}
	if r.SineOfAngle > 0 {
		return SameSide
	}
	return OppositeSide
}

// ZeroOrder reports the behaviour of the m = 0 term. The grating equation
// forces sin(theta_0) = sin(theta_i), so the zero order always continues
// along the incident direction regardless of groove spacing or wavelength.
func (g Grating) ZeroOrder() OrderResult {
	res, _ := g.Eq(0)
	return res
}

// Degrees converts radians to degrees for display.
func Degrees(angleRad float64) float64 {
	return angleRad * 180 / math.Pi
}

// Radians converts degrees to radians for input.
func Radians(angleDeg float64) float64 {
	return angleDeg * math.Pi / 180
}

// AngleDistance returns the absolute angular separation between the order
// and the incident beam, in radians.
func (r OrderResult) AngleDistance(g Grating) float64 {
	if !r.Exists {
		return math.Inf(1)
	}
	return math.Abs(r.AngleRad - g.IncidentAngleRad)
}

// ReflectivitySide is a placeholder name kept out of the public API; the
// package never computes reflected intensities.
type reflectivitySide int
