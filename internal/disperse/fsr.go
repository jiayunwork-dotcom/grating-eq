package disperse

import (
	"math"

	"grating-eq/internal/grating"
)

// FSRResult carries the free spectral range of one order.
type FSRResult struct {
	// Order is the order whose free spectral range is reported.
	Order int

	// ValueNm is the FSR in nanometres, defined as lambda/|m|.
	ValueNm float64

	// Defined reports whether |m| >= 1 so the FSR is meaningful.
	Defined bool
}

// FreeSpectralRange returns lambda/|m|, the wavelength span over which
// orders m and m+1 do not overlap. The definition is pinned to adjacent
// order overlap, not to any device-specific finesse factor.
func FreeSpectralRange(wavelengthNm float64, m int) FSRResult {
	out := FSRResult{Order: m}
	if m == 0 {
		return out
	}
	out.ValueNm = wavelengthNm / math.Abs(float64(m))
	out.Defined = true
	return out
}

// FSRString renders the FSR for tables.
func (r FSRResult) FSRString() string {
	if !r.Defined {
		return "—"
	}
	return formatValue(r.ValueNm) + " nm"
}

// OverlapFSR derives the free spectral range from the overlap of adjacent
// order wavelength windows. It is a cross-check for lambda/|m|: at normal
// incidence the two definitions agree.
func OverlapFSR(g grating.Grating, m int) (float64, error) {
	ov, err := g.AdjacentOverlap(m)
	if err != nil {
		return 0, err
	}
	return ov.OverlapNm, nil
}
