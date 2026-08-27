package envelope

import (
	"fmt"
	"math"

	"grating-eq/internal/grating"
)

type HalfWidth struct {
	Order      int
	Slits      int
	Radians    float64
	Degrees    float64
	Defined    bool
	Cosine     float64
	GrazingCut bool
}

func PrincipalHalfWidth(g grating.Grating, m, n int) (HalfWidth, error) {
	out := HalfWidth{Order: m, Slits: n}
	if n < 1 {
		return out, grating.ErrBadSlits(n)
	}
	if m == 0 {
		return out, fmt.Errorf("envelope: zero order has no finite principal-peak width from the array factor")
	}
	res, err := g.Eq(m)
	if err != nil {
		return out, err
	}
	if !res.Exists {
		return out, fmt.Errorf("envelope: order %d is not visible", m)
	}
	cos := grating.CosineOf(res.AngleRad)
	out.Cosine = cos
	if math.Abs(cos) < 1e-12 {
		out.GrazingCut = true
		return out, fmt.Errorf("envelope: order %d sits at grazing, half-width diverges", m)
	}
	hw := g.WavelengthNm / (float64(n) * g.GrooveSpacingNm * math.Abs(cos))
	if math.IsNaN(hw) || math.IsInf(hw, 0) || hw <= 0 {
		return out, fmt.Errorf("envelope: half-width is not a positive finite value")
	}
	out.Radians = hw
	out.Degrees = grating.Degrees(hw)
	out.Defined = true
	return out, nil
}

func FirstZeroOffset(g grating.Grating, m, n int) (float64, error) {
	hw, err := PrincipalHalfWidth(g, m, n)
	if err != nil {
		return 0, err
	}
	return hw.Radians, nil
}

func FirstZeroPhase(n int) (float64, error) {
	if n < 1 {
		return 0, grating.ErrBadSlits(n)
	}
	return 2 * math.Pi / float64(n), nil
}

func WidthShrinksWithN(g grating.Grating, m, n1, n2 int) (bool, error) {
	if n2 <= n1 {
		return false, fmt.Errorf("envelope: n2 must be larger than n1")
	}
	a, err := PrincipalHalfWidth(g, m, n1)
	if err != nil {
		return false, err
	}
	b, err := PrincipalHalfWidth(g, m, n2)
	if err != nil {
		return false, err
	}
	if !a.Defined || !b.Defined {
		return false, fmt.Errorf("envelope: half-width undefined for the comparison")
	}
	ratio := a.Radians / b.Radians
	expect := float64(n2) / float64(n1)
	return math.Abs(ratio-expect) < 1e-9 && b.Radians < a.Radians, nil
}

func PeakToNullSeparation(g grating.Grating, m, n int) (float64, float64, error) {
	res, err := g.Eq(m)
	if err != nil {
		return 0, 0, err
	}
	if !res.Exists {
		return 0, 0, fmt.Errorf("envelope: order %d is not visible", m)
	}
	hw, err := PrincipalHalfWidth(g, m, n)
	if err != nil {
		return 0, 0, err
	}
	lo := res.AngleRad - hw.Radians
	hi := res.AngleRad + hw.Radians
	return lo, hi, nil
}
