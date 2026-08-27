package disperse

import (
	"math"

	"grating-eq/internal/grating"
)

type ScanPoint struct {
	WavelengthNm float64

	AngleRad float64

	Exists bool

	DispersionDegNm float64
}

func ScanOrder(g grating.Grating, m, samples int, lo, hi float64) ([]ScanPoint, error) {
	if samples < 2 {
		samples = 2
	}
	if hi < lo {
		lo, hi = hi, lo
	}
	out := make([]ScanPoint, 0, samples)
	for i := 0; i < samples; i++ {
		frac := float64(i) / float64(samples-1)
		lambda := lo + frac*(hi-lo)
		gg := g
		gg.WavelengthNm = lambda
		res, err := gg.Eq(m)
		if err != nil {
			return nil, err
		}
		p := ScanPoint{WavelengthNm: lambda, Exists: res.Exists}
		if res.Exists {
			p.AngleRad = res.AngleRad
			d := AngularDispersion(gg, res)
			p.DispersionDegNm = d.ValueDegNm
		}
		out = append(out, p)
	}
	return out, nil
}

type BandSummary struct {
	Order int

	Points int

	MaxAbsAngleDeg float64

	FadesOut bool

	PeakDispersionDegNm float64
}

func SummarizeBand(points []ScanPoint, m int) BandSummary {
	s := BandSummary{Order: m, Points: len(points)}
	for _, p := range points {
		if !p.Exists {
			s.FadesOut = true
			continue
		}
		if a := math.Abs(grating.Degrees(p.AngleRad)); a > s.MaxAbsAngleDeg {
			s.MaxAbsAngleDeg = a
		}
		if d := math.Abs(p.DispersionDegNm); d > s.PeakDispersionDegNm {
			s.PeakDispersionDegNm = d
		}
	}
	return s
}

func WavelengthForAngleInverse(g grating.Grating, m int, targetAngleRad float64) (float64, error) {
	if m == 0 {
		if math.Abs(targetAngleRad-g.IncidentAngleRad) > grating.AngleTolerance {
			return 0, grating.ErrBadIncident("zero order is pinned to the incident angle")
		}
		return 0, nil
	}
	sTarget := math.Sin(targetAngleRad)
	sInc := math.Sin(g.IncidentAngleRad)
	lambda := g.GrooveSpacingNm * (sTarget - sInc) / float64(m)
	if lambda <= 0 || math.IsNaN(lambda) || math.IsInf(lambda, 0) {
		return 0, grating.ErrBadWavelength(lambda)
	}
	return lambda, nil
}

func InverseScan(g grating.Grating, m, samples int) ([]ScanPoint, error) {
	if samples < 2 {
		samples = 2
	}
	out := make([]ScanPoint, 0, samples)
	for i := 0; i < samples; i++ {
		frac := float64(i) / float64(samples-1)
		theta := -math.Pi/2 + frac*math.Pi
		lambda, err := WavelengthForAngleInverse(g, m, theta)
		if err != nil {
			continue
		}
		out = append(out, ScanPoint{
			WavelengthNm: lambda,
			AngleRad:     theta,
			Exists:       true,
		})
	}
	return out, nil
}
