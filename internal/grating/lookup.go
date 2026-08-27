package grating

import "math"

type Incidence struct {
	FromNormalRad float64

	FromSurfaceRad float64
}

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

func NormalFromSurface(normalAngleRad float64) float64 {
	return math.Pi/2 - normalAngleRad
}

type GrazingError struct {
	FromSurfaceDeg float64
}

func (e *GrazingError) Error() string {
	return "incidence is grazing (surface angle " + formatNumber(e.FromSurfaceDeg, 4) +
		" deg); normal-based equation is numerically unstable there"
}

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

func WavelengthForOrder(spacingNm, incidentRad, targetAngleRad float64, m int) (float64, error) {
	if m == 0 {
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
