package grating

import (
	"fmt"
	"math"
)

// Summary is a compact human-readable digest of the grating geometry used
// as the header of CLI reports.
type Summary struct {
	// SpacingNm is the groove spacing in nanometres.
	SpacingNm float64

	// DensityName is the formatted line density, e.g. "600 grooves/mm".
	DensityName string

	// WavelengthNm is the illumination wavelength.
	WavelengthNm float64

	// IncidentDeg is the incident angle in degrees.
	IncidentDeg float64

	// Slits is the finite slit count (0 = unbounded).
	Slits int

	// MaxOrder is the largest |m| that still exists for this geometry.
	MaxOrder int
}

// Summarize builds a Summary from the grating.
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

// Header renders a single-line header for tables.
func (s Summary) Header() string {
	return fmt.Sprintf(
		"grating: d=%.6g nm (%s)  lambda=%.6g nm  theta_i=%.6g deg  N=%d  max|m|=%d",
		s.SpacingNm, s.DensityName, s.WavelengthNm, s.IncidentDeg, s.Slits, s.MaxOrder)
}

// ScaleLaw describes how a change to one input moves the order angles. It
// turns the qualitative cross rules into a queryable quantity so tests can
// assert strict monotonicity.
type ScaleLaw struct {
	// WavelengthUp makes |theta_m| larger for a fixed order.
	WavelengthUp bool

	// SpacingDown makes |theta_m| larger when the grating is denser.
	SpacingDown bool

	// ResolvingPowerRisesWithN reports whether R = m*N grows linearly.
	ResolvingPowerRisesWithN bool
}

// ProbeScaleLaws evaluates the two monotonicity rules at the current
// geometry for a given order m. The rules are stated as inequalities that
// must hold for every realizable configuration.
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

// ErrBadStep reports a non-positive probe step.
func ErrBadStep() error {
	return fmt.Errorf("probe step must be positive")
}

// WavelengthRange returns the wavelength window for which a given order
// remains visible, given fixed spacing and incident angle. It is a closed
// interval [lo, hi] in nanometres; an empty interval means the order never
// appears at any wavelength.
func (g Grating) WavelengthRange(m int) (WavelengthRangeResult, error) {
	order := float64(m)
	si := g.IncidentSine()
	if order == 0 {
		return WavelengthRangeResult{Lo: 0, Hi: math.Inf(1)}, nil
	}
	// |m*lambda/d + sin(theta_i)| <= 1.
	hi := (1 - si) * g.GrooveSpacingNm / order
	lo := (-1 - si) * g.GrooveSpacingNm / order
	if hi < lo {
		hi, lo = lo, hi
	}
	lo = math.Max(0, lo)
	return WavelengthRangeResult{Lo: lo, Hi: hi}, nil
}

// WavelengthRangeResult carries the visibility window of one order.
type WavelengthRangeResult struct {
	// Lo is the lower wavelength bound in nanometres.
	Lo float64

	// Hi is the upper wavelength bound in nanometres.
	Hi float64
}

// Empty reports whether the interval is empty.
func (w WavelengthRangeResult) Empty() bool {
	return w.Hi < w.Lo || math.IsNaN(w.Lo) || math.IsNaN(w.Hi)
}
