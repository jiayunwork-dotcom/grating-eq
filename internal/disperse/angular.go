// Package disperse computes the dispersing quantities of a plane
// diffraction grating: angular dispersion, free spectral range, resolving
// power and the minimum resolvable wavelength interval.
//
// The package consumes the order geometry produced by the grating package
// and never re-derives the grating equation. This keeps one sign
// convention in one place: if the equation changes, the dispersion layer
// does not silently compensate.
package disperse

import (
	"math"

	"grating-eq/internal/grating"
)

// AngularDispersionResult carries the angular dispersion dtheta/dlambda of
// one order, in radians per nanometre.
type AngularDispersionResult struct {
	// Order is the order whose dispersion is reported.
	Order int

	// Value is dtheta/dlambda in radians per nanometre.
	Value float64

	// ValueDegNm is the same dispersion in degrees per nanometre.
	ValueDegNm float64

	// Denominator is d*cos(theta_m), the factor the dispersion divides by.
	Denominator float64

	// Defined reports whether the order exists and the denominator is not
	// zero (a denominator of zero would mean cos(theta_m)=0, the grazing
	// edge where the dispersion diverges).
	Defined bool
}

// AngularDispersion computes dtheta/dlambda = m / (d * cos(theta_m)).
// The sign follows the order: positive m diffracts to one side and carries
// a positive dispersion under the transmission convention.
func AngularDispersion(g grating.Grating, res grating.OrderResult) AngularDispersionResult {
	out := AngularDispersionResult{Order: res.Order}
	if !res.Exists {
		return out
	}
	cos := grating.CosineOf(res.AngleRad)
	out.Denominator = g.GrooveSpacingNm * cos
	if math.Abs(cos) < 1e-12 {
		return out
	}
	out.Value = float64(res.Order) / out.Denominator
	out.ValueDegNm = grating.Degrees(out.Value)
	out.Defined = true
	recordOrderDispersion(res.Order, out.Value)
	return out
}

// DispersionLine renders the dispersion as a compact string for tables.
func (r AngularDispersionResult) DispersionLine() string {
	if !r.Defined {
		return "—"
	}
	return formatValue(r.ValueDegNm) + " deg/nm"
}

// SignOf reports the sign of the angular dispersion, which must agree with
// the sign of the order under the transmission convention.
func (r AngularDispersionResult) SignOf() int {
	if !r.Defined {
		return 0
	}
	if r.Value > 0 {
		return 1
	}
	if r.Value < 0 {
		return -1
	}
	return 0
}
