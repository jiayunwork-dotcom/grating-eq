package grating

import "math"

// VisibleRegion classifies an order angle against the physical constraint
// that every diffracted direction must lie within 90 degrees of the normal.
type VisibleRegion int

const (
	// RegionVisible means |sin(theta_m)| <= 1 and the order exists.
	RegionVisible VisibleRegion = iota

	// RegionGrazing means the order exists but sits at nearly +/- 90
	// degrees; it is at the very edge of the visible region.
	RegionGrazing

	// RegionHidden means |sin(theta_m)| > 1 and the order cannot exist.
	RegionHidden
)

// Classify returns the visible-region class of an evaluated order.
func (r OrderResult) Classify() VisibleRegion {
	if !r.Exists {
		return RegionHidden
	}
	if math.Abs(r.SineOfAngle) > 1-1e-9 {
		return RegionGrazing
	}
	return RegionVisible
}

// RegionName renders the classification for the CLI table.
func (r OrderResult) RegionName() string {
	switch r.Classify() {
	case RegionGrazing:
		return "grazing"
	case RegionVisible:
		return "visible"
	default:
		return "hidden"
	}
}

// ExistenceString renders the existence column: "yes" when the order
// exists, "no" otherwise. Borderline round-off cases render as "edge".
func (r OrderResult) ExistenceString() string {
	if r.Exists {
		return "yes"
	}
	if r.JustMissing {
		return "edge"
	}
	return "no"
}

// Overlap describes two adjacent orders whose wavelength windows overlap,
// which defines the free spectral range of the grating.
type Overlap struct {
	// LowerOrder and UpperOrder are the pair (m, m+1).
	LowerOrder int
	UpperOrder int

	// OverlapNm is the wavelength span where both orders diffract.
	OverlapNm float64
}

// AdjacentOverlap computes the overlap of order m and order m+1. At normal
// incidence the two windows begin to overlap when lambda/d reaches 1/|m|,
// which is the classic FSR criterion.
func (g Grating) AdjacentOverlap(m int) (Overlap, error) {
	a, err := g.WavelengthRange(m)
	if err != nil {
		return Overlap{}, err
	}
	b, err := g.WavelengthRange(m + 1)
	if err != nil {
		return Overlap{}, err
	}
	lo := math.Max(a.Lo, b.Lo)
	hi := math.Min(a.Hi, b.Hi)
	span := hi - lo
	if span < 0 || a.Empty() || b.Empty() {
		span = 0
	}
	return Overlap{LowerOrder: m, UpperOrder: m + 1, OverlapNm: span}, nil
}
