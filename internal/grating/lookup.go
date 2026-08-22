package grating

import "math"

// Incidence describes how the light meets the grating, expressed in the
// two conventions a user might hand over: the angle from the normal, or
// the angle from the surface plane.
type Incidence struct {
	// FromNormalRad is theta_i measured from the grating normal.
	FromNormalRad float64

	// FromSurfaceRad is the complementary angle measured from the surface.
	FromSurfaceRad float64
}

// IncidenceFromNormal builds an Incidence from an angle measured from the
// normal. The surface angle is the complement.
func IncidenceFromNormal(normalAngleRad float64) (Incidence, error) {
	sine, err := SineOf(normalAngleRad)
	if err != nil {
		return Incidence{}, err
	}
	if math.Abs(sine) > 1 {
		return Incidence{}, ErrBadIncident("|sin(incident)| must not exceed 1")
	}
	return Incidence{
		FromNormalRad:  normalAngleRad,
		FromSurfaceRad: math.Pi/2 - normalAngleRad,
	}, nil
}

// FromSurfaceNormalized converts an angle measured from the surface into
// the angle from the normal that the grating equation expects. Mixing the
// two conventions is a classic error: a surface angle fed straight into the
// normal-based equation produces sin(theta_surface) instead of
// sin(theta_normal), which shifts every order and breaks the m=0 rule.
func FromSurfaceNormalized(surfaceAngleRad float64) (float64, error) {
	sine, err := SineOf(surfaceAngleRad)
	if err != nil {
		return 0, err
	}
	if math.Abs(sine) > 1 {
		return 0, ErrBadIncident("|sin(surface angle)| must not exceed 1")
	}
	return math.Pi/2 - surfaceAngleRad, nil
}

// NormalFromSurface is the inverse of FromSurfaceNormalized.
func NormalFromSurface(normalAngleRad float64) float64 {
	return math.Pi/2 - normalAngleRad
}

// GrazingError reports an incidence so close to grazing that the equation
// is numerically unreliable.
type GrazingError struct {
	FromSurfaceDeg float64
}

// Error implements the error interface.
func (e *GrazingError) Error() string {
	return "incidence is grazing (surface angle " + formatNumber(e.FromSurfaceDeg, 4) +
		" deg); normal-based equation is numerically unstable there"
}

// LookupAngle computes theta_m for a wavelength and order, given spacing
// and incident angle. This is the inverse-direction helper used when a
// spectrum is swept instead of a fixed wavelength.
func LookupAngle(spacingNm, wavelengthNm, incidentRad float64, m int) (float64, error) {
	g := Grating{
		GrooveSpacingNm:  spacingNm,
		WavelengthNm:     wavelengthNm,
		IncidentAngleRad: incidentRad,
	}
	if err := g.Validate(); err != nil {
		return 0, err
	}
	res, err := g.Eq(m)
	if err != nil {
		return 0, err
	}
	if !res.Exists {
		return 0, ErrSineOutOfRange(res.SineOfAngle)
	}
	return res.AngleRad, nil
}

// WavelengthForOrder solves the grating equation for the wavelength that
// places order m at the requested angle. It returns the wavelength in
// nanometres.
func WavelengthForOrder(spacingNm, incidentRad, targetAngleRad float64, m int) (float64, error) {
	if m == 0 {
		// The zero order sits at the incident angle for every wavelength.
		if math.Abs(targetAngleRad-incidentRad) > AngleTolerance {
			return 0, ErrBadIncident("zero order always follows the incident direction")
		}
		return 0, nil
	}
	sTarget := math.Sin(targetAngleRad)
	lambda := spacingNm * (sTarget - math.Sin(incidentRad)) / float64(m)
	if lambda <= 0 {
		return 0, ErrBadWavelength(lambda)
	}
	return lambda, nil
}
