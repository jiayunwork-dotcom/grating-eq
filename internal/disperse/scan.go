package disperse

import (
	"math"

	"grating-eq/internal/grating"
)

// ScanPoint is one wavelength sample in a spectral scan.
type ScanPoint struct {
	// WavelengthNm is the sampled wavelength.
	WavelengthNm float64

	// AngleRad is theta_m at that wavelength.
	AngleRad float64

	// Exists reports whether the order survives at this wavelength.
	Exists bool

	// DispersionDegNm is the angular dispersion at this wavelength.
	DispersionDegNm float64
}

// ScanOrder samples order m across a wavelength window [lo, hi] with a
// given number of samples, at normal incidence. Orders that leave the
// visible region inside the window are reported with Exists=false so a
// caller can show where the order fades out.
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

// BandSummary aggregates a spectral scan of order m.
type BandSummary struct {
	// Order is the scanned order.
	Order int

	// Points is the number of samples.
	Points int

	// MaxAbsAngleDeg is the largest |theta| seen among existing points.
	MaxAbsAngleDeg float64

	// FadesOut reports whether the order vanished inside the window.
	FadesOut bool

	// PeakDispersionDegNm is the largest |dtheta/dlambda| seen.
	PeakDispersionDegNm float64
}

// SummarizeBand condenses a scan into the quantities a report prints.
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

// WavelengthForAngleInverse is the inverse scan: given a target angle and
// order, return the wavelength that lands there. It inverts
// theta(lambda) = asin(m*lambda/d + sin(theta_i)).
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

// InverseScan builds a table of wavelength versus angle for one order,
// sweeping the angle across the visible region.
func InverseScan(g grating.Grating, m, samples int) ([]ScanPoint, error) {
	if samples < 2 {
		samples = 2
	}
	// Angle window: (-pi/2, pi/2) is the full visible region.
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
