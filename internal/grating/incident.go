package grating

import "math"

type Convention int

const (
	Transmission Convention = iota
	Reflection
)

func (c Convention) String() string {
	if c == Reflection {
		return "reflection"
	}
	return "transmission"
}

type SignOf int

const (
	SameSide SignOf = iota
	OppositeSide
	AlongNormal
)

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

func (g Grating) ZeroOrder() OrderResult {
	res, _ := g.Eq(0)
	return res
}

func Degrees(angleRad float64) float64 {
	return angleRad * 180 / math.Pi
}

func Radians(angleDeg float64) float64 {
	return angleDeg * math.Pi / 180
}

func (r OrderResult) AngleDistance(g Grating) float64 {
	if !r.Exists {
		return math.Inf(1)
	}
	return math.Abs(r.AngleRad - g.IncidentAngleRad)
}

type reflectivitySide int
