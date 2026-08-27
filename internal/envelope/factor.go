package envelope

import (
	"fmt"
	"math"

	"grating-eq/internal/grating"
)

func PhaseDelta(g grating.Grating, thetaRad float64) (float64, error) {
	if err := g.Validate(); err != nil {
		return 0, err
	}
	if math.IsNaN(thetaRad) || math.IsInf(thetaRad, 0) {
		return 0, fmt.Errorf("envelope: observation angle must be finite")
	}
	sObs, err := grating.SineOf(thetaRad)
	if err != nil {
		return 0, err
	}
	sInc := g.IncidentSine()
	return 2 * math.Pi * (g.GrooveSpacingNm / g.WavelengthNm) * (sObs - sInc), nil
}

func PhaseAtOrder(m int) float64 {
	return 2 * math.Pi * float64(m)
}

func ArrayFactor(n int, delta float64) (float64, error) {
	if n < 1 {
		return 0, grating.ErrBadSlits(n)
	}
	if math.IsNaN(delta) || math.IsInf(delta, 0) {
		return 0, fmt.Errorf("envelope: phase delta must be finite")
	}
	half := delta / 2
	s := math.Sin(half)
	if math.Abs(s) < 1e-14 {
		return float64(n), nil
	}
	return math.Sin(float64(n)*half) / s, nil
}

func ArrayIntensity(n int, delta float64) (float64, error) {
	af, err := ArrayFactor(n, delta)
	if err != nil {
		return 0, err
	}
	return af * af, nil
}

func PeakArrayIntensity(n int) (float64, error) {
	if n < 1 {
		return 0, grating.ErrBadSlits(n)
	}
	nn := float64(n)
	return nn * nn, nil
}

func NormalizedArray(n int, delta float64) (float64, error) {
	peak, err := PeakArrayIntensity(n)
	if err != nil {
		return 0, err
	}
	i, err := ArrayIntensity(n, delta)
	if err != nil {
		return 0, err
	}
	return i / peak, nil
}

func AngleFromPhase(g grating.Grating, delta float64) (float64, error) {
	if err := g.Validate(); err != nil {
		return 0, err
	}
	if math.IsNaN(delta) || math.IsInf(delta, 0) {
		return 0, fmt.Errorf("envelope: phase delta must be finite")
	}
	sInc := g.IncidentSine()
	sObs := sInc + (delta*g.WavelengthNm)/(2*math.Pi*g.GrooveSpacingNm)
	theta, err := grating.ArcSine(sObs)
	if err != nil {
		return 0, err
	}
	return theta, nil
}

func IntensityAtAngle(g grating.Grating, n int, thetaRad float64) (float64, error) {
	delta, err := PhaseDelta(g, thetaRad)
	if err != nil {
		return 0, err
	}
	return ArrayIntensity(n, delta)
}

func IntensityAtOrder(g grating.Grating, n, m int) (float64, error) {
	res, err := g.Eq(m)
	if err != nil {
		return 0, err
	}
	if !res.Exists {
		return 0, fmt.Errorf("envelope: order %d is not visible", m)
	}
	return IntensityAtAngle(g, n, res.AngleRad)
}
