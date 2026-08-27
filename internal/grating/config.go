package grating

import "math"

const (
	NanoPerMilli = 1e6

	MaxWavelengthNm = 1e12

	MinGrooveSpacingNm = 0.1

	DefaultMaxOrderScan = 64

	AngleTolerance = 1e-9
)

type Grating struct {
	GrooveSpacingNm float64

	WavelengthNm float64

	IncidentAngleRad float64

	Slits int
}

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

func (g Grating) GrooveDensity() float64 {
	return NanoPerMilli / g.GrooveSpacingNm
}

func SpacingFromDensity(groovesPerMm float64) float64 {
	return NanoPerMilli / groovesPerMm
}

func (g Grating) LineDensityName() string {
	return formatDensity(g.GrooveDensity())
}

func (g Grating) IncidentSine() float64 {
	return math.Sin(g.IncidentAngleRad)
}
