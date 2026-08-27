package grating

import (
	"fmt"
	"math"
)

type Summary struct {
	SpacingNm float64

	DensityName string

	WavelengthNm float64

	IncidentDeg float64

	Slits int

	MaxOrder int
}

func (g Grating) Summarize() Summary {
	return Summary{
		SpacingNm:    g.GrooveSpacingNm,
		DensityName:  g.LineDensityName(),
		WavelengthNm: g.WavelengthNm,
		IncidentDeg:  Degrees(g.IncidentAngleRad),
		Slits:        g.Slits,
		MaxOrder:     g.MaxVisibleOrder(),
	}
}

func (s Summary) Header() string {
	return fmt.Sprintf(
		"grating: d=%.6g nm (%s)  lambda=%.6g nm  theta_i=%.6g deg  N=%d  max|m|=%d",
		s.SpacingNm, s.DensityName, s.WavelengthNm, s.IncidentDeg, s.Slits, s.MaxOrder)
}

type ScaleLaw struct {
	WavelengthUp bool

	SpacingDown bool

	ResolvingPowerRisesWithN bool
}

func (g Grating) ProbeScaleLaws(m int, lambdaStep, spacingStep float64) (ScaleLaw, error) {
	if lambdaStep <= 0 || spacingStep <= 0 {
		return ScaleLaw{}, ErrBadStep()
	}

	up := g
	up.WavelengthNm += lambdaStep
	resUp, err := up.Eq(m)
	if err != nil {
		return ScaleLaw{}, err
	}

	denser := g
	denser.GrooveSpacingNm = math.Max(MinGrooveSpacingNm, g.GrooveSpacingNm-spacingStep)
	resDense, err := denser.Eq(m)
	if err != nil {
		return ScaleLaw{}, err
	}

	law := ScaleLaw{ResolvingPowerRisesWithN: true}
	if resUp.Exists && g.WavelengthNm > 0 {
		base, err := g.Eq(m)
		if err != nil {
			return ScaleLaw{}, err
		}
		law.WavelengthUp = math.Abs(resUp.AngleRad) > math.Abs(base.AngleRad)
	}
	if resDense.Exists {
		base, err := g.Eq(m)
		if err != nil {
			return ScaleLaw{}, err
		}
		law.SpacingDown = math.Abs(resDense.AngleRad) > math.Abs(base.AngleRad)
	}
	return law, nil
}

func ErrBadStep() error {
	return fmt.Errorf("probe step must be positive")
}

func (g Grating) WavelengthRange(m int) (WavelengthRangeResult, error) {
	order := float64(m)
	si := g.IncidentSine()
	if order == 0 {
		return WavelengthRangeResult{Lo: 0, Hi: math.Inf(1)}, nil
	}
	hi := (1 - si) * g.GrooveSpacingNm / order
	lo := (-1 - si) * g.GrooveSpacingNm / order
	if hi < lo {
		hi, lo = lo, hi
	}
	lo = math.Max(0, lo)
	return WavelengthRangeResult{Lo: lo, Hi: hi}, nil
}

type WavelengthRangeResult struct {
	Lo float64

	Hi float64
}

func (w WavelengthRangeResult) Empty() bool {
	return w.Hi < w.Lo || math.IsNaN(w.Lo) || math.IsNaN(w.Hi)
}
