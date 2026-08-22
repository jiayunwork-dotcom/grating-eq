package disperse

import (
	"math"

	"grating-eq/internal/grating"
)

// LinearDispersionResult carries the linear dispersion dx/dlambda on a
// focal plane located a distance f from the grating, projected on the
// normal. When the screen is placed perpendicular to the diffracted beam
// the angular dispersion converts with 1/cos(theta_m).
type LinearDispersionResult struct {
	// Order is the diffraction order.
	Order int

	// FocalMm is the detector distance in millimetres.
	FocalMm float64

	// ValueNmPerMm is the reciprocal linear dispersion in nm/mm, the
	// common instrument spec. It is f*cos(theta_m)/(m/d) when finite.
	ValueNmPerMm float64

	// Defined reports whether the value is computable.
	Defined bool
}

// LinearDispersion computes the reciprocal linear dispersion in nm/mm for a
// detector at focal distance f (mm) set perpendicular to the normal. The
// geometry converts angular dispersion dtheta/dlambda = m/(d*cos(theta_m))
// to wavelength spread per millimetre on the detector.
func LinearDispersion(g grating.Grating, res grating.OrderResult, focalMm float64) LinearDispersionResult {
	out := LinearDispersionResult{Order: res.Order, FocalMm: focalMm}
	if !res.Exists || focalMm <= 0 || res.Order == 0 {
		return out
	}
	cos := grating.CosineOf(res.AngleRad)
	if math.Abs(cos) < 1e-12 || focalMm <= 0 {
		return out
	}
	// dtheta/dlambda = m/(d*cos). Reciprocal linear dispersion in nm/mm:
	// nm per mm of detector = 1/(f * dtheta/dlambda) when the detector is
	// perpendicular to the normal.
	disp := float64(res.Order) / (g.GrooveSpacingNm * cos)
	out.ValueNmPerMm = 1.0 / (focalMm * math.Abs(disp))
	out.Defined = true
	return out
}

// LinearDispersionString renders the reciprocal linear dispersion.
func (r LinearDispersionResult) LinearDispersionString() string {
	if !r.Defined {
		return "—"
	}
	return formatValue(r.ValueNmPerMm) + " nm/mm"
}

// SpectralSpread returns the total angular spread in radians between the
// +m and -m orders for a wavelength band [lambda1, lambda2], at normal
// incidence.
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

// OrderNotFoundError reports an order that disappears inside a scan.
type OrderNotFoundError struct {
	Order  int
	Detail string
}

// Error implements the error interface.
func (e *OrderNotFoundError) Error() string {
	return "order " + itoa(e.Order) + " not found: " + e.Detail
}
