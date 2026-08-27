package disperse

import (
	"math"

	"grating-eq/internal/grating"
)

type FSRResult struct {
	Order int

	ValueNm float64

	Defined bool
}

func FreeSpectralRange(wavelengthNm float64, m int) FSRResult {
	out := FSRResult{Order: m}
	if m == 0 {
		return out
	}
	out.ValueNm = wavelengthNm / math.Abs(float64(m))
	out.Defined = true
	return out
}

func (r FSRResult) FSRString() string {
	if !r.Defined {
		return "—"
	}
	return formatValue(r.ValueNm) + " nm"
}

func OverlapFSR(g grating.Grating, m int) (float64, error) {
	ov, err := g.AdjacentOverlap(m)
	if err != nil {
		return 0, err
	}
	return ov.OverlapNm, nil
}
