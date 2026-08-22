// Package grating implements the plane diffraction grating core: the
// groove geometry, the grating equation and the order-existence test.
//
// Angles are measured from the grating normal. The transmission convention
// is pinned: the diffracted order m satisfies
//
//	m*lambda = d*(sin(theta_m) - sin(theta_i))
//
// where d is the groove spacing, theta_i the incident angle and theta_m the
// angle of order m. A positive theta points to the same side of the normal
// as the incident beam; the zero order always continues along the incident
// direction.
package grating

import "math"

// Physical constants used for unit conversions. The tool accepts groove
// geometry either as a groove density (grooves per millimetre) or directly
// as a groove spacing d in nanometres.
const (
	// NanoPerMilli is 1e6 nanometres in one millimetre.
	NanoPerMilli = 1e6

	// MaxWavelengthNm is an arbitrary upper bound used to reject
	// nonsensical wavelengths before they overflow intermediate maths.
	MaxWavelengthNm = 1e12

	// MinGrooveSpacingNm is the smallest groove spacing accepted. Below a
	// fraction of a nanometre no optical diffraction grating is physically
	// meaningful.
	MinGrooveSpacingNm = 0.1

	// DefaultMaxOrderScan bounds how far the orders command scans away from
	// zero. Orders beyond the scan window are reported as not scanned.
	DefaultMaxOrderScan = 64

	// Tolerance used when comparing angles that must be equal up to float
	// round-off, for example the zero order against the incident angle.
	AngleTolerance = 1e-9
)

// Grating holds the groove geometry and the illumination of a plane
// transmission grating. All fields are in SI units except WavelengthNm and
// GrooveSpacingNm which use nanometres for readability.
type Grating struct {
	// GrooveSpacingNm is the centre-to-centre distance between adjacent
	// grooves, in nanometres. It is the reciprocal of the groove density.
	GrooveSpacingNm float64

	// WavelengthNm is the illumination wavelength in nanometres.
	WavelengthNm float64

	// IncidentAngleRad is the incident angle measured from the grating
	// normal, positive on the same side as the positive diffracted orders.
	IncidentAngleRad float64

	// Slits is the finite number of illuminated grooves. Zero means an
	// unbounded grating where only angular dispersion is reported.
	Slits int
}

// New returns a Grating whose groove spacing and wavelength are given in
// nanometres and whose incident angle is given in radians. The constructor
// validates the parameters and reports a descriptive error otherwise.
func New(grooveSpacingNm, wavelengthNm, incidentAngleRad float64, slits int) (Grating, error) {
	g := Grating{
		GrooveSpacingNm:  grooveSpacingNm,
		WavelengthNm:     wavelengthNm,
		IncidentAngleRad: incidentAngleRad,
		Slits:            slits,
	}
	if err := g.Validate(); err != nil {
		return Grating{}, err
	}
	return g, nil
}

// Validate checks every invariant that the rest of the package relies on.
// It returns a non-nil error with a human-readable message for the CLI when
// any input falls outside the physical domain.
func (g Grating) Validate() error {
	if !isFinite(g.GrooveSpacingNm) || g.GrooveSpacingNm <= 0 {
		return ErrBadSpacing(g.GrooveSpacingNm)
	}
	if !isFinite(g.WavelengthNm) || g.WavelengthNm <= 0 {
		return ErrBadWavelength(g.WavelengthNm)
	}
	if !isFinite(g.IncidentAngleRad) {
		return ErrBadIncident("incident angle must be finite")
	}
	if s := math.Sin(g.IncidentAngleRad); math.IsNaN(s) || math.Abs(s) > 1 {
		return ErrBadIncident("|sin(incident)| must not exceed 1")
	}
	if g.Slits < 0 {
		return ErrBadSlits(g.Slits)
	}
	return nil
}

// GrooveDensity returns the number of grooves per millimetre. It is the
// exact reciprocal of the groove spacing expressed in millimetres. Callers
// that accept a density from JSON must obtain the spacing through this
// relation so that d and the line density always stay consistent.
func (g Grating) GrooveDensity() float64 {
	return NanoPerMilli / g.GrooveSpacingNm
}

// SpacingFromDensity converts a groove density in grooves per millimetre to
// a groove spacing in nanometres. The conversion is the inverse of
// GrooveDensity and is used when a JSON file names the density instead of
// the spacing.
func SpacingFromDensity(groovesPerMm float64) float64 {
	return NanoPerMilli / groovesPerMm
}

// LineDensityName returns the display form of the groove density, for
// example "600 grooves/mm" for the bundled example.
func (g Grating) LineDensityName() string {
	return formatDensity(g.GrooveDensity())
}

// IncidentSine returns sin(theta_i). It is the quantity that enters the
// grating equation and the order-existence inequality.
func (g Grating) IncidentSine() float64 {
	return math.Sin(g.IncidentAngleRad)
}
