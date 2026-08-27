package envelope

import (
	"fmt"
	"math"

	"grating-eq/internal/grating"
)

type Groove struct {
	SpacingNm float64
	WidthNm   float64
}

func (g Groove) Validate() error {
	if math.IsNaN(g.SpacingNm) || math.IsInf(g.SpacingNm, 0) || g.SpacingNm <= 0 {
		return grating.ErrBadSpacing(g.SpacingNm)
	}
	if math.IsNaN(g.WidthNm) || math.IsInf(g.WidthNm, 0) || g.WidthNm <= 0 {
		return fmt.Errorf("envelope: groove width must be positive and finite, got %v", g.WidthNm)
	}
	if g.WidthNm > g.SpacingNm {
		return fmt.Errorf("envelope: groove width %v exceeds spacing %v", g.WidthNm, g.SpacingNm)
	}
	return nil
}

func (g Groove) Fill() (float64, error) {
	if err := g.Validate(); err != nil {
		return 0, err
	}
	return g.WidthNm / g.SpacingNm, nil
}

func (g Groove) Missing(m int) (bool, error) {
	if m == 0 {
		return false, nil
	}
	fill, err := g.Fill()
	if err != nil {
		return false, err
	}
	prod := math.Abs(float64(m)) * fill
	nearest := math.Round(prod)
	if nearest < 1 {
		return false, nil
	}
	return math.Abs(prod-nearest) < 1e-9, nil
}

func SingleSlitFactor(beta float64) (float64, error) {
	if math.IsNaN(beta) || math.IsInf(beta, 0) {
		return 0, fmt.Errorf("envelope: single-slit phase must be finite")
	}
	half := beta / 2
	if math.Abs(half) < 1e-14 {
		return 1, nil
	}
	return math.Sin(half) / half, nil
}

func SingleSlitIntensity(beta float64) (float64, error) {
	s, err := SingleSlitFactor(beta)
	if err != nil {
		return 0, err
	}
	return s * s, nil
}

func EnvelopeBeta(g Groove, lambdaNm, sinDiff float64) (float64, error) {
	if err := g.Validate(); err != nil {
		return 0, err
	}
	if lambdaNm <= 0 || math.IsNaN(lambdaNm) || math.IsInf(lambdaNm, 0) {
		return 0, grating.ErrBadWavelength(lambdaNm)
	}
	if math.IsNaN(sinDiff) || math.IsInf(sinDiff, 0) {
		return 0, fmt.Errorf("envelope: sine difference must be finite")
	}
	return 2 * math.Pi * (g.WidthNm / lambdaNm) * sinDiff, nil
}

func BetaAtOrder(g Groove, m int) (float64, error) {
	if err := g.Validate(); err != nil {
		return 0, err
	}
	fill, err := g.Fill()
	if err != nil {
		return 0, err
	}
	return 2 * math.Pi * float64(m) * fill, nil
}

func CombinedIntensity(n int, arrayDelta, slitBeta float64) (float64, error) {
	arr, err := ArrayIntensity(n, arrayDelta)
	if err != nil {
		return 0, err
	}
	slit, err := SingleSlitIntensity(slitBeta)
	if err != nil {
		return 0, err
	}
	return arr * slit, nil
}

func CombinedAtOrder(gr Groove, n, m int) (float64, error) {
	delta := PhaseAtOrder(m)
	beta, err := BetaAtOrder(gr, m)
	if err != nil {
		return 0, err
	}
	return CombinedIntensity(n, delta, beta)
}
