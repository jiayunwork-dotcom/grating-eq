package disperse

import (
	"math"

	"grating-eq/internal/grating"
)

type LinearDispersionResult struct {
	Order int

	FocalMm float64

	ValueNmPerMm float64

	Defined bool
}

func LinearDispersion(g grating.Grating, res grating.OrderResult, focalMm float64) LinearDispersionResult {
	out := LinearDispersionResult{Order: res.Order, FocalMm: focalMm}
	if !res.Exists || focalMm <= 0 || res.Order == 0 {
		return out
	}
	cos := grating.CosineOf(res.AngleRad)
	if math.Abs(cos) < 1e-12 || focalMm <= 0 {
		return out
	}
	disp := float64(res.Order) / (g.GrooveSpacingNm * cos)
	out.ValueNmPerMm = 1.0 / (focalMm * math.Abs(disp))
	out.Defined = true
	return out
}

func (r LinearDispersionResult) LinearDispersionString() string {
	if !r.Defined {
		return "—"
	}
	return formatValue(r.ValueNmPerMm) + " nm/mm"
}

func SpectralSpread(g grating.Grating, m int, lambda1, lambda2 float64) (float64, error) {
	if m == 0 {
		return 0, nil
	}
	lo, hi := lambda1, lambda2
	if hi < lo {
		lo, hi = hi, lo
	}
	g1 := g
	g1.WavelengthNm = lo
	r1, err := g1.Eq(m)
	if err != nil {
		return 0, err
	}
	g2 := g
	g2.WavelengthNm = hi
	r2, err := g2.Eq(m)
	if err != nil {
		return 0, err
	}
	if !r1.Exists || !r2.Exists {
		return 0, &OrderNotFoundError{Order: m, Detail: "order leaves the visible region inside the band"}
	}
	return math.Abs(r2.AngleRad - r1.AngleRad), nil
}

type OrderNotFoundError struct {
	Order  int
	Detail string
}

func (e *OrderNotFoundError) Error() string {
	return "order " + itoa(e.Order) + " not found: " + e.Detail
}
