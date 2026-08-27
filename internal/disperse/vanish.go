package disperse

import (
	"math"

	"grating-eq/internal/grating"
)

type VanishReport struct {
	Existing int

	MaxOrder int

	FirstMissing int

	Ratio float64
}

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
	rep.Existing *= 2
	rep.MaxOrder = maxOrder
	rep.FirstMissing = firstMissing
	return rep, nil
}

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

func CriticalDensity(wavelengthNm float64, m int) (float64, error) {
	if wavelengthNm <= 0 {
		return 0, grating.ErrBadWavelength(wavelengthNm)
	}
	if m == 0 {
		return math.Inf(1), nil
	}
	d := math.Abs(float64(m)) * wavelengthNm
	return grating.NanoPerMilli / d, nil
}

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
