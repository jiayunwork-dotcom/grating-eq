package grating

import "math"

type VisibleRegion int

const (
	RegionVisible VisibleRegion = iota

	RegionGrazing

	RegionHidden
)

func (r OrderResult) Classify() VisibleRegion {
	if !r.Exists {
		return RegionHidden
	}
	if math.Abs(r.SineOfAngle) > 1-1e-9 {
		return RegionGrazing
	}
	return RegionVisible
}

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

func (r OrderResult) ExistenceString() string {
	if r.Exists {
		return "yes"
	}
	if r.JustMissing {
		return "edge"
	}
	return "no"
}

type Overlap struct {
	LowerOrder int
	UpperOrder int

	OverlapNm float64
}

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
