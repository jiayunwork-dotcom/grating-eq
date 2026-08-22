package disperse

import (
	"math"

	"grating-eq/internal/grating"
)

// VanishReport describes how many orders survive on a grating, and where
// the highest existing order sits.
type VanishReport struct {
	// Existing is the number of orders with |sin(theta_m)| <= 1.
	Existing int

	// MaxOrder is the largest |m| that still exists.
	MaxOrder int

	// FirstMissing is the smallest |m| that no longer exists (0 when the
	// scan found no missing order).
	FirstMissing int

	// Ratio is lambda/d; larger values hide more orders.
	Ratio float64
}

// VanishAnalysis scans orders up to a limit and reports the dense-grating
// behaviour: as lambda/d grows towards 1 the first order approaches the
// grazing angle, and beyond 1 all non-zero orders disappear.
func VanishAnalysis(g grating.Grating, limit int) (VanishReport, error) {
	if limit < 1 {
		limit = grating.DefaultMaxOrderScan
	}
	rep := VanishReport{Ratio: g.WavelengthNm / g.GrooveSpacingNm}
	maxOrder := 0
	firstMissing := 0
	for m := 1; m <= limit; m++ {
		res, err := g.Eq(m)
		if err != nil {
			return VanishReport{}, err
		}
		if res.Exists {
			rep.Existing++
			if m > maxOrder {
				maxOrder = m
			}
		} else if firstMissing == 0 {
			firstMissing = m
		}
	}
	rep.Existing *= 2 // +m and -m share the same fate at normal incidence.
	rep.MaxOrder = maxOrder
	rep.FirstMissing = firstMissing
	return rep, nil
}

// EdgeAngle returns the angle of the last surviving order before the
// grating goes dense, in degrees. It is the observable symptom of the
// "high orders vanish" cross rule.
func EdgeAngle(g grating.Grating, limit int) (float64, error) {
	rep, err := VanishAnalysis(g, limit)
	if err != nil {
		return 0, err
	}
	if rep.MaxOrder == 0 {
		return 0, nil
	}
	res, err := g.Eq(rep.MaxOrder)
	if err != nil {
		return 0, err
	}
	if !res.Exists {
		return 0, nil
	}
	return grating.Degrees(res.AngleRad), nil
}

// CriticalDensity returns the groove density in grooves/mm at which order
// m reaches the grazing angle for a given wavelength. Beyond it the order
// disappears. This makes the vanishing rule quantitative.
func CriticalDensity(wavelengthNm float64, m int) (float64, error) {
	if wavelengthNm <= 0 {
		return 0, grating.ErrBadWavelength(wavelengthNm)
	}
	if m == 0 {
		return math.Inf(1), nil
	}
	// lambda/d * |m| = 1  =>  d = |m|*lambda, density = 1e6/(|m|*lambda).
	d := math.Abs(float64(m)) * wavelengthNm
	return grating.NanoPerMilli / d, nil
}

// DeepestOrder returns the order |m| that sits closest to the grazing
// angle at the current geometry.
func DeepestOrder(g grating.Grating, limit int) (int, float64, error) {
	if limit < 1 {
		limit = grating.DefaultMaxOrderScan
	}
	best := 0
	bestSine := 0.0
	for m := 1; m <= limit; m++ {
		res, err := g.Eq(m)
		if err != nil {
			return 0, 0, err
		}
		if !res.Exists {
			continue
		}
		absSine := math.Abs(res.SineOfAngle)
		if absSine > bestSine {
			bestSine = absSine
			best = m
		}
	}
	return best, bestSine, nil
}
